package task

import "time"

type Request struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	ActiveAt string `json:"active_at"`
}

type Response struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	ActiveAt time.Time `json:"active_at"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:       data.ID,
		Title:    *data.Title,
		ActiveAt: *data.ActiveAt,
	}
	return
}

func ParseFromEntities(data []Entity) (res []Response) {
	res = make([]Response, 0)
	for _, object := range data {
		res = append(res, ParseFromEntity(object))
	}
	return
}
