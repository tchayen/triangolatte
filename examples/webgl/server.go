package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// apiHandler tries to respond with matching API route.
func apiHandler(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.RequestURI, "/")

	// Any path shorter than [<empty>, "api", "data"] must be incorrect.
	if len(pathSegments) < 3 {
		http.Error(w, "Wrong path", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	switch pathSegments[2] {
	case "data":
		data, err := ioutil.ReadFile("../../assets/json_tmp")

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Write(data)

	case "report":
		type Body struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		}
		data, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var body Body
		err = json.Unmarshal(data, &body)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		log.Printf("id: %s, status: %s", body.ID, body.Status)

	default:
		w.Write([]byte("Error: wrong API path"))
	}
}

// Directly open file based on URL.
func serveFile(fileName string, w http.ResponseWriter) {
	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Printf("%s", err)
		http.Error(w, "Not found", 404)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

// staticHandler aims to return static file based on request URL.
func staticHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" {
		index, err := ioutil.ReadFile("./index.html")

		if err != nil {
			log.Fatal(err)
		}

		w.Write(index)
		return
	}

	serveFile(r.RequestURI[1:], w)
}

func main() {
	http.HandleFunc("/", staticHandler)
	http.HandleFunc("/api/", apiHandler)

	port := 3000
	log.Printf("Listening on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		log.Fatal(err)
	}
}
