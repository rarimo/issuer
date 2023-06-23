package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"

	"gitlab.com/rarimo/identity/issuer/internal/service/api/responses"
)

func GetIdentifier(w http.ResponseWriter, r *http.Request) {
	ape.Render(w, responses.NewGetIdentifier(Issuer(r).GetIdentifier()))
}
