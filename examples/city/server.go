package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"triangolatte"
)

type Building struct {
	Properties map[string]string
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
		b := Building{Properties: map[string]string{}}

		// Rewrite properties.
		for k, v := range feature["properties"].(map[string]interface{}) {
			switch value := v.(type) {
			case string:
				b.Properties[k] = value
			}
		}

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

func normalizeCoordinates(buildings []Building) {
	// TODO: normalize coordinates of the buildings in range [0.0, 1.0].
}

func triangulate(buildings []Building) (triangles [][]float64, totalSuccesses, totalErrors int) {
	triangles = make([][]float64, len(buildings))

	for i, b := range buildings {
		if len(b.Points) == 0 {
			continue
		}

		error := false
		cleaned, err := triangolatte.EliminateHoles(b.Points)

		if err != nil {
			error = true
		}

		t, err := triangolatte.EarCut(cleaned)

		if err != nil {
			error = true
		}

		var h [][]triangolatte.Point
		if len(b.Points) > 1 {
			h = b.Points[1:]
		} else {
			h = [][]triangolatte.Point{}
		}
		_, _, deviation := triangolatte.Deviation(b.Points[0], h, t)

		triangles[i] = t
		// 1e-6 was chosen arbitrarily as a frontier between low and high error
		// rate.
		if deviation > 1e-6 {
			error = true
		}

		if error {
			totalErrors += 1
		} else {
			totalSuccesses += 1
		}
	}
	return
}

func main() {
	// Load data from file.
	data, err := ioutil.ReadFile("assets/cracow_tmp")

	if err != nil {
		log.Fatal("Could not read file")
	}

	// Parse JSON.
	var m map[string]interface{}
	json.Unmarshal(data, &m)

	// Translate data to a more handy format.
	buildings := parseData(m)

	// Normalize coordinates.
	normalizeCoordinates(buildings)

	// Check out what went right and what wrong.
	_, successes, errors := triangulate(buildings)

	// Brag about success (or admit to poor performance, who knows...).
	fmt.Printf("success: %d failure: %d", successes, errors)
}
