package main

import (
	"fmt"
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

func serveFile(fileName string, w http.ResponseWriter) {
	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Printf("%s", err)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("File not found"))
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	index, err := ioutil.ReadFile("examples/webgl/index.html")

	if err != nil {
		log.Fatal(err)
	}

	if r.RequestURI == "/" {
		w.Write(index)
		return
	}

	serveFile(fmt.Sprintf("./examples/webgl/%s", r.RequestURI[1:]), w)
}

func main() {
	http.HandleFunc("/", staticHandler)
	http.HandleFunc("/data", dataHandler)

	port := 3000
	log.Printf("Listening on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		log.Fatal(err)
	}
}
