module github.com/diamondburned/guth-ls-web

go 1.17

replace github.com/diamondburned/tmplutil => ../tmplutil

require (
	github.com/diamondburned/listener v0.0.0-20201025235325-c458deda1495
	github.com/diamondburned/tmplutil v0.0.0-20220109065244-7ae8b7fad5dd
	github.com/go-sql-driver/mysql v1.6.0
)

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/go-chi/chi/v5 v5.0.7
	github.com/joho/godotenv v1.4.0
	github.com/pkg/errors v0.9.1
)

require github.com/leighmacdonald/steamid v1.2.0

require github.com/keegancsmith/sqlf v1.1.1
