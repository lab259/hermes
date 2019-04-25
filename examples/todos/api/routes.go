package api

import (
	"github.com/lab259/http"
	"github.com/lab259/http/examples/todos/api/todos"
)

func SetupRoutes(r http.Routable) {
	r.Prefix("/todos").Group(func(r http.Routable) {
		r.Get("", todos.Index)
		r.Post("", todos.Create)
		r.Prefix("/:id").Group(func(r http.Routable) {
			r.Get("", todos.Show)
			r.Put("", todos.Update)
			r.Delete("", todos.Delete)
		})
	})
}
