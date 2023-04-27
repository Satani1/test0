package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/tinrab/retry"
	"log"
	"net/http"
	"os"
	"os/signal"
	"test0/internal"
	"test0/internal/db"
	"test0/internal/models"
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
		addr := fmt.Sprintf("postgres://postgres:12345678@localhost:5432/postgres?sslmode=disable")
		repo, err := db.NewPostgres(addr)
		if err != nil {
			log.Println(err)
			return err
		}
		db.SetRepository(repo)
		return nil
	})
	defer db.Close()
	//connect to stan
	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID, stan.NatsURL(cfg.NatsURL), stan.MaxPubAcksInflight(1000))
	if err != nil {
		App.ErrorLog.Println("cant connect to stan")

	}

	sub, err := sc.Subscribe("orders", func(msg *stan.Msg) {
		//insert data from message into cash and postgres
		if err := InsertDataFromMessage(msg.Data); err != nil {
			App.ErrorLog.Println(err)
		}
		fmt.Printf("Received a message: %s\n", string(msg.Data))
	}, stan.DeliverAllAvailable(), stan.DurableName("my-durable"))
	if err != nil {
		App.ErrorLog.Println(err)

	}
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
	//stan off
	sub.Unsubscribe()
	sc.Close()
	App.InfoLog.Printf("Server closing...")
}

func InsertDataFromMessage(data []byte) error {
	or := new(models.Order)
	err := json.Unmarshal(data, or)
	if err != nil {
		return err
	}

	co := db.CreateOrder{
		OrderUID: or.OrderUID,
		Data:     data,
	}

	if err = db.InsertRow(context.Background(), co); err != nil {
		return err
	}

	return nil
}
