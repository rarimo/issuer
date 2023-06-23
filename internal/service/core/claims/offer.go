package claims

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-schema-processor/verifiable"
	"github.com/iden3/iden3comm/packers"
	"github.com/iden3/iden3comm/protocol"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/rarimo/identity/issuer/internal/data"
)

func NewClaimOffer(callBackURL string, from, to *core.DID, claim *data.Claim) *protocol.CredentialsOfferMessage {
	return &protocol.CredentialsOfferMessage{
		ID:       uuid.NewString(),
		Typ:      packers.MediaTypePlainMessage,
		Type:     protocol.CredentialOfferMessageType,
		ThreadID: uuid.NewString(),
		Body: protocol.CredentialsOfferMessageBody{
			URL: callBackURL,
			Credentials: []protocol.CredentialOffer{
				protocol.CredentialOffer{
					ID:          claim.ID,
					Description: ClaimSchemaList[ClaimSchemaTypeList[claim.ClaimType]].ClaimSchemaName,
				},
			},
		},
		From: from.String(),
		To:   to.String(),
	}
}

func ClaimOfferToRaw(
	claimOffer *protocol.CredentialsOfferMessage,
	createdAt time.Time,
	fromID core.ID,
	toID core.ID,
) *data.ClaimOffer {
	claimOfferRaw := data.ClaimOffer{
		ID:         claimOffer.ThreadID,
		From:       fromID.String(),
		To:         toID.String(),
		CreatedAt:  createdAt,
		IsReceived: false,
	}

	if len(claimOffer.Body.Credentials) > 0 {
		claimOfferRaw.ClaimID = claimOffer.Body.Credentials[0].ID
	}

	return &claimOfferRaw
}

func ClaimModelToW3Credential(claim *data.Claim) (*verifiable.W3CCredential, error) {
	proofs, err := compactProofs(claim)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compact proofs")
	}

	res := &verifiable.W3CCredential{}
	err = json.Unmarshal(claim.Credential, res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal credential")
	}

	res.Proof = proofs

	return res, nil
}

func getClaimIDPosition(claim *core.Claim) (string, error) {
	claimIDPosition, err := claim.GetIDPosition()
	if err != nil {
		return "", errors.Wrap(err, "failed to get subject id position")
	}

	switch claimIDPosition {
	case core.IDPositionIndex:
		return SubjectPositionIndex, nil
	case core.IDPositionValue:
		return SubjectPositionValue, nil
	default:
		return "", ErrIDPositionIsNotSpecified
	}
}

func compactProofs(claim *data.Claim) (verifiable.CredentialProofs, error) {
	proofs := verifiable.CredentialProofs{}

	signatureProof := &verifiable.BJJSignatureProof2021{}
	if claim.SignatureProof != nil {
		if err := json.Unmarshal(claim.SignatureProof, signatureProof); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal signature proof")
		}

		proofs = append(proofs, signatureProof)
	}

	mtp := &verifiable.Iden3SparseMerkleProof{}
	if claim.MTP != nil {
		if err := json.Unmarshal(claim.MTP, mtp); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal merkle tree proof")
		}

		proofs = append(proofs, mtp)
	}

	return proofs, nil
}
