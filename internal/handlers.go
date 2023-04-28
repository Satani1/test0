package internal

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nats-io/stan.go"
	"net/http"
	"test0/internal/db"
	"text/template"
)

func (app *Application) RenderIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//get order id from html form
		orderUID := r.FormValue("order_uid")

		//redirect to page with render of order information
		url := "/order/" + orderUID
		http.Redirect(w, r, url, http.StatusSeeOther)
	} else {
		//find html template
		ts, err := template.ParseFiles("./web/html/index.html")
		if err != nil {
			app.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		//execute html page
		err = ts.Execute(w, nil)
		if err != nil {
			app.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (app *Application) RenderOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//get order uid from url
		params := mux.Vars(r)
		order_uid := params["order_uid"]

		//get order info from db
		//orderData, err := db.GetOrder(order_uid)
		//if err != nil {
		//	app.ErrorLog.Println(err)
		//	http.Error(w, errors.New("not found any orders with that ID").Error(), http.StatusBadRequest)
		//	return
		//}

		//get order info from memory cache
		orderData, err := app.AppCache.GetOrder(order_uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		//find html template
		w.Header().Set("Content-Type", "text/html")
		ts, err := template.ParseFiles("./web/html/order.html")
		if err != nil {
			app.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//execute html template with order data
		err = ts.Execute(w, orderData)
		if err != nil {
			app.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (app *Application) GetOrder(w http.ResponseWriter, r *http.Request) {
	//orderData, err := db.GetOrder("b563feb7b2b84b6test")
	//if err != nil {
	//	app.ErrorLog.Fatalln(err)
	//}

	//get order info from memory cache
	orderData, err := app.AppCache.GetOrder("b563feb7b2b84b6test")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//encoding that data in json
	if err := json.NewEncoder(w).Encode(orderData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)

}

func (app *Application) CreateOrder(w http.ResponseWriter, r *http.Request) {
	//read json data from request and place it to an order struct
	var params db.CreateOrder
	if err := json.NewDecoder(r.Body).Decode(&params.Data); err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//err := db.InsertRow(context.Background(), params)
	//if err != nil {
	//	app.ErrorLog.Println(err)
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//

	//connect to stan
	sc, err := stan.Connect("test-cluster", "client-124", stan.NatsURL("nats://localhost:6060"))
	if err != nil {
		app.ErrorLog.Println("cant connect stan", err)
	}
	//publish message with order data
	if err := sc.Publish("orders", params.Data); err != nil {
		app.ErrorLog.Println("Cant publish msg", err)
	}
	defer sc.Close()

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte("Order was created with ID: " + params.OrderUID)); err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
