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

	r.Route("/integrations/qid-issuer", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/state", func(r chi.Router) {
				r.Post("/publish", handlers.PublishState)
			})

			r.Route("/claims", func(r chi.Router) {
				r.Route("/offers", func(r chi.Router) {
					r.Get("/{user-id}/{claim-id}", handlers.ClaimOffer)
					r.Post("/callback", handlers.OfferCallback)
				})

				r.Post("/issue/{user-id}", handlers.IssueClaim)
			})
		})
	})

	return r
}
