package handler

import (
	"net/http"
	"sandhi/internal/db"
	"sandhi/internal/logging"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db *pgxpool.Pool
	q  *db.Queries
	l  *logging.Logger
}

func NewHandler(db *pgxpool.Pool, q *db.Queries, l *logging.Logger) *Handler {
	h := &Handler{
		db: db,
		q:  q,
		l:  l,
	}
	return h
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/api/v1/health is up and running"))
	})

	// define unprotected routes

	// define protected routes
	r.Group(func(r chi.Router) {
		// use a authentiaction middleware
		// define protected routes
	})

	return r
}
