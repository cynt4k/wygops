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
	"github.com/cynt4k/wygops/cmd/config"
	"github.com/spf13/viper"
)

var (
	//Version : Populated during build time
	Version string = "Unknown"
	//Build : Populated during build time
	Build string = "Unknown"
)

func main() {

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	log.Printf("starting wygops version: %s - build: %s", Version, Build)
	shutdownCtx, shutdown := context.WithCancel(context.Background())
	defer shutdown()

	config := config.Config{}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.ServeServer()

	if err != nil {
		log.Fatalf("error while serving application %s", err)
	}

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
