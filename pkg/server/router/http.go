package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"time"
)

func New() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)

	r.Use(middleware.RealIP)

	r.Use(middleware.Logger)

	r.Use(middleware.Recoverer)

	r.Use(middleware.CleanPath)

	r.Use(middleware.Heartbeat("/"))

	r.Use(middleware.Timeout(time.Second * 60))

	r.Use(middleware.AllowContentType("application/json"))

	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	return r
}
