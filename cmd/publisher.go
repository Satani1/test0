package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
	"log"
	"test0/internal/models"
	"time"
)

type PublisherApp struct {
	StanConnection stan.Conn
	ClientID       string
	ClusterID      string
	NatsURL        string
	MessageDelay   time.Duration
}

func (pub *PublisherApp) PublishMessage() {
	o := new(models.Order)
	o.OrderUID = uuid.New().String()

	bytesOrder, err := json.Marshal(o)
	if err != nil {
		log.Println(err)
	}

	err = pub.StanConnection.Publish("orders", bytesOrder)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Order with ID %s is sent", o.OrderUID)
}
