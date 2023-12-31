package issuer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-schema-processor/verifiable"
	"github.com/iden3/iden3comm/packers"
	"github.com/iden3/iden3comm/protocol"
	"github.com/pkg/errors"

	"github.com/rarimo/issuer/internal/data"
	"github.com/rarimo/issuer/internal/service/api/requests"
	"github.com/rarimo/issuer/internal/service/core/claims"
	"github.com/rarimo/issuer/internal/service/core/identity/state"
)

func (isr *issuer) CreateClaimOffer(
	userDID *core.DID, claimID string,
) (*protocol.CredentialsOfferMessage, error) {
	claim, err := isr.Identity.State.DB.ClaimsQ().GetBySchemaType(claimID, userDID.ID.String())
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
	if !userDID.ID.Equals(&recipient) {
		return nil, ErrClaimRetrieverIsNotClaimOwner
	}

	claimOffer := NewClaimOffer(
		fmt.Sprint(isr.baseURL, ClaimIssueCallBackPath), isr.Identifier, userDID, claim,
	)

	err = isr.claimsOffersQ.Insert(
		ClaimOfferToRaw(claimOffer, time.Now(), isr.Identifier.ID, userDID.ID),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert claim offer to db")
	}

	return claimOffer, nil
}

func (isr *issuer) IssueClaim(
	ctx context.Context,
	userDID *core.DID,
	expiration *time.Time,
	claimType claims.ClaimSchemaType,
	credentialRaw []byte,
) (string, error) {
	claim, err := isr.compactClaim(ctx, userDID, expiration, claimType, credentialRaw)
	if err != nil {
		return "", errors.Wrap(err, "failed to compact the requested claim")
	}

	hi, hv, err := claim.CoreClaim.HiHv()
	if err != nil {
		return "", errors.Wrap(err, "failed to get claim index and value hash")
	}

	// we clone it to save the atomicity of the TX and to
	// avoid this issue https://github.com/lib/pq/issues/635
	// merkletree.Add() contains update statement after the select
	db := isr.State.DB.New()
	err = db.Transaction(func() error {
		err = db.ClaimsQ().Insert(claim)
		if err != nil {
			return errors.Wrap(err, "failed to insert claim into db")
		}

		err = isr.State.ClaimsTree.Add(ctx, hi, hv)
		if err != nil {
			return errors.Wrap(err, "failed to add claim to the claims merkle tree")
		}

		return nil
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to execute db transaction")
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

	claim, err := isr.State.DB.ClaimsQ().Get(request.FetchMessage.Body.ID)
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

	cred, err := ClaimModelToW3Credential(claim)
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
		Typ:      packers.MediaTypePlainMessage,
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
	lastCommittedStateRaw, err := isr.State.DB.CommittedStatesQ().WhereStatus(data.StatusCompleted).GetLatest()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last committed state")
	}

	lastCommittedState, err := state.CommittedStateFromRaw(lastCommittedStateRaw)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last committed state")
	}

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
			RootOfRoots        *string `json:"rootOfRoots,omitempty"`
			ClaimsTreeRoot     *string `json:"claimsTreeRoot,omitempty"`
			RevocationTreeRoot *string `json:"revocationTreeRoot,omitempty"`
		}{
			State:              strptr(stateHash.Hex()),
			RevocationTreeRoot: strptr(lastCommittedState.RevocationsTreeRoot.Hex()),
			RootOfRoots:        strptr(lastCommittedState.RootsTreeRoot.Hex()),
			ClaimsTreeRoot:     strptr(lastCommittedState.ClaimsTreeRoot.Hex()),
		},
	}, nil
}

func (isr *issuer) GetInclusionMTP(
	ctx context.Context,
	claimID uuid.UUID,
) (*ClaimInclusionMTP, error) {
	lastCommittedStateRaw, err := isr.State.DB.CommittedStatesQ().WhereStatus(data.StatusCompleted).GetLatest()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last committed state")
	}

	lastCommittedState, err := state.CommittedStateFromRaw(lastCommittedStateRaw)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last committed state")
	}

	claim, err := isr.State.DB.ClaimsQ().Get(claimID.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get claim from db")
	}

	hi, err := claim.CoreClaim.HIndex()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get claim index hash")
	}

	mtp, _, err := isr.State.ClaimsTree.GenerateProof(ctx, hi, lastCommittedState.ClaimsTreeRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate claim inclusion proof")
	}

	stateHash, err := lastCommittedState.StateHash()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate state hash")
	}

	return &ClaimInclusionMTP{
		MTP: *mtp,
		Issuer: struct {
			State              *string `json:"state,omitempty"`
			RootOfRoots        *string `json:"rootOfRoots,omitempty"`
			ClaimsTreeRoot     *string `json:"claimsTreeRoot,omitempty"`
			RevocationTreeRoot *string `json:"revocationTreeRoot,omitempty"`
		}{
			State:              strptr(stateHash.Hex()),
			RevocationTreeRoot: strptr(lastCommittedState.RevocationsTreeRoot.Hex()),
			RootOfRoots:        strptr(lastCommittedState.RootsTreeRoot.Hex()),
			ClaimsTreeRoot:     strptr(lastCommittedState.ClaimsTreeRoot.Hex()),
		},
	}, nil
}

func (isr *issuer) GetIdentifier() string {
	return isr.Identifier.ID.String()
}

func (isr *issuer) RevokeClaim(
	ctx context.Context,
	userID *core.ID,
	schemaType claims.ClaimSchemaType,
) error {
	claim, err := isr.State.DB.ClaimsQ().GetBySchemaType(schemaType.ToRaw(), userID.String())
	if err != nil {
		return errors.Wrap(err, "failed to get claim from db")
	}
	if claim == nil {
		return ErrClaimIsNotExist
	}

	if claim.Revoked {
		return ErrClaimIsAlreadyRevoked
	}

	claim.Revoked = true

	// we clone it to save the atomicity of the TX and to
	// avoid this issue https://github.com/lib/pq/issues/635
	// merkletree.Add() contains update statement after the select
	db := isr.State.DB.New()
	err = db.Transaction(func() error {
		err = db.ClaimsQ().Update(claim)
		if err != nil {
			return errors.Wrap(err, "failed to update claim in db")
		}

		err = isr.State.RevocationsTree.Add(ctx, new(big.Int).SetUint64(
			claim.CoreClaim.GetRevocationNonce()), merkletree.HashZero.BigInt(),
		)
		if err != nil {
			return errors.Wrap(err, "failed to add revocation nonce to the revocations merkle tree")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to execute db transaction")
	}

	return nil
}

func (isr *issuer) checkCallbackRequest(
	claim *data.Claim,
	claimOffer *data.ClaimOffer,
	request *requests.OfferCallbackRequest,
) error {
	userDID, err := core.ParseDID(request.FetchMessage.From)
	if err != nil {
		return errors.Wrap(err, "failed to parse user did")
	}

	ok, err := isr.checkClaimRetriever(claim, userDID.ID.String(), request.Token)
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
