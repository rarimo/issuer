package responses

import "gitlab.com/rarimo/identity/issuer/resources"

func NewPublishState(txHash string) *resources.PublishStateResponse {
	return &resources.PublishStateResponse{
		Data: resources.PublishState{
			Key: resources.Key{
				Type: resources.PUBLISH_STATE,
			},
			Attributes: resources.PublishStateAttributes{
				TxHash: txHash,
			},
		},
	}
}
