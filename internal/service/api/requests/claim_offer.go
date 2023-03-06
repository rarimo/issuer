package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	core "github.com/iden3/go-iden3-core"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	claimTypePathParam = "claim-type"
	UserIDPathParam    = "user-id"
)

type ClaimOfferRequest struct {
	UserDID   *core.DID
	ClaimType string
}

type claimOfferRequestRaw struct {
	UserID    string
	ClaimType string
}

func NewClaimOffer(r *http.Request) (*ClaimOfferRequest, error) {
	requestBody := claimOfferRequestRaw{
		UserID:    chi.URLParam(r, UserIDPathParam),
		ClaimType: chi.URLParam(r, claimTypePathParam),
	}

	if err := requestBody.validate(); err != nil {
		return nil, err
	}

	return requestBody.parse(), nil
}

// nolint
func (req *claimOfferRequestRaw) validate() error {
	return validation.Errors{
		"path/{claim-type}": validation.Validate(
			req.ClaimType, validation.Required, validation.By(MustBeClaimType),
		),
		"path/{user-id}": validation.Validate(
			req.UserID, validation.Required, validation.By(MustBeValidID),
		),
	}.Filter()
}

func MustBeValidID(src interface{}) error {
	identifierRawBase58, ok := src.(string)
	if !ok {
		return errors.New("it is not an identifier")
	}

	id, err := core.IDFromString(identifierRawBase58)
	if err != nil {
		return errors.New("it is not a string")
	}

	_, err = core.ParseDIDFromID(id)
	if err != nil {
		return errors.New("it is not a valid did")
	}

	return nil
}

func MustBeValidDID(src interface{}) error {
	didRawBase58, ok := src.(string)
	if !ok {
		return errors.New("it is not a string")
	}

	_, err := core.ParseDID(didRawBase58)
	if err != nil {
		return errors.New("it is not a valid did")
	}

	return nil
}

func (req *claimOfferRequestRaw) parse() *ClaimOfferRequest {
	userID := core.ID{}
	_ = userID.UnmarshalText([]byte(req.UserID))
	did, _ := core.ParseDIDFromID(userID)

	return &ClaimOfferRequest{
		UserDID:   did,
		ClaimType: req.ClaimType,
	}
}
