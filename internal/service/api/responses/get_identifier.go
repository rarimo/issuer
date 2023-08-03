package responses

import (
	"github.com/rarimo/issuer/resources"
)

func NewGetIdentifier(identifier string) *resources.GetIdentifierResponse {
	return &resources.GetIdentifierResponse{
		Data: resources.GetIdentifier{
			Key: resources.Key{
				Type: resources.IDENTIFIER,
			},
			Attributes: &resources.GetIdentifierAttributes{
				Identifier: identifier,
			},
		},
	}
}
