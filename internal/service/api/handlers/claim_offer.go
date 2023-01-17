package handlers

import (
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/q-dev/q-id/issuer/internal/service/api/requests"
	"gitlab.com/q-dev/q-id/issuer/internal/service/api/responses"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/issuer"
)

func ClaimOffer(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewClaimOffer(r)
	if err != nil {
		Log(r).WithField("reason", err).Debug("Bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	claimOffer, err := Issuer(r).CreateClaimOffer(req.UserID, req.ClaimID)
	switch {
	case errors.Is(err, issuer.ErrClaimIsNotExist):
		Log(r).WithField("reason", err).Debug("Not found")
		ape.RenderErr(w, problems.NotFound())
		return
	case errors.Is(err, issuer.ErrClaimRetrieverIsNotClaimOwner):
		Log(r).WithField("reason", err).Debug("Forbidden")
		ape.RenderErr(w, problems.Forbidden())
		return
	case err != nil:
		Log(r).WithError(err).
			WithField("claim-id", req.ClaimID).
			WithField("user-id", req.UserID).
			Error("Failed get claim offer")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, responses.NewClaimOffer(claimOffer))
}
