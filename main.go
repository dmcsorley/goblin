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

func dumpRequest(r *http.Request) string {
	bytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(bytes)
}

func main() {
	cfg, err := loadConfig(CONFIG_FILE)
	if err !=nil {
		log.Fatal("Error loading server config: " + err.Error())
	}

	hostname, _ := os.Hostname()
	log.Println(fmt.Sprintf("Listening on %s%s", hostname, LISTEN_ADDR))

	r := mux.NewRouter()
	posts := r.Methods("POST").Subrouter()

	for _, bc := range cfg.Builds {
		log.Println("Build configured on /" + bc.Name)
		posts.HandleFunc("/" + bc.Name, func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			joblog("DEBUG", dumpRequest(r), os.Stdout)
			job := NewJob(now, &bc)
			joblog(job.Id, "Received build for " + r.URL.Path, os.Stdout)
			w.WriteHeader(http.StatusOK)
			go job.Run()
		})
	}
	log.Fatal("Error starting http server: " + http.ListenAndServe(LISTEN_ADDR, r).Error())
}
