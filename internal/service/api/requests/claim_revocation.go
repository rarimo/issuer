package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	core "github.com/iden3/go-iden3-core"
	claimResources "gitlab.com/q-dev/q-id/resources/claim_resources"
)

type ClaimRevocationRequest struct {
	UserID    *core.ID
	ClaimType claimResources.ClaimSchemaType
}

type claimRevocationRequestRaw struct {
	UserID    string
	ClaimType string
}

func NewRevocationClaim(r *http.Request) (*ClaimRevocationRequest, error) {
	requestRaw := claimRevocationRequestRaw{
		UserID:    chi.URLParam(r, UserIDPathParam),
		ClaimType: chi.URLParam(r, claimTypePathParam),
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
		"path/{claim-type}": validation.Validate(
			req.ClaimType, validation.Required, validation.By(MustBeClaimType),
		),
	}.Filter()
}

func (req *claimRevocationRequestRaw) parse() *ClaimRevocationRequest {
	userID := &core.ID{}

	_ = userID.UnmarshalText([]byte(req.UserID))

	return &ClaimRevocationRequest{
		UserID:    userID,
		ClaimType: claimResources.ClaimSchemaTypeList[req.ClaimType],
	}
}
