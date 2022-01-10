package root

import (
	"net/http"

	"github.com/diamondburned/guth-ls-web/internal/frontend"
	"github.com/diamondburned/guth-ls-web/internal/guthls"
	"github.com/go-chi/chi/v5"
)

var indexTmpl = frontend.Templater.Register("index", "root/index.html")

func init() {
	flags := map[string]guthls.LeaderboardQueryFlags{
		"QueryUser": guthls.LeaderboardQueryUser,
		"QueryRank": guthls.LeaderboardQueryRank,
	}

	for name, flag := range flags {
		flag := flag // copy for closure
		frontend.Templater.Func(name, func() guthls.LeaderboardQueryFlags { return flag })
	}
}

type index struct {
	*chi.Mux
	Prov guthls.Provider
}

func newIndex(prov guthls.Provider) http.Handler {
	index := index{
		Mux:  chi.NewMux(),
		Prov: prov,
	}
	index.Get("/", index.render)
	return index
}

func (index index) render(w http.ResponseWriter, r *http.Request) {
	indexTmpl.Execute(w, index)
}
