package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// TODO: serve data API endpoint at /data:
// 1. load file with JSON
// 2. encode the JSON
// 3. serve the file

func dataHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("assets/json_tmp")

	var f [][][]float64
	json.Unmarshal(data, &f)

	if err != nil {
		log.Fatal(err)
	}
	// ^ doesn't really work

	json.NewEncoder(w).Encode(data)
}

func main() {
	// // Triangulate data
	// v := []Point{{50, 110}, {150, 30}, {240, 115}, {320, 65}, {395, 170}, {305, 160}, {265, 240}, {190, 100}, {95, 125}, {100, 215}}
	// triangulated, err := EarCut(v)
	// if err != nil {
	// 	log.Fatal("Failed to triangulate polygon")
	// }
	//
	// // Marshal to JSON
	// data, err := json.Marshal(triangulated)
	// if err != nil {
	// 	log.Fatal("Failed to marshal JSON")
	// }
	//
	// // Save in a file
	// err = ioutil.WriteFile("polygon_tmp", data, 0644)
	//
	// if err != nil {
	// 	log.Fatal("Failed to save file")
	// }
	//
	// fs := http.FileServer(http.Dir(""))

	http.HandleFunc("/data", dataHandler)

	log.Println("Listening on :3000...")
	http.ListenAndServe(":3000", nil)
}
