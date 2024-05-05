package atlas

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"net/http"
)

func (a *Atlas) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)

	if a.Debug {
		mux.Use(middleware.Logger)
	}
	mux.Use(middleware.Recoverer)

	mux.Use(a.SessionLoad)

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to %s", a.AppName)
	})

	return mux
}
