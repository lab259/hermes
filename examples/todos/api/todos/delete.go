package todos

import (
	"github.com/lab259/http"
	"github.com/lab259/http/examples/todos/errors"
)

func Delete(req http.Request, res http.Response) http.Result {
	id := req.Param("id")
	if _, found := db[id]; !found {
		return res.Status(404).Error(errors.ErrNotFound)
	}
	delete(db, id)
	return res.Status(204).Data("")
}
