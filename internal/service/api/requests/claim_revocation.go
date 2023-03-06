package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	core "github.com/iden3/go-iden3-core"
	claimResources "gitlab.com/q-dev/q-id/resources/claim_resources"
)

type ClaimRevocationRequest struct {
	UserID  *core.ID
	ClaimID claimResources.ClaimSchemaType
}

type claimRevocationRequestRaw struct {
	UserID  string
	ClaimID string
}

func NewRevocationClaim(r *http.Request) (*ClaimRevocationRequest, error) {
	requestRaw := claimRevocationRequestRaw{
		UserID:  chi.URLParam(r, UserIDPathParam),
		ClaimID: chi.URLParam(r, claimTypePathParam),
	}

	if err := requestRaw.validate(); err != nil {
		return nil, err
	}

	return requestRaw.parse(), nil
}

// nolint
func (req *claimRevocationRequestRaw) validate() error {
	return validation.Errors{
		"path/{user-id}": validation.Validate(
			req.UserID, validation.Required, validation.By(MustBeValidID),
		),
		"path/{claim-id}": validation.Validate(
			req.ClaimID, validation.Required, validation.By(MustBeClaimID),
		),
	}.Filter()
}

func (req *claimRevocationRequestRaw) parse() *ClaimRevocationRequest {
	userID := &core.ID{}

	_ = userID.UnmarshalText([]byte(req.UserID))

	return &ClaimRevocationRequest{
		UserID:  userID,
		ClaimID: claimResources.ClaimSchemaTypeList[req.ClaimID],
	}
}
