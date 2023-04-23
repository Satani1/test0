package models

type Order struct {
	orderUID    string `json:"order_uid"`
	trackNumber string `json:"track_number"`
	entry       string `json:"entry "`
	delivery    struct {
		name    string `json:"name"`
		phone   string `json:"phone"`
		zip     string `json:"zip"`
		city    string `json:"city"`
		address string `json:"address"`
		region  string `json:"region"`
		email   string `json:"email"`
	} `json:"delivery"`
	payment struct {
		transaction  string `json:"transaction"`
		requestID    string `json:"request_id"`
		currency     string `json:"currency"`
		provider     string `json:"provider"`
		amount       int    `json:"amount"`
		paymentDT    int    `json:"payment_dt"`
		bank         string `json:"bank"`
		deliveryCost int    `json:"delivery_cost"`
		goodsTotal   int    `json:"goods_total"`
		customFee    int    `json:"custom_fee"`
	} `json:"payment"`
	items []struct {
		chrtID      int    `json:"chrt_id"`
		trackNumber string `json:"track_number"`
		price       int    `json:"price"`
		rid         string `json:"rid"`
		name        string `json:"name"`
		sale        int    `json:"sale"`
		size        string `json:"size"`
		totalPrice  int    `json:"total_price"`
		nmID        int    `json:"nm_id"`
		brand       string `json:"brand"`
		status      int    `json:"status"`
	} `json:"items"`
	locale            string `json:"locale"`
	internalSignature string `json:"internal_signature"`
	customerID        string `json:"customer_id"`
	deliveryService   string `json:"delivery_service"`
	shardkey          string `json:"shardkey"`
	smID              int    `json:"sm_id"`
	dateCreated       string `json:"date_created"`
	oofShard          string `json:"oof_shard"`
}
