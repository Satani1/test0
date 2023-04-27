package internal

import (
	"log"
	"os"
)

type Application struct {
	Addr     string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

func NewApplication() *Application {
	//logs
	InfoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	ErrorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	App := &Application{
		ErrorLog: ErrorLog,
		InfoLog:  InfoLog,
	}
	return App
}
