package api

import (
	"github.com/lab259/hermes"
	"github.com/lab259/hermes/examples/todos/api/todos"
)

func SetupRoutes(r hermes.Routable) {
	r.Prefix("/todos").Group(func(r hermes.Routable) {
		r.Get("/", todos.Index)
		r.Post("/", todos.Create)
		r.Prefix("/:id").Group(func(r hermes.Routable) {
			r.Get("/", todos.Show)
			r.Put("/", todos.Update)
			r.Delete("/", todos.Delete)
		})
	})
}
