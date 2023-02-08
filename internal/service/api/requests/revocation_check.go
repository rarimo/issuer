package requests

import (
	"math/big"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/iden3/iden3comm/protocol"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	revocationIDPathParam = "rev-id"
)

type revocationCheckRequest struct {
	RevocationID      string
	RevocationRequest *protocol.RevocationStatusRequestMessage
}

type RevocationCheckRequest struct {
	RevocationID      *big.Int
	RevocationRequest *protocol.RevocationStatusRequestMessage
}

func NewRevocationCheck(r *http.Request) (*RevocationCheckRequest, error) {
	requestRaw := revocationCheckRequest{
		RevocationID: chi.URLParam(r, revocationIDPathParam),
	}

	if err := requestRaw.validate(); err != nil {
		return nil, err
	}

	return requestRaw.parse(), nil
}

// nolint
func (r *revocationCheckRequest) validate() error {
	return validation.Errors{
		"path/{revocation-id}": validation.Validate(
			r.RevocationID, validation.Required, validation.By(MustBeCorrectRevocationID),
		),
	}.Filter()
}

func (r *revocationCheckRequest) parse() *RevocationCheckRequest {
	revID, _ := new(big.Int).SetString(r.RevocationID, 10) //nolint

	return &RevocationCheckRequest{
		RevocationID: revID,
	}
}

func MustBeCorrectRevocationID(src interface{}) error {
	revocationIDRaw, ok := src.(string)
	if !ok {
		return errors.New("it is not a revocation id")
	}

	_, err := strconv.ParseUint(revocationIDRaw, 10, 64)
	if err != nil {
		return errors.New("it is not valid uint64")
	}

	return nil
}
