package internal

import (
	"net/http"
	"test0/internal/models"
	"text/template"
)

func (app *Application) RenderIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		orderUID := r.FormValue("order_uid")
		http.Redirect(w, r, "/order?id="+orderUID, http.StatusSeeOther)
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
		order_uid := r.URL.Query().Get("id")
		app.InfoLog.Printf("%s\n", order_uid)

		ts, err := template.ParseFiles("./web/html/order.html")
		if err != nil {
			app.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		o := models.Order{
			OrderUID: order_uid,
		}
		err = ts.Execute(w, o)
		if err != nil {
			app.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
