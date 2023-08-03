package handlers

import (
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"github.com/rarimo/issuer/internal/service/api/requests"
	"github.com/rarimo/issuer/internal/service/core/issuer"
)

func RevocationCheck(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewRevocationCheck(r)
	if err != nil {
		Log(r).WithField("reason", err).Debug("Bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	response, err := Issuer(r).GetRevocationStatus(r.Context(), req.RevocationID)
	if err != nil {
		if errors.Is(err, issuer.ErrMessageRecipientIsNotIssuer) {
			Log(r).WithField("reason", err).Debug("Bad request")
			ape.RenderErr(w, problems.BadRequest(errors.Cause(err))...)
			return
		}
		Log(r).WithError(err).Error("Failed to get revocation status")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, response)
}
