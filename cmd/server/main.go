package main

import (
	"context"
	"github.com/BinaryArchaism/mc-srv/internal/server"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	srv, _ := server.New()
	go func() {
		err := srv.Accept(ctx)
		if err != nil {
			panic(err)
		}
	}()

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)
	<-osSignal
	cancel()
	time.Sleep(1 * time.Second)
}
