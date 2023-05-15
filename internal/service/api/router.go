package api

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"

	"gitlab.com/q-dev/q-id/issuer/internal/service/api/handlers"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.ContentType("application/vnd.api+json"),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxIssuer(s.issuer),
		),
	)

	r.Route("/integrations/issuer", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/public", func(r chi.Router) {
				r.Route("/claims", func(r chi.Router) {
					r.Route("/offers", func(r chi.Router) {
						r.Get("/{user-id}/{claim-type}", handlers.ClaimOffer)
						r.Post("/callback", handlers.OfferCallback)
					})

					r.Route("/revocations", func(r chi.Router) {
						r.Get("/check/{rev-id}", handlers.RevocationCheck)
					})
				})

				r.Route("/identity", func(r chi.Router) {
					r.Get("/identifier", handlers.GetIdentifier)
				})
			})

			r.Route("/private", func(r chi.Router) {
				r.Route("/claims", func(r chi.Router) {
					r.Route("/revocations", func(r chi.Router) {
						r.Post("/{user-id}/{claim-type}", handlers.ClaimRevocation)
					})

					r.Post("/issue/{user-id}/{claim-type}", handlers.IssueClaim)
				})
			})
		})
	})

	return r
}
