package main

import (
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/nats-io/stan.go"
	"log"
	"test0/internal/models"
)

func main() {
	sc, err := stan.Connect("test-cluster", "client-125", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Println(err)
	}
	fmt.Println("send incorrect data")
	//send incorrect data
	for i := 0; i < 5; i++ {
		if err := sc.Publish("orders", []byte("{\"sale\": 30, \"size\": \"0\",  \"price\": 453, \"status\": 202, \"chrt_id\": 9934930, \"total_price\": 317, \"track_number\": \"WBILMTESTTRACK\"}], \"sm_id\": 99, \"locale\": \"en\", \"payment\": {\"bank\": \"alpha\", \"amount\": 1817, \"currency\": \"RUB\", \"provider\": \"wbpay\", \"custom_fee\": 0, \"payment_dt\": 1637907727, \"request_id\": \"\", \"goods_total\": 317, \"transaction\": \"teasrasd\", \"delivery_cost\": 1500}, \"delivery\": {\"zip\": \"2639809\", \"city\": \"Kiryat Mozkin\", \"name\": \"Test Testov\", \"email\": \"test@gmail.com\", \"phone\": \"+9720000000\", \"region\": \"Kraiot\", \"address\": \"Ploshad Mira 15\"}, \"shardkey\": \"9\", \"oof_shard\": \"1\", \"order_uid\": \"test\", \"customer_id\": \"test\", \"date_created\": \"2021-11-26T06:22:19Z\", \"track_number\": \"WBILMTESTTRACK\", \"delivery_service\": \"meest\", \"internal_signature\": \"\"}")); err != nil {
			log.Println(err)
		}
	}
	fmt.Println("another incorrect data")
	//another incorrect data
	for i := 0; i < 5; i++ {
		if err := sc.Publish("orders", []byte("some fake data :3")); err != nil {
		}
	}

	//random data

	fmt.Println("send random data")
	var testOrder models.Order
	err = gofakeit.Struct(&testOrder)
	if err != nil {
		log.Println(err)
	}
	nOrder, err := json.Marshal(testOrder)
	if err != nil {
		log.Println(err)
	}
	if err := sc.Publish("orders", []byte(nOrder)); err != nil {
		log.Println(err)
	}
	fmt.Println("close")
}
