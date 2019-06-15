package todos

import (
	"github.com/lab259/hermes"
	"github.com/lab259/hermes/examples/todos/errors"
)

func Show(req hermes.Request, res hermes.Response) hermes.Result {
	id := req.Param("id")
	todo, found := db[id]
	if !found {
		return res.Status(404).Error(errors.ErrTodoNotFound)
	}
	return res.Data(&todo)
}
