package internal

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"test0/internal/db"
	"text/template"
)

func (app *Application) RenderIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		orderUID := r.FormValue("order_uid")
		http.Redirect(w, r, "/order/"+orderUID, http.StatusSeeOther)
	} else {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ts, err := template.ParseFiles("./web/html/index.html")
		if err != nil {
			app.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = ts.Execute(w, nil)
		if err != nil {
			app.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (app *Application) RenderOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		params := mux.Vars(r)
		order_uid := params["order_uid"]

		orderData, err := db.GetOrder(order_uid)
		if err != nil {
			app.ErrorLog.Println(err)
			http.Error(w, errors.New("not found any orders with that ID").Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		ts, err := template.ParseFiles("./web/html/order.html")
		if err != nil {
			app.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
	orderData, err := db.GetOrder("b563feb7b2b84b6test")
	if err != nil {
		app.ErrorLog.Fatalln(err)
	}

	if err := json.NewEncoder(w).Encode(orderData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)

}

func (app *Application) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var params db.CreateOrder

	if err := json.NewDecoder(r.Body).Decode(&params.Data); err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err := db.InsertRow(context.Background(), params)
	if err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
