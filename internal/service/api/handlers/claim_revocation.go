package handlers

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/rarimo/issuer/internal/service/api/requests"
	"github.com/rarimo/issuer/internal/service/core/issuer"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func ClaimRevocation(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewRevocationClaim(r)
	if err != nil {
		Log(r).WithField("reason", err).Debug("Bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	err = Issuer(r).RevokeClaim(r.Context(), req.UserID, req.ClaimType)
	switch {
	case errors.Is(err, issuer.ErrClaimIsNotExist):
		Log(r).WithField("reason", err).Debug("Not found")
		ape.RenderErr(w, problems.NotFound())
		return
	case errors.Is(err, issuer.ErrClaimIsAlreadyRevoked):
		Log(r).WithField("reason", err).Debug("Conflict")
		ape.RenderErr(w, problems.Conflict())
		return
	case err != nil:
		Log(r).WithError(err).
			WithField("claim-id", req.ClaimType).
			WithField("user-id", req.UserID).
			Error("Failed get claim offer")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
