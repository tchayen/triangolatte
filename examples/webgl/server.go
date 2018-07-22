package main

import (
    "net/http"
    "log"
    . "triangolatte"
    "io/ioutil"
    "encoding/json"
)

func main() {
    // Triangulate data
    v := [] Point{{50, 110}, {150, 30}, {240, 115}, {320, 65}, {395, 170}, {305, 160}, {265, 240}, {190, 100}, {95, 125}, {100, 215}}
    triangulated, err := EarCut(v, [][]Point{})
    if err != nil {
        log.Fatal("Failed to triangulate polygon")
    }

    // Marshal to JSON
    data, err := json.Marshal(triangulated)
    if err != nil {
        log.Fatal("Failed to marshal JSON")
    }

    // Save in a file
    err = ioutil.WriteFile("polygon_tmp", data, 0644)
    if err != nil {
        log.Fatal("Failed to save file")
    }

    fs := http.FileServer(http.Dir(""))
    http.Handle("/", fs)

    log.Println("Listening on :3000...")
    http.ListenAndServe(":3000", nil)
}
