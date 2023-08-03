package handlers

import (
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"github.com/rarimo/issuer/internal/service/api/requests"
	"github.com/rarimo/issuer/internal/service/core/issuer"
)

func InclusionMTPUpdate(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewInclusionMTPUpdate(r)
	if err != nil {
		Log(r).WithField("reason", err).Debug("Bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	response, err := Issuer(r).GetInclusionMTP(r.Context(), req.CredentialID)
	if err != nil {
		if errors.Is(err, issuer.ErrClaimIsNotExist) {
			Log(r).WithField("reason", err).Debug("Not found")
			ape.RenderErr(w, problems.NotFound())
			return
		}
		Log(r).WithError(err).Error("Failed to get inclusion mtp proof")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, response)
}
