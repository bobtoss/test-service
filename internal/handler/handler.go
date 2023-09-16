package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"test-service/docs"
	"test-service/internal/config"
	"test-service/internal/handler/http"
	"test-service/internal/service"
	"test-service/pkg/server/router"
)

type Dependencies struct {
	Configs     config.Configs
	UserService *service.Service
}

type Configuration func(h *Handler) error

type Handler struct {
	dependencies Dependencies

	HTTP *chi.Mux
}

func New(d Dependencies, configs ...Configuration) (h *Handler, err error) {
	h = &Handler{
		dependencies: d,
	}

	for _, cfg := range configs {
		if err = cfg(h); err != nil {
			return
		}
	}

	return
}

// @title Gin Swagger Example API
// @version 1.0
// @description Testing Swagger APIs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:80
// @BasePath /api/v1
// @schemes http

func WithHTTPHandler() Configuration {
	return func(h *Handler) (err error) {
		h.HTTP = router.New()

		h.HTTP.Use(middleware.Timeout(h.dependencies.Configs.APP.Timeout))

		docs.SwaggerInfo.BasePath = h.dependencies.Configs.APP.Path
		h.HTTP.Get("/swagger/*", httpSwagger.WrapHandler)

		// Init service handlers
		userHandler := http.NewUserHandler(h.dependencies.UserService)

		h.HTTP.Route("/api/v1", func(r chi.Router) {
			r.Mount("/users", userHandler.Routes())
		})

		return
	}
}
