package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/q-dev/q-id/issuer/internal/service/api/responses"
)

func PublishState(w http.ResponseWriter, r *http.Request) {
	txHash, err := Issuer(r).PublishStateOnChain(r.Context())
	if err != nil {
		Log(r).WithError(err).Error("Failed to publish issuer state on chain")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if txHash == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	ape.Render(w, responses.NewPublishState(txHash))
}
