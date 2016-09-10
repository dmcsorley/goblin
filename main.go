package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

const (
	CONFIG_FILE = "config.json"
	LISTEN_ADDR = ":80"
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
		log.Fatal("Error loading server config: " + err.Error())
	}

	log.Println(cfg)

	hostname, _ := os.Hostname()
	log.Println(fmt.Sprintf("Listening on %s%s", hostname, LISTEN_ADDR))

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
	log.Fatal("Error starting http server: " + http.ListenAndServe(LISTEN_ADDR, r).Error())
}
