package app

import (
	"net/http"
	"os"
	"sandhi/internal/db"
	"sandhi/internal/logging"
	"slices"

	handlerV1 "sandhi/internal/app/handler/v1"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	db *pgxpool.Pool
	q  *db.Queries
	l  *logging.Logger
}

func NewApp(db *pgxpool.Pool, q *db.Queries, l *logging.Logger) *App {
	return &App{
		db: db,
		q:  q,
		l:  l,
	}
}

func (a *App) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(a.requestMiddleware)
	r.Use(middleware.Recoverer)

	AllowedOrigins := []string{
		"http://localhost:3000",  // optional (web dev)
		"http://localhost:8000",  // optional (web dev)
		"http://localhost:5173",  // optional (web dev)
		"https://localhost",      // Android Capacitor
		"capacitor://localhost",  // iOS Capacitor
		"http://localhost:19006", // Expo dev
	}
	// Apply CORS globally to all handlers
	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc: func(_ *http.Request, origin string) bool {
			// 1. In DEV, allow everything (including your Mac's IP)
			// while still supporting AllowCredentials: true
			if os.Getenv("APP_ENV") == "DEV" {
				return true
			}

			// 2. In PROD, match against your strict list
			return slices.Contains(AllowedOrigins, origin)
		},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "ngrok-skip-browser-warning"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server is up and running"))
	})

	r.Get("/dbHealth", func(w http.ResponseWriter, r *http.Request) {
		if err := a.db.Ping(r.Context()); err != nil {
			RespondWithError(w, http.StatusServiceUnavailable, "database unavailable")
			return
		}
		RespondWithJSON(w, 200, map[string]string{
			"msg": "server is up and running and database in alive",
		})
	})

	handlerV1 := handlerV1.NewHandler(a.db, a.q, a.l)

	r.Mount("/api/v1", handlerV1.Routes())

	return r
}
