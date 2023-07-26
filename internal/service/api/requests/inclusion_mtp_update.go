package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

const (
	credentialIDPathParam = "credential-id"
)

type inclusionMTPUpdateRequest struct {
	CredentialID string
}

type InclusionMTPUpdateRequest struct {
	CredentialID uuid.UUID
}

func NewInclusionMTPUpdate(r *http.Request) (*InclusionMTPUpdateRequest, error) {
	requestRaw := inclusionMTPUpdateRequest{
		CredentialID: chi.URLParam(r, credentialIDPathParam),
	}

	if err := requestRaw.validate(); err != nil {
		return nil, err
	}

	return requestRaw.parse(), nil
}

func (r *inclusionMTPUpdateRequest) validate() error {
	return validation.Errors{
		"path/{credential-id}": validation.Validate(
			r.CredentialID, validation.Required, validation.By(MustBeValidUUID),
		),
	}.Filter()
}

func (r *inclusionMTPUpdateRequest) parse() *InclusionMTPUpdateRequest {
	return &InclusionMTPUpdateRequest{
		CredentialID: uuid.MustParse(r.CredentialID),
	}
}
