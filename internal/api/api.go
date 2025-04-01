package api

import (
	auth_router "github.com/Antonious-Stewart/Full-Scale-Ecommerce/internal/api/routes/authRouter"
	"github.com/Antonious-Stewart/Full-Scale-Ecommerce/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Handler struct {
	DbPool    db.DB
	Validator *validator.Validate
}

func New(dbPool db.DB, validate *validator.Validate) *Handler {

	return &Handler{
		DbPool:    dbPool,
		Validator: validate,
	}
}

func (h *Handler) Routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/api", func(childRoute chi.Router) {
		childRoute.Mount("/auth", auth_router.New(h.DbPool, h.Validator).Routes())
	})

	return router
}
