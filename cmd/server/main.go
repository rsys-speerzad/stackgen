package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rsys-speerzad/stackgen/pkg/configs"
	"github.com/rsys-speerzad/stackgen/pkg/router"
	"github.com/rsys-speerzad/stackgen/pkg/store"
	"github.com/rsys-speerzad/stackgen/pkg/testing"
)

func main() {
	Run()
}

// Run starts the HTTP server
func Run() {
	// parse the environment variables
	configPath := flag.String("config", "config.json", "path to config file")
	createTestData := flag.Bool("testdata", false, "create test data")
	flag.Parse()
	configs.ParseEnv(*configPath)
	// initialize and auto migrate the db schemas
	store.AutoMigrate()
	// create test data if needed
	if *createTestData {
		if err := testing.CreateTestData(store.GetDB()); err != nil {
			log.Fatalf("Failed to create test data: %v", err)
		}
	}
	// start API server
	server := router.NewServer()
	// gracefully close the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		println()
		log.Println("Closing db connection...")
		store.CloseDB()
		log.Println("Shutting down server...")
		if err := gracefulShutdown(server, 25*time.Second); err != nil {
			log.Printf("Server stopped: %s", err.Error())
		}
		os.Exit(0)
	}()
	// start the server
	log.Printf("Listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func gracefulShutdown(server *http.Server, maximumTime time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), maximumTime)
	defer cancel()
	return server.Shutdown(ctx)
}
