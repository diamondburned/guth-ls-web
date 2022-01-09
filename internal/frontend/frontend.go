package frontend

import (
	"embed"
	"os"

	"github.com/diamondburned/tmplutil"
)

//go:embed *
var webFS embed.FS

var Templater = tmplutil.Templater{
	FileSystem: tmplutil.OverrideFS(webFS, os.DirFS(".")),
	Includes: map[string]string{
		"css":    "components/css.html",
		"header": "components/header.html",
		"footer": "components/footer.html",
	},
}
