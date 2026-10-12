package app

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) requestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		route := chi.RouteContext(r.Context()).RoutePattern()
		if route == "" {
			route = "unmatched"
		}

		a.l.Info(
			"Http Request",
			"req_id", middleware.GetReqID(r.Context()),
			"method", r.Method,
			"route", route,
			"path", r.URL.Path,
			"status", ww.Status(),
			"duration", time.Since(startTime).Milliseconds(),
			"bytes_written", ww.BytesWritten(),
		)
	})
}
