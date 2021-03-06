package frontend

import (
	"embed"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/diamondburned/guth-ls-web/internal/duration"
	"github.com/diamondburned/tmplutil"
	"github.com/dustin/go-humanize"
	"github.com/go-chi/chi/v5"
)

//go:embed *
var webFS embed.FS

var Templater = tmplutil.Templater{
	FileSystem: tmplutil.OverrideFS(webFS, os.DirFS("frontend")),
	Includes: map[string]string{
		"css":    "components/css.html",
		"error":  "components/error.html",
		"header": "components/header.html",
		"footer": "components/footer.html",
	},
	Functions: template.FuncMap{
		"RFC3339": func() string {
			return time.RFC3339
		},
		"RelTime": func(t time.Time) string {
			return humanize.Time(t)
		},
		"RelDurationShort": func(d time.Duration) string {
			return duration.Short(d)
		},
		"RelDurationLong": func(d time.Duration) string {
			return duration.Long(d)
		},
	},
	OnRenderFail: func(sub *tmplutil.Subtemplate, w io.Writer, err error) {
		sub.Templater().Execute(w, "error", err)
	},
}

// MountStatic mounts the /static directory. Note that the returned handler will
// only handle /static paths, so it's safe to mount this to /static directly.
func MountStatic() http.Handler {
	mux := chi.NewMux()
	mux.Use(staticCache)
	mux.Handle("/static*", http.FileServer(http.FS(webFS)))
	return mux
}

func staticCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// https://cache-control.sdgluck.vercel.app/
		w.Header().Set(
			"Cache-Control",
			"public, max-age 86400, max-stale 604800, stale-while-revalidate 604800",
		)
		next.ServeHTTP(w, r)
	})
}
