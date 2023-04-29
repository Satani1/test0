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
	mCache "test0/internal/cache"
	"test0/internal/db"
	"test0/internal/models"
	"time"
)

func main() {
	//read environmental variables
	cfg := LoadEnvVariables()

	//connect to database
	addr := fmt.Sprintf("postgres://postgres:12345678@localhost:5432/postgres?sslmode=disable")
	repo, err := db.NewPostgres(addr)
	if err != nil {
		log.Println(err)
		return
	}
	db.SetRepository(repo)
	defer db.Close()

	//cache memory create
	CacheMemoryApp := mCache.NewCache(0*time.Minute, 0*time.Minute)
	//restore all orders data in cache memory
	if err := CacheMemoryApp.Restore(); err != nil {
		log.Println("Cat restore memory in cache", err)
	}
	log.Println(CacheMemoryApp.Cache.ItemCount())
	//connect to stan
	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID, stan.NatsURL(cfg.NatsURL), stan.MaxPubAcksInflight(1000))
	if err != nil {
		log.Println("cant connect to stan", err)

	}
	//subscribe to channel "orders"
	_, err = sc.Subscribe("orders", func(msg *stan.Msg) {
		//insert data from message into cash and postgres
		err := InsertDataFromMessage(msg.Data, CacheMemoryApp)
		if err != nil {
			log.Println("cant insert message in repo", err)
		}
		fmt.Printf("Received a message: %s\n", string(msg.Data))
	}, stan.DeliverAllAvailable(), stan.DurableName("my-durable"))
	if err != nil {
		log.Println("Cant subscribe to channel", err)

	}
	//application and server setup
	App := internal.NewApplication(*CacheMemoryApp)
	srv := http.Server{
		Addr:     cfg.ServerAddress,
		ErrorLog: App.ErrorLog,
		Handler:  App.Routes(),
	}
	log.Printf("Launch server on %s", srv.Addr)
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

func InsertDataFromMessage(data []byte, c *mCache.CacheMemory) error {
	log.Println("insert data in db")
	var or models.CreateOrder
	err := json.Unmarshal(data, &or)
	if err != nil {
		return err
	}
	co := models.CreateOrder{
		OrderUID: uuid.New().String(),
		Data:     data,
	}
	//save to cache
	var order models.Order
	if err := json.Unmarshal(co.Data, &order); err != nil {
		return err
	}
	err = c.Put(co.OrderUID, order)
	if err != nil {
		return err
	}
	//save to db
	if err = db.InsertRow(context.Background(), co); err != nil {
		return err
	}

	return nil
}
