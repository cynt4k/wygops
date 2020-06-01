package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cynt4k/wygops/cmd"
)

var (
	//Version : Populated during build time
	Version string
	//Build : Populated during build time
	Build string
)

func main() {

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	log.Printf("starting wygops version: %s - build: %s", Version, Build)
	shutdownCtx, shutdown := context.WithCancel(context.Background())
	defer shutdown()

	err = waitForInterrupt(shutdownCtx)
	log.Printf("shuting down: %s", err)
}

func waitForInterrupt(ctx context.Context) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-c:
		return fmt.Errorf("received signal %s", sig)
	case <-ctx.Done():
		return errors.New("canceled")
	}
}
