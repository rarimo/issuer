package responses

import (
	"fmt"

	"gitlab.com/q-dev/q-id/issuer/resources"
)

type IssueClaimResponse struct {
	Data IssueClaimResponseData `json:"data"`
}

type IssueClaimResponseData struct {
	resources.Key
}

func NewIssueClaim(claimID uint64) *IssueClaimResponse {
	return &IssueClaimResponse{
		Data: IssueClaimResponseData{
			Key: resources.Key{
				ID:   fmt.Sprint(claimID),
				Type: resources.CLAIM_ID,
			},
		},
	}
}
