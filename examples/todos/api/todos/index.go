package todos

import "github.com/lab259/hermes"

func Index(req hermes.Request, res hermes.Response) hermes.Result {
	list := make([]*Todo, len(db))
	idx := 0
	for _, todo := range db {
		list[idx] = &todo
		idx++
	}
	return res.Data(list)
}
