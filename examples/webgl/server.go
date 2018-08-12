package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func dataHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("assets/json_tmp")

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func main() {
	http.HandleFunc("/data", dataHandler)

	log.Println("Listening on :3000...")
	http.ListenAndServe(":3000", nil)
}
