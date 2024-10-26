package main

import (
	"context"
	"github.com/BinaryArchaism/mc-srv/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	ctx, cancel := context.WithCancel(context.Background())
	srv, _ := server.New()
	go func() {
		err := srv.Accept(ctx)
		if err != nil {
			panic(err)
		}
	}()

	log.Info().Msg("server started")
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)
	<-osSignal
	log.Info().Msg("server shutdown signal received")
	cancel()
	time.Sleep(1 * time.Second)
	log.Info().Msg("server shutdown")
}
