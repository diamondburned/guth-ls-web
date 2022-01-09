package frontend

import (
	"embed"

	"github.com/diamondburned/tmplutil"
)

//go:embed *
var webFS embed.FS

var Templater = tmplutil.Templater{
	FileSystem: webFS,
	Includes: map[string]string{
		"css":    "components/css.html",
		"header": "components/header.html",
		"footer": "components/footer.html",
	},
}
