package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/diamondburned/guth-ls-web/internal/frontend/root"
	"github.com/diamondburned/guth-ls-web/internal/guthls"
	"github.com/diamondburned/listener"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalln("cannot load .env:", err)
	}

	var prov guthls.Provider
	var server http.Server

	parseEnv(map[string]func(string) error{
		"GUTH_MYSQL": func(val string) (err error) {
			prov, err = guthls.NewMySQLProvider(val)
			return
		},
		"GUTH_HTTP": func(val string) (err error) {
			server.Addr = val
			return nil
		},
	})

	if prov == nil {
		log.Fatalln("no provider given, see .env")
	}

	server.Handler = root.Mount(prov)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Println("Listening at", server.Addr)

	if err := listener.HTTPListenAndServeCtx(ctx, &server); err != nil {
		log.Fatalln("cannot listen and serve HTTP:", err)
	}
}

func parseEnv(envs map[string]func(string) error) {
	for env, fn := range envs {
		// TODO: optional, whatever, no need for now
		val := os.Getenv(env)
		if val == "" {
			log.Fatalf("missing $%s", env)
		}
		if err := fn(val); err != nil {
			log.Fatalf("$%s error: %v", env, err)
		}
	}
}
