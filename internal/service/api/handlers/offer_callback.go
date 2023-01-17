package handlers

import (
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/q-dev/q-id/issuer/internal/service/api/requests"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/issuer"
)

func OfferCallback(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewOfferCallback(r)
	if err != nil {
		Log(r).WithField("reason", err).Debug("Bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	response, err := Issuer(r).OfferCallback(r.Context(), req.FetchMessage)
	switch {
	case errors.Is(err, issuer.ErrClaimRetrieverIsNotClaimOwner):
		Log(r).WithField("reason", err).Error("Forbidden")
		ape.RenderErr(w, problems.Forbidden())
		return
	case errors.Is(err, issuer.ErrMessageRecipientIsNotIssuer):
		Log(r).WithField("reason", err).Error("Bad request")
		ape.RenderErr(w, problems.BadRequest(errors.Cause(err))...)
		return
	case err != nil:
		Log(r).WithError(err).Error("Failed to issue claim")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, response)
}
