package todos

var db = make(map[string]Todo)

type Todo struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}
