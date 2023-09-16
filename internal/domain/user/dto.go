package user

import (
	"errors"
	"net/http"
)

type Request struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func (s *Request) Bind(r *http.Request) error {
	if s.Name == "" {
		return errors.New("name: cannot be blank")
	}

	if s.Email == "" {
		return errors.New("email: cannot be blank")
	}

	if s.Status != "online" && s.Status != "offline" && s.Status != "delete" {
		return errors.New("status: incorrect")
	}

	return nil
}

type UpdateResponse struct {
	ID       string `json:"id"`
	New      string `json:"new"`
	Previous string `json:"previous"`
}

type Response struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:     data.ID.Hex(),
		Name:   *data.Name,
		Email:  *data.Email,
		Status: *data.Status,
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
