package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fileName := r.RequestURI[1:]

		if r.RequestURI == "/" {
			fileName = "index.html"
		}

		log.Printf("opening %s", fileName)

		file, err := ioutil.ReadFile(fileName)

		if err != nil {
			log.Print(err)
			http.Error(w, "Not found", 404)
			return
		}

		w.Write(file)
	})

	port := 3010
	log.Printf("Listening on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		log.Fatal(err)
	}
}
