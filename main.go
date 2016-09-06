package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

const (
	CONFIG_FILE = "config.json"
)

func dumpRequest(r *http.Request) {
	bytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(bytes))
}

func main() {
	cfg, err := loadConfig(CONFIG_FILE)
	if err !=nil {
		log.Fatal(err)
	}

	log.Println(cfg)

	r := mux.NewRouter()
	posts := r.Methods("POST").Subrouter()

	for _, bc := range cfg.Builds {
		posts.HandleFunc("/" + bc.Name, func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			dumpRequest(r)
			job := NewJob(now, &bc)
			log.Println("Received request to " + r.URL.Path + " for job " + job.Id)
			w.WriteHeader(http.StatusOK)
			go job.Run()
		})
	}
	log.Fatal(http.ListenAndServe(":80", r))
}
