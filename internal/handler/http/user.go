package http

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"test-service/internal/domain/user"
	"test-service/internal/service"
	"test-service/pkg/server/status"
	"test-service/pkg/store"
)

type UserHandler struct {
	userService *service.Service
}

func NewUserHandler(s *service.Service) *UserHandler {
	return &UserHandler{userService: s}
}

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.add)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update)
	})

	return r
}

// @Summary	list of users from the repository
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		email		query	string	false	"query param"
// @Param		name		query	string 	false	"query param"
// @Param		status		query	string	false	"query param"
// @Success	200			{array}		user.Response
// @Failure	500			{object}	status.Object
// @Router		/users 	[get]
func (h *UserHandler) list(w http.ResponseWriter, r *http.Request) {
	res, err := h.userService.ListUsers(r.Context(), r)
	if err != nil {
		status.InternalServerError(w, r, err)
		return
	}

	status.OK(w, r, res)
}

// @Summary	add a new user to the repository
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		request	body		user.Request	true	"body param"
// @Success	200		{object}	user.Response
// @Failure	400		{object}	status.Object
// @Failure	500		{object}	status.Object
// @Router		/users [post]
func (h *UserHandler) add(w http.ResponseWriter, r *http.Request) {
	req := user.Request{}
	if err := render.Bind(r, &req); err != nil {
		status.BadRequest(w, r, err, req)
		return
	}
	res, err := h.userService.AddUser(r.Context(), req)
	if err != nil {
		status.InternalServerError(w, r, err)
		return
	}

	status.OK(w, r, res)
}

// @Summary	get the user from the repository
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		id	path		string	true	"path param"
// @Success	200	{object}	user.Response
// @Failure	404	{object}	status.Object
// @Failure	500	{object}	status.Object
// @Router		/users/{id} [get]
func (h *UserHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		status.BadRequest(w, r, err, "")
		return
	}

	res, err := h.userService.GetUser(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			status.NotFound(w, r, err)
		default:
			status.InternalServerError(w, r, err)
		}
		return
	}

	status.OK(w, r, res)
}

// @Summary	update the user in the repository
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		id		path	string				true	"path param"
// @Param		request	body	user.Request	true	"body param"
// @Success	200
// @Failure	400	{object}	status.Object
// @Failure	404	{object}	status.Object
// @Failure	500	{object}	status.Object
// @Router		/users/{id} [put]
func (h *UserHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		status.BadRequest(w, r, err, "")
		return
	}

	req := user.Request{}
	if err := render.Bind(r, &req); err != nil {
		status.BadRequest(w, r, err, req)
		return
	}

	res, err := h.userService.UpdateStatus(r.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			status.NotFound(w, r, err)
		default:
			status.InternalServerError(w, r, err)
		}
		return
	}
	status.OK(w, r, res)
}
