package root

import (
	"net/http"

	"github.com/diamondburned/guth-ls-web/internal/frontend"
	"github.com/diamondburned/guth-ls-web/internal/guthls"
	"github.com/diamondburned/tmplutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Mount returns a new http.Handler containing routes.
func Mount(prov guthls.Provider) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.CleanPath)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Compress(5))

	r.Handle("/static*", frontend.MountStatic())

	r.Group(func(r chi.Router) {
		r.Use(tmplutil.AlwaysFlush)
		r.Use(middleware.NoCache)
		r.Use(contentType("text/html"))
		r.Handle("/*", newIndex(prov))
	})

	return r
}

// AlwaysFlush is the middleware to always flush after a write.
func contentType(typ string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", typ)
			next.ServeHTTP(w, r)
		})
	}
}
