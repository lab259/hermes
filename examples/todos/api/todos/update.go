package todos

import (
	"github.com/lab259/hermes"
	"github.com/lab259/hermes/examples/todos/errors"
)

func Update(req hermes.Request, res hermes.Response) hermes.Result {
	id := req.Param("id")
	_, found := db[id]
	if !found {
		return res.Status(404).Error(errors.ErrTodoNotFound)
	}

	var todo Todo
	if err := req.Data(&todo); err != nil {
		return res.Status(400).Error(err)
	}

	if todo.Description == "" {
		return res.Status(400).Error(errors.ErrDescriptionRequired)
	}

	todo.ID = id
	db[id] = todo

	return res.Data(&todo)
}
