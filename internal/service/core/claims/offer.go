package claims

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-schema-processor/verifiable"
	"github.com/iden3/iden3comm/packers"
	"github.com/iden3/iden3comm/protocol"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
)

func NewClaimOffer(callBackURL string, from, to *core.ID, claim *data.Claim) *protocol.CredentialsOfferMessage {
	return &protocol.CredentialsOfferMessage{
		ID:       uuid.NewString(),
		Typ:      packers.MediaTypePlainMessage,
		Type:     protocol.CredentialOfferMessageType,
		ThreadID: uuid.NewString(),
		Body: protocol.CredentialsOfferMessageBody{
			URL: callBackURL,
			Credentials: []protocol.CredentialOffer{
				protocol.CredentialOffer{
					ID:          fmt.Sprint(claim.ID),
					Description: claim.SchemaType,
				},
			},
		},
		From: from.String(),
		To:   to.String(),
	}
}

func ClaimOfferToRaw(claimOffer *protocol.CredentialsOfferMessage, createdAt time.Time) *data.ClaimOffer {
	claimOfferRaw := data.ClaimOffer{
		ID:        claimOffer.ID,
		From:      claimOffer.From,
		To:        claimOffer.To,
		CreatedAt: createdAt,
	}

	if len(claimOffer.Body.Credentials) > 0 {
		claimOfferRaw.ClaimID = claimOffer.Body.Credentials[0].ID
	}

	return &claimOfferRaw
}

func ClaimModelToIden3Credential(claim *data.Claim) (*verifiable.Iden3Credential, error) {
	claimIDPosition, err := getClaimIDPosition(claim.CoreClaim.Claim)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get claim id position")
	}

	credentialData, err := compactCredentialData(claim)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compact credential subject")
	}

	proofs, err := compactProofs(claim)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compact proofs")
	}

	credentialStatus := &verifiable.CredentialStatus{}
	if claim.CredentialStatus != nil && string(claim.CredentialStatus) != "{}" {
		if err := json.Unmarshal(claim.CredentialStatus, credentialStatus); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal credential status")
		}
	}

	res := &verifiable.Iden3Credential{
		ID:        fmt.Sprint(claim.ID),
		Type:      []string{Iden3CredentialSchema.ToRaw()},
		RevNonce:  claim.CoreClaim.GetRevocationNonce(),
		Updatable: claim.CoreClaim.GetFlagUpdatable(),
		Version:   claim.CoreClaim.GetVersion(),
		Context: []string{
			ClaimSchemaList[ClaimSchemaTypeList[claim.SchemaType]].ClaimSchemaURL,
			claim.SchemaURL,
		},
		CredentialSchema: struct {
			ID   string `json:"@id"`
			Type string `json:"type"`
		}{
			ID:   claim.SchemaURL,
			Type: claim.SchemaType,
		},
		SubjectPosition:   claimIDPosition,
		CredentialSubject: credentialData,
		Proof:             proofs,
		CredentialStatus:  credentialStatus,
	}

	expiration, ok := claim.CoreClaim.GetExpirationDate()
	if ok {
		res.Expiration = expiration
	}

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

func compactProofs(claim *data.Claim) ([]interface{}, error) {
	proofs := make([]interface{}, 0)

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

func compactCredentialData(claim *data.Claim) (map[string]interface{}, error) {
	subjectID, err := claim.CoreClaim.GetID()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get subject id")
	}

	credentialData := make(map[string]interface{})
	if err := json.Unmarshal(claim.Data, &credentialData); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal credential subject")
	}

	credentialData["type"] = claim.SchemaType
	if len(subjectID.String()) > 0 {
		credentialData["id"] = subjectID.String()
	}

	return credentialData, nil
}
