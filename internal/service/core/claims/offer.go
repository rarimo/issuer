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
	res := &verifiable.W3CCredential{
		Proof: verifiable.CredentialProofs{
			claim.SignatureProof,
			claim.MTP,
		},
	}
	err := json.Unmarshal(claim.Credential, res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal credential")
	}

	return res, nil
}
