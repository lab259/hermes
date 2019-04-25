package todos

import (
	"fmt"
	"time"

	"github.com/lab259/http"
	"github.com/lab259/http/examples/todos/errors"
)

func Create(req http.Request, res http.Response) http.Result {
	now := time.Now()

	var todo Todo
	if err := req.Data(&todo); err != nil {
		return res.Status(400).Error(err)
	}

	if todo.Description == "" {
		return res.Status(400).Error(errors.ErrDescriptionRequired)
	}

	todo.ID = fmt.Sprintf("%d", now.Unix())
	db[todo.ID] = todo

	return res.Status(201).Data(&todo)
}
