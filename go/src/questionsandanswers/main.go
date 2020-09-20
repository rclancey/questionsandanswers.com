package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var db *DB
var accessLogger *AccessLogger

var port int
var logdir string
var docroot string
var proxy string
var dbfn string

func init() {
	flag.IntVar(&port, "port", 8080, "HTTP Port")
	flag.StringVar(&logdir, "log-dir", "./var/log", "Log directory")
	flag.StringVar(&docroot, "document-root", ".", "Static file document root")
	flag.StringVar(&proxy, "default-proxy", "", "Default proxy URL for unmatched paths")
	flag.StringVar(&dbfn, "database", "./questions.db", "Questions database SQLite3 filename")
}

func main() {
	flag.Parse()
	var err error
	accessLogger, err = NewAccessLogger(filepath.Join(logdir, "access.log"))
	if err != nil {
		log.Fatal(err)
		return
	}
	errLogFn := filepath.Join(logdir, "error.log")
	errLogF, err := os.OpenFile(errLogFn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	db, err = NewDB(dbfn)
	if err != nil {
		log.Fatal(err)
		return
	}
	router, err := NewRouter(docroot, proxy)
	if err != nil {
		log.Fatal(err)
		return
	}
	router.AddRoute(http.MethodGet, "/api/questions", ListQuestions)
	router.AddRoute(http.MethodGet, "/api/question", GetQuestion)
	router.AddRoute(http.MethodPost, "/api/question", CreateQuestion)
	router.AddRoute(http.MethodPut, "/api/question", AnswerQuestion)
	log.SetOutput(errLogF)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	log.Fatal(err)
}
