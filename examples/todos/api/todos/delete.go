package todos

import (
	"github.com/lab259/hermes"
	"github.com/lab259/hermes/examples/todos/errors"
)

func Delete(req hermes.Request, res hermes.Response) hermes.Result {
	id := req.Param("id")
	if _, found := db[id]; !found {
		return res.Status(404).Error(errors.ErrTodoNotFound)
	}
	delete(db, id)
	return res.Status(204).Data("")
}
