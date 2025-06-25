package http

import (
	"context"
	"go.uber.org/fx"
	"log"
	"net/http"
)

func RegisterRoutes(lc fx.Lifecycle, gameHandler *GameHandler) {
	mux := http.NewServeMux()

	mux.HandleFunc("/new-game", gameHandler.HandleNewGame)

	mux.HandleFunc("/game/", gameHandler.HandleGameMove)
	mux.Handle("/", http.FileServer(http.Dir("static")))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting HTTP Server on :8080")
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping HTTP Server")
			return server.Shutdown(ctx)
		},
	})
}
