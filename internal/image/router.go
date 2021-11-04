package image

import "github.com/go-chi/chi/v5"

// Routes add new routes to chi Router.
func (s *Service) Routes(r chi.Router) {
	r.Route("/image", func(r chi.Router) {
		r.Post("/", s.Upload)
		r.Get("/", s.GetImage)
	})
}