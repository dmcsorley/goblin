package main

import (
	"flag"
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
	ENV_IMAGE = "IMAGE"
	RUN_FLAG = "run"
	TIME_FLAG = "time"
)

func dumpRequest(r *http.Request) string {
	bytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(bytes)
}

var runFlag string
var timeFlag string

func init() {
	flag.StringVar(&runFlag, RUN_FLAG, "", "build to run from config file")
	flag.StringVar(&timeFlag, TIME_FLAG, "", "timestamp of job to run")
}

func main() {
	flag.Parse()

	cfg, err := loadConfig(CONFIG_FILE)
	if err !=nil {
		log.Fatal("Error loading server config: " + err.Error())
	}

	if runFlag != "" && timeFlag != "" {
		runJob(cfg, runFlag, timeFlag)
	} else {
		serve(cfg)
	}
}

func serve(cfg *ServerConfig) {
	image := os.Getenv(ENV_IMAGE)
	if image == "" {
		log.Fatal(ENV_IMAGE + " environment variable is required")
	}

	hostname, _ := os.Hostname()
	log.Println(fmt.Sprintf("Listening on %s%s", hostname, LISTEN_ADDR))

	r := mux.NewRouter()
	posts := r.Methods("POST").Subrouter()

	for _, bc := range cfg.Builds {
		localConfig := bc
		log.Println("Build configured on /" + localConfig.Name)
		posts.HandleFunc("/" + localConfig.Name, func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			joblog("DEBUG", dumpRequest(r), os.Stdout)
			job := NewJob(now, &localConfig)
			joblog(job.Id, "Received build for " + r.URL.Path, os.Stdout)
			w.WriteHeader(http.StatusOK)
			go job.DockerRun()
		})
	}
	log.Fatal("Error starting http server: " + http.ListenAndServe(LISTEN_ADDR, r).Error())
}

func runJob(cfg *ServerConfig, buildName string, timeStamp string) {
	bc := cfg.FindBuildByName(buildName)
	if bc == nil {
		log.Fatal("No build found with name " + buildName)
	}

	t, err := time.Parse(JOB_TIME_FORMAT, timeStamp)
	if err != nil {
		log.Fatal(err)
	}

	job := NewJob(t, bc)
	job.Run()
}
