package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"homework9/internal/adapters/repo"
	"homework9/internal/app"
	grpc "homework9/internal/ports/grpcPort"
	"homework9/internal/ports/httpgin"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	grpcPort = ":50054"
	httpPort = ":9000"
)

func main() {

	a := app.NewApp(repo.NewMapAdRepo(), repo.NewUserRepo())

	eg, ctx := errgroup.WithContext(context.Background())

	grpcServ := grpc.NewGrpcServer(ctx, grpcPort, a)

	httpServer := httpgin.NewHTTPServer(ctx, httpPort, a)

	sigQuit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			log.Printf("captured signal: %v\n", s)
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	// run grpc server
	eg.Go(grpcServ.Run)

	eg.Go(httpServer.Run)

	if err := eg.Wait(); err != nil {
		log.Printf("gracefully shutting down the servers: %s\n", err.Error())
	}

	log.Println("servers were successfully shutdown")
}
