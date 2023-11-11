package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/env"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func main() {
	if err := realMain(); err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	env, err := env.GetEnv()
	if err != nil {
		return err
	}

	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	todoDB, err := db.NewDB(env.DBPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	mux := router.NewRouter(env, todoDB)
	srv := &http.Server{
		Addr:    env.Port,
		Handler: mux,
	}

	return startServer(srv)
}

func startServer(srv *http.Server) error {
	// set up signal handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	// start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// wait for signal
	<-ctx.Done()

	// shutdown the server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown failed: %v", err)
		return err
	}

	log.Println("Server gracefully shutdown")
	return nil
}
