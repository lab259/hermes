package todos

import "github.com/lab259/http"

func Index(req http.Request, res http.Response) http.Result {
	list := make([]*Todo, len(db))
	idx := 0
	for _, todo := range db {
		list[idx] = &todo
		idx++
	}
	return res.Data(list)
}
