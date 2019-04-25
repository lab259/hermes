package api

import (
	"github.com/lab259/http"
	"github.com/lab259/http/examples/todos/api/todos"
)

func SetupRoutes(r http.Routable) {
	r.Prefix("/todos").Group(func(r http.Routable) {
		r.GET("", todos.Index)
		r.POST("", todos.Create)
		r.Prefix("/:id").Group(func(r http.Routable) {
			r.GET("", todos.Show)
			r.PUT("", todos.Update)
			r.DELETE("", todos.Delete)
		})
	})
}
