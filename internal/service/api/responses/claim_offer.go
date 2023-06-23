package responses

import (
	"github.com/iden3/iden3comm/protocol"
	"gitlab.com/rarimo/identity/issuer/resources"
)

func NewClaimOffer(claimOffer *protocol.CredentialsOfferMessage) *resources.ClaimOfferResponse {
	credentials := make([]resources.ClaimOfferBodyCredentials, 0, len(claimOffer.Body.Credentials))
	for _, credential := range claimOffer.Body.Credentials {
		credentials = append(credentials, resources.ClaimOfferBodyCredentials{
			Description: credential.Description,
			Id:          credential.ID,
		})
	}

	return &resources.ClaimOfferResponse{
		Data: resources.ClaimOffer{
			Key: resources.Key{
				Type: resources.CLAIM_OFFER,
			},
			Attributes: resources.ClaimOfferAttributes{
				Body: resources.ClaimOfferBody{
					Credentials: credentials,
					Url:         claimOffer.Body.URL,
				},
				From: claimOffer.From,
				Id:   claimOffer.ID,
				Thid: claimOffer.ThreadID,
				To:   claimOffer.To,
				Typ:  string(claimOffer.Typ),
				Type: string(claimOffer.Type),
			},
		},
		Included: resources.Included{},
	}
}
