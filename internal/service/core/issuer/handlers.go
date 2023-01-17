package issuer

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/iden3comm/protocol"
	"github.com/pkg/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims"
)

func (isr *issuer) CreateClaimOffer(
	userID *core.ID, claimID uint64,
) (*protocol.CredentialsOfferMessage, error) {
	claim, err := isr.Identity.State.ClaimsQ.Get(claimID)
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
	schemaType claims.ClaimSchemaType,
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
	fetchMessage *protocol.CredentialFetchRequestMessage,
) (*protocol.CredentialIssuanceMessage, error) {
	if isr.Identifier.String() != fetchMessage.To {
		return nil, ErrMessageRecipientIsNotIssuer
	}

	claim, err := isr.retrieveClaim(fetchMessage.Body.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve claim by fetch message id")
	}

	ok, err := checkClaimRetriever(claim, fetchMessage.From)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check claim retriever")
	}
	if !ok {
		return nil, ErrClaimRetrieverIsNotClaimOwner
	}

	if err := isr.generateMTPForOfferCallback(ctx, claim); err != nil {
		return nil, errors.Wrap(err, "failed to generate mtp")
	}

	cred, err := claims.ClaimModelToIden3Credential(claim)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create iden3 credential from claim model")
	}

	return &protocol.CredentialIssuanceMessage{
		ID:       uuid.NewString(),
		Type:     protocol.CredentialIssuanceResponseMessageType,
		ThreadID: fetchMessage.ThreadID,
		Body:     protocol.IssuanceMessageBody{Credential: *cred},
		From:     fetchMessage.From,
		To:       fetchMessage.To,
	}, nil
}
