package responses

import (
	"github.com/rarimo/issuer/resources"
)

type IssueClaimResponse struct {
	Data IssueClaimResponseData `json:"data"`
}

type IssueClaimResponseData struct {
	resources.Key
}

func NewIssueClaim(claimID string) *IssueClaimResponse {
	return &IssueClaimResponse{
		Data: IssueClaimResponseData{
			Key: resources.Key{
				ID:   claimID,
				Type: resources.CLAIM_ID,
			},
		},
	}
}
