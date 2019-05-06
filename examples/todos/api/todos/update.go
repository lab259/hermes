package todos

import (
	"github.com/lab259/http"
	"github.com/lab259/http/examples/todos/errors"
)

func Update(req http.Request, res http.Response) http.Result {
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
