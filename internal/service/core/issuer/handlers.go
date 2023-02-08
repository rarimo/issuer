package issuer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-schema-processor/verifiable"
	"github.com/iden3/iden3comm/protocol"
	"github.com/pkg/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/api/requests"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state"
	resources "gitlab.com/q-dev/q-id/resources/claim_resources"
)

func (isr *issuer) CreateClaimOffer(
	userID *core.ID, claimID string,
) (*protocol.CredentialsOfferMessage, error) {
	claim, err := isr.Identity.State.ClaimsQ.GetBySchemaType(claimID, userID.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get claim from db")
	}
	if claim == nil {
		return nil, ErrClaimIsNotExist
	}

	recipient, err := claim.CoreClaim.GetID()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the claim recipient id")
	}
	if !userID.Equals(&recipient) {
		return nil, ErrClaimRetrieverIsNotClaimOwner
	}

	claimOffer := claims.NewClaimOffer(
		fmt.Sprint(isr.domain, ClaimIssueCallBackPath), isr.Identifier, userID, claim,
	)

	claimOfferRaw := claims.ClaimOfferToRaw(claimOffer, time.Now())
	err = isr.claimsOffersQ.Insert(claimOfferRaw)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert claim offer to db")
	}

	return claimOffer, nil
}

func (isr *issuer) IssueClaim(
	ctx context.Context,
	userID *core.ID,
	expiration time.Time,
	schemaType resources.ClaimSchemaType,
	schemaData []byte,
) (uint64, error) {
	claim, err := isr.compactClaim(ctx, userID, expiration, schemaType, schemaData)
	if err != nil {
		return 0, errors.Wrap(err, "failed to compact the requested claim")
	}

	hi, hv, err := claim.CoreClaim.HiHv()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get claim index and value hash")
	}

	err = isr.State.ClaimsQ.Insert(claim)
	if err != nil {
		return 0, errors.Wrap(err, "failed to insert claim into db")
	}

	err = isr.State.ClaimsTree.Add(ctx, hi, hv)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add claim to the claims merkle tree")
	}

	return claim.ID, nil
}

func (isr *issuer) OfferCallback(
	ctx context.Context,
	request *requests.OfferCallbackRequest,
) (*protocol.CredentialIssuanceMessage, error) {
	claimOffer, err := isr.claimsOffersQ.Get(request.FetchMessage.ThreadID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get claim offer")
	}
	if claimOffer == nil {
		return nil, ErrClaimOfferIsNotExist
	}

	claim, err := isr.State.ClaimsQ.GetBySchemaType(request.FetchMessage.Body.ID, request.FetchMessage.From)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get claim from db")
	}
	if claim == nil {
		return nil, ErrClaimIsNotExist
	}

	err = isr.checkCallbackRequest(claim, claimOffer, request)
	if err != nil {
		return nil, errors.Wrap(err, "invalid callback request")
	}

	if err := isr.generateProofs(ctx, claim); err != nil {
		return nil, errors.Wrap(err, "failed to generate mtp")
	}

	cred, err := claims.ClaimModelToIden3Credential(claim)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create iden3 credential from claim model")
	}

	claimOffer.IsReceived = true
	err = isr.claimsOffersQ.Update(claimOffer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update claim offer in db")
	}

	return &protocol.CredentialIssuanceMessage{
		ID:       uuid.NewString(),
		Type:     protocol.CredentialIssuanceResponseMessageType,
		ThreadID: request.FetchMessage.ThreadID,
		Body:     protocol.IssuanceMessageBody{Credential: *cred},
		From:     request.FetchMessage.To,
		To:       request.FetchMessage.From,
	}, nil
}

func (isr *issuer) GetRevocationStatus(
	ctx context.Context,
	revID *big.Int,
) (*verifiable.RevocationStatus, error) {
	lastCommittedStateRaw, err := isr.State.CommittedStateQ.WhereStatus(data.StatusCompleted).GetLatest()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last committed state")
	}

	lastCommittedState := state.CommittedStateFromRaw(lastCommittedStateRaw)

	mtp, _, err := isr.State.RevocationsTree.GenerateProof(ctx, revID, lastCommittedState.RevocationsTreeRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate non revocation proof")
	}

	stateHash, err := lastCommittedState.StateHash()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate state hash")
	}

	return &verifiable.RevocationStatus{
		MTP: *mtp,
		Issuer: struct {
			State              *string `json:"state,omitempty"`
			RootOfRoots        *string `json:"root_of_roots,omitempty"`
			ClaimsTreeRoot     *string `json:"claims_tree_root,omitempty"`
			RevocationTreeRoot *string `json:"revocation_tree_root,omitempty"`
		}{
			State:              strptr(stateHash.Hex()),
			RevocationTreeRoot: strptr(lastCommittedState.RevocationsTreeRoot.Hex()),
			RootOfRoots:        strptr(lastCommittedState.RootsTreeRoot.Hex()),
			ClaimsTreeRoot:     strptr(lastCommittedState.ClaimsTreeRoot.Hex()),
		},
	}, nil
}

func (isr *issuer) GetIdentifier() string {
	return isr.Identifier.String()
}

func (isr *issuer) checkCallbackRequest(
	claim *data.Claim,
	claimOffer *data.ClaimOffer,
	request *requests.OfferCallbackRequest,
) error {
	ok, err := isr.checkClaimRetriever(claim, request.FetchMessage.From, request.Token)
	if err != nil {
		return errors.Wrap(err, "failed to check claim retriever")
	}
	if !ok {
		return ErrClaimRetrieverIsNotClaimOwner
	}

	if isr.Identifier.String() != request.FetchMessage.To {
		return ErrMessageRecipientIsNotIssuer
	}

	if claimOffer.IsReceived {
		return ErrRepeatedCallbackRequest
	}

	return nil
}
