package handlers

import (
	"net/http"

	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/q-dev/q-id/issuer/internal/service/api/requests"
	"gitlab.com/q-dev/q-id/issuer/internal/service/api/responses"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims/schemas"
)

func IssueClaim(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewIssueClaim(r)
	if err != nil {
		Log(r).WithField("reason", err).Debug("Bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	claimID, err := Issuer(r).IssueClaim(r.Context(), req.UserDID, &req.Expiration, req.ClaimType, req.Credential)
	switch {
	case errors.Is(err, schemas.ErrValidationData):
		Log(r).WithField("reason", err).Debug("Bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	case errors.Is(err, merkletree.ErrEntryIndexAlreadyExists):
		Log(r).WithField("reason", err).Debug("Conflict")
		ape.RenderErr(w, problems.Conflict())
		return
	case err != nil:
		Log(r).WithError(err).
			WithField("schema-type", req.ClaimType).
			WithField("schema-data", string(req.Credential)).
			WithField("user-id", req.UserDID.String()).
			Error("Failed to issue claim")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, responses.NewIssueClaim(claimID))
}
