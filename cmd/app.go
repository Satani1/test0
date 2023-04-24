package main

import (
	"context"
	"errors"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Application struct {
	Addr     string
	errorLog *log.Logger
	infoLog  *log.Logger
}

func NewApplication() *Application {
	//addr config from terminal
	addr := flag.String("addr", "localhost:8080", "Server Address")
	flag.Parse()
	//logs
	InfoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	ErrorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	App := &Application{
		errorLog: ErrorLog,
		infoLog:  InfoLog,
		Addr:     *addr,
	}
	return App
}

func (app *Application) Routes() *mux.Router {
	rMux := mux.NewRouter()

	rMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("hello users!")); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return rMux
}

func main() {
	App := NewApplication()

	srv := http.Server{
		Addr:     App.Addr,
		ErrorLog: App.errorLog,
		Handler:  App.Routes(),
	}

	App.infoLog.Printf("Launch server on %s", App.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			App.errorLog.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		App.errorLog.Fatalln(err)
	}
	App.infoLog.Printf("Server closing...")
}
