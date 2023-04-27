package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
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
	addr := fmt.Sprintf("postgres://postgres:12345678@localhost:5432/postgres?sslmode=disable")
	repo, err := db.NewPostgres(addr)
	if err != nil {
		App.ErrorLog.Println(err)
		return
	}
	db.SetRepository(repo)
	defer db.Close()

	//connect to stan
	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID, stan.NatsURL(cfg.NatsURL), stan.MaxPubAcksInflight(1000))
	if err != nil {
		App.ErrorLog.Println("cant connect to stan", err)

	}

	_, err = sc.Subscribe("orders", func(msg *stan.Msg) {
		//insert data from message into cash and postgres
		err := InsertDataFromMessage(msg.Data)
		if err != nil {
			App.ErrorLog.Println("cant insert message in repo", err)
		}
		fmt.Printf("Received a message: %s\n", string(msg.Data))
	}, stan.DeliverAllAvailable(), stan.DurableName("my-durable"))
	if err != nil {
		App.ErrorLog.Println("Cant subscribe to channel", err)

	}

	//running http server
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			App.ErrorLog.Fatalln(err)
		}
	}()

	//graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	//close nats-stream connection and stop the server
	if err := sc.Close(); err != nil {
		App.ErrorLog.Println(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		App.ErrorLog.Fatalln("Cannot shutdown the server", err)
	}

	App.InfoLog.Printf("Server closing...")
}

func InsertDataFromMessage(data []byte) error {
	log.Println("insert data in db")
	var or db.CreateOrder
	err := json.Unmarshal(data, &or)
	if err != nil {
		return err
	}

	co := db.CreateOrder{
		OrderUID: uuid.New().String(),
		Data:     data,
	}
	fmt.Println(co.OrderUID)
	if err = db.InsertRow(context.Background(), co); err != nil {
		return err
	}

	return nil
}
