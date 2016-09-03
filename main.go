package main

import (
        "log"
        "net/http"
        "net/http/httputil"
)

func main() {
        http.HandleFunc("/",  func(w http.ResponseWriter, r *http.Request) {
		bytes, err := httputil.DumpRequest(r, true)

		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(bytes))
                w.WriteHeader(http.StatusOK)
        })
        log.Fatal(http.ListenAndServe(":80", nil))
}
