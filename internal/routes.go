package internal

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (app *Application) Routes() *mux.Router {
	rMux := mux.NewRouter()

	rMux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("hello users!")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	rMux.HandleFunc("/", app.RenderIndex)
	rMux.HandleFunc("/order/{order_uid}", app.RenderOrder)
	rMux.HandleFunc("/test", app.GetOrder)

	fileServer := http.FileServer(http.Dir("./web"))

	rMux.PathPrefix("/web/").Handler(http.StripPrefix("/web", fileServer))

	return rMux
}
