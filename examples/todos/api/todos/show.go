package todos

import (
	"github.com/lab259/http"
	"github.com/lab259/http/examples/todos/errors"
)

func Show(req http.Request, res http.Response) http.Result {
	id := req.Param("id")
	todo, found := db[id]
	if !found {
		return res.Status(404).Error(errors.ErrNotFound)
	}
	return res.Data(&todo)
}
