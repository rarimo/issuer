package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"gitlab.com/rarimo/identity/issuer/internal/service/api/requests"
	"gitlab.com/rarimo/identity/issuer/internal/service/core/issuer"
)

func OfferCallback(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewOfferCallback(r)
	if err != nil {
		Log(r).WithField("reason", err).Debug("Bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	response, err := Issuer(r).OfferCallback(r.Context(), req)
	switch {
	case errors.Is(err, issuer.ErrClaimRetrieverIsNotClaimOwner),
		errors.Is(err, issuer.ErrRepeatedCallbackRequest),
		errors.Is(err, issuer.ErrProofVerifyFailed):
		Log(r).WithField("reason", err).Debug("Forbidden")
		ape.RenderErr(w, problems.Forbidden())
		return
	case errors.Is(err, issuer.ErrMessageRecipientIsNotIssuer):
		Log(r).WithField("reason", err).Debug("Bad request")
		ape.RenderErr(w, problems.BadRequest(errors.Cause(err))...)
		return
	case errors.Is(err, issuer.ErrClaimOfferIsNotExist), errors.Is(err, issuer.ErrClaimIsNotExist):
		Log(r).WithField("reason", err).Debug("Not found")
		ape.RenderErr(w, problems.NotFound())
		return
	case err != nil:
		Log(r).WithError(err).Error("Failed to issue claim")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	marshaled, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(marshaled))

	ape.Render(w, response)
}
