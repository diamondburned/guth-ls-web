package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/diamondburned/listener"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	server := http.Server{
		Addr:    os.Getenv("GUTH_HTTP"),
		Handler: nil,
	}

	if err := listener.HTTPListenAndServeCtx(ctx, &server); err != nil {
		log.Fatalln("cannot listen and serve HTTP:", err)
	}
}
