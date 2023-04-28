package internal

import (
	"log"
	"os"
	mCache "test0/internal/cache"
)

type Application struct {
	Addr     string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	AppCache mCache.CacheMemory
}

func NewApplication(c mCache.CacheMemory) *Application {
	//logs
	InfoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	ErrorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	App := &Application{
		ErrorLog: ErrorLog,
		InfoLog:  InfoLog,
		AppCache: c,
	}
	return App
}
