package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"triangolatte"
)

func main() {
	type Point triangolatte.Point

	// Load data
	data, err := ioutil.ReadFile("assets/cracow_tmp")

	if err != nil {
		log.Fatal("Could not read file")
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	buildings := make([][][]Point, 0)
	for i, f := range m["features"].([]interface{}) {
		buildings = append(buildings, [][]Point{})
		geometry := f.(map[string]interface{})["geometry"].(map[string]interface{})
		switch geometry["type"] {
		case "Polygon":
			for j, polygon := range geometry["coordinates"].([]interface{}) {
				buildings[i] = append(buildings[i], []Point{})
				for _, p := range polygon.([]interface{}) {
					buildings[i][j]
					point := p.([]float64)
					Point{point[0], point[1]}
				}
				//append(buildings, )
			}
		case "LineString":
		case "Point":
		}
	}

	// fmt.Println(m)

	// // Define types
	// type Geometry struct {
	// 	Type        string        `json:"type"`
	// 	Coordinates [][][]float64 `json:"coordinates"`
	// }

	// type Feature struct {
	// 	Geometry Geometry `json:"geometry"`
	// }

	// type FeatureCollection struct {
	// 	Features []Feature `json:"features"`
	// }

	// type Point triangolatte.Point

	// // Parse
	// var parsed FeatureCollection
	// json.Unmarshal(data, &parsed)

	// // Translate to our data format
	// var points [][][]Point

	// points = make([][][]Point, len(parsed.Features))

	// // Buildings
	// for i := 0; i < len(parsed.Features); i++ {
	// 	coords := parsed.Features[i].Geometry.Coordinates
	// 	points[i] = make([][]Point, len(coords))

	// 	// Features (shape and holes)
	// 	for j := 0; j < len(coords); j++ {
	// 		points[i][j] = make([]Point, len(coords[j]))
	// 		for k := 0; k < len(coords[j]); k++ {
	// 			points[i][j][k] = Point{coords[j][k][0], coords[j][k][1]}
	// 		}
	// 	}
	// }
	//
	// fmt.Println(points)
}
