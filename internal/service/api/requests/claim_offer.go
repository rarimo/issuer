package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	core "github.com/iden3/go-iden3-core"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	claimIDPathParam = "claim-id"
	UserIDPathParam  = "user-id"
)

type ClaimOfferRequest struct {
	UserID  *core.ID
	ClaimID string
}

type claimOfferRequestRaw struct {
	UserID  string
	ClaimID string
}

func NewClaimOffer(r *http.Request) (*ClaimOfferRequest, error) {
	requestBody := claimOfferRequestRaw{
		UserID:  chi.URLParam(r, UserIDPathParam),
		ClaimID: chi.URLParam(r, claimIDPathParam),
	}

	if err := requestBody.validate(); err != nil {
		return nil, err
	}

	return requestBody.parse(), nil
}

// nolint
func (req *claimOfferRequestRaw) validate() error {
	return validation.Errors{
		"path/{claim-id}": validation.Validate(
			req.ClaimID, validation.Required, validation.By(MustBeClaimID),
		),
		"path/{user-id}": validation.Validate(
			req.UserID, validation.Required, validation.By(MustBeIden3Identifier),
		),
	}.Filter()
}

func MustBeIden3Identifier(src interface{}) error {
	identifierRawBase58, ok := src.(string)
	if !ok {
		return errors.New("it is not an identifier")
	}

	_, err := core.IDFromString(identifierRawBase58)
	if err != nil {
		return errors.New("it is not an identifier")
	}

	return nil
}

func (req *claimOfferRequestRaw) parse() *ClaimOfferRequest {
	userID := &core.ID{}
	_ = userID.UnmarshalText([]byte(req.UserID))

	return &ClaimOfferRequest{
		UserID:  userID,
		ClaimID: req.ClaimID,
	}
}
