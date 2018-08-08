package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"triangolatte"
)

type Building struct {
	Name string
	Points [][]triangolatte.Point
}

func parseData(m map[string]interface{}) (buildings []Building) {
	// This part is really ugly, but gets the job done with converting
	// unstructured JSON to GO.
	buildings = make([]Building, 0)
	for i, f := range m["features"].([]interface{}) {
		// Extract 'feature'.
		feature := f.(map[string]interface{})

		// Initialize new building.
		b := Building{Name: feature["properties"].(map[string]interface{})["name"].(string)}
		buildings = append(buildings, b)

		// Extract 'geometry'.
		geometry := feature["geometry"].(map[string]interface{})

		// Pay price for strict typing with no algebraic data types, i.e. switch
		// handle different geometry types that might happen.
		switch geometry["type"] {
		case "Polygon":
			for j, polygon := range geometry["coordinates"].([]interface{}) {
				// Initialize points array in the building.
				buildings[i].Points = append(buildings[i].Points, []triangolatte.Point{})

				for _, p := range polygon.([]interface{}) {
					// Cast from interface{} to []interface{}.
					pointArray := p.([]interface{})

					point := triangolatte.Point{
						X: pointArray[0].(float64),
						Y: pointArray[1].(float64),
					}

					// Convert coordinates.
					pointInMeters := triangolatte.DegreesToMeters(point)
					buildings[i].Points[j] = append(buildings[i].Points[j], pointInMeters)
				}
			}
		case "LineString":
		case "Point":
		}
	}
	return
}

func calculateStats(buildings []Building) (totalSuccesses, totalErrors int) {
	for _, b := range buildings {
		if len(b.Points) == 0 {
			continue
		}

		error := false
		cleaned, err := triangolatte.EliminateHoles(b.Points)

		if err != nil {
			fmt.Printf("Error in hole removal: %s\n", err)
			error = true
		}

		triangles, err := triangolatte.EarCut(cleaned)

		if err != nil {
			fmt.Printf("Error in triangulation: %s\n", err)
			error = true
		}

		var holes [][]triangolatte.Point
		if len(b.Points) > 1 {
			holes = b.Points[1:]
		} else {
			holes = [][]triangolatte.Point{}
		}
		_, _, deviation := triangolatte.Deviation(b.Points[0], holes, triangles)

		if deviation > 1e-6 {
			fmt.Printf("Error detected in deviation\n")
			error = true
		}

		if error {
			totalErrors += 1
		} else {
			totalSuccesses += 1
			fmt.Printf("%#v\n", triangles)
		}
	}
	return
}

func main() {
	// Load data
	data, err := ioutil.ReadFile("assets/cracow_tmp")

	if err != nil {
		log.Fatal("Could not read file")
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	buildings := parseData(m)
	totalSuccesses, totalErrors := calculateStats(buildings)
	fmt.Printf("success: %d failure: %d", totalSuccesses, totalErrors)
}
