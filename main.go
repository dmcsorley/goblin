package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

const (
	CONFIG_FILE = "config.json"
)

func isValidRequest(r *http.Request, cfg *ServerConfig) *BuildConfig {
	if strings.ToUpper(r.Method) != "POST" {
		return nil
	}

	return cfg.BuildConfigForPath(r.URL.Path)
}

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		dumpRequest(r)
		bc := isValidRequest(r, cfg)
		if bc != nil {
			job := NewJob(now, bc)
			log.Println("Received request to " + r.URL.Path + " for job " + job.Id)
			w.WriteHeader(http.StatusOK)
			go job.Run()
		} else {
			log.Println("Bad request for " + r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
		}

	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
