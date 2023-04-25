package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/tinrab/retry"
	"log"
	"net/http"
	"os"
	"os/signal"
	"test0/internal"
	"test0/internal/db"
	"time"
)

func main() {
	App := internal.NewApplication()
	cfg := LoadEnvVariables()
	srv := http.Server{
		Addr:     cfg.ServerAddress,
		ErrorLog: App.ErrorLog,
		Handler:  App.Routes(),
	}

	App.InfoLog.Printf("Launch server on %s", srv.Addr)

	//connect to database

	retry.ForeverSleep(2*time.Second, func(attempt int) error {
		addr := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
		repo, err := db.NewPostgres(addr)
		if err != nil {
			log.Println(err)
			return err
		}
		db.SetRepository(repo)
		return nil
	})
	defer db.Close()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			App.ErrorLog.Fatalln(err)
		}
	}()

	//graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		App.ErrorLog.Fatalln(err)
	}
	App.InfoLog.Printf("Server closing...")
}
