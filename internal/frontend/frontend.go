package frontend

import (
	"embed"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/diamondburned/tmplutil"
	"github.com/dustin/go-humanize"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed *
var webFS embed.FS

var Templater = tmplutil.Templater{
	FileSystem: tmplutil.OverrideFS(webFS, os.DirFS(".")),
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
		"RelDuration": func(d time.Duration) string {
			now := time.Unix(0, 0)
			return humanize.RelTime(now.Add(d), now, "", "")
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
	mux.Use(middleware.Compress(5))
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
