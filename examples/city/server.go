package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"triangolatte"
)

func main() {
	// Load data
	data, err := ioutil.ReadFile("assets/cracow_tmp")

	if err != nil {
		log.Fatal("Could not read file")
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	// This part is really ugly, but gets the job done with converting
	// unstructured JSON to GO.

	maxX, maxY := 0.0, 0.0
	minX, minY := math.MaxFloat64, math.MaxFloat64

	buildings := make([][][]triangolatte.Point, 0)
	for i, f := range m["features"].([]interface{}) {
		buildings = append(buildings, [][]triangolatte.Point{})
		geometry := f.(map[string]interface{})["geometry"].(map[string]interface{})
		switch geometry["type"] {
		case "Polygon":
			for j, polygon := range geometry["coordinates"].([]interface{}) {
				buildings[i] = append(buildings[i], []triangolatte.Point{})
				for _, p := range polygon.([]interface{}) {
					pointArray := p.([]interface{})
					point := triangolatte.Point{
						X: pointArray[0].(float64),
						Y: pointArray[1].(float64),
					}

					pointInMeters := triangolatte.DegreesToMeters(point)

					if pointInMeters.X < minX {
						minX = pointInMeters.X
					}

					if pointInMeters.X > maxX {
						maxX = pointInMeters.X
					}

					if pointInMeters.Y < minY {
						minY = pointInMeters.Y
					}

					if pointInMeters.Y > maxY {
						maxY = pointInMeters.Y
					}

					buildings[i][j] = append(buildings[i][j], pointInMeters)
				}
			}
		case "LineString":
		case "Point":
		}
	}

	for i, p2 := range buildings {
		for j, p1 := range p2 {
			for k, p := range p1 {
				buildings[i][j][k] = triangolatte.Point{X: p.X - minX, Y: p.Y - minY}
			}
		}
	}

	totalSuccesses, totalErrors := 0, 0
	for _, b := range buildings {
		if len(b) == 0 {
			continue
		}

		errored := false
		cleaned, err := triangolatte.EliminateHoles(b)

		if err != nil {
			fmt.Printf("Error in hole removal: %s\n", err)
			errored = true
		}

		triangles, err := triangolatte.EarCut(cleaned)

		if err != nil {
			fmt.Printf("Error in triangulation: %s\n", err)
			errored = true
		}

		var holes [][]triangolatte.Point
		if len(b) > 1 {
			holes = b[1:]
		} else {
			holes = [][]triangolatte.Point{}
		}
		_, _, deviation := triangolatte.Deviation(b[0], holes, triangles)

		if deviation > 1e-10 {
			fmt.Printf("Error detected in deviation\n")
			errored = true
		}

		if errored {
			totalErrors += 1
		} else {
			totalSuccesses += 1
			fmt.Printf("%#v\n", triangles)
		}
	}

	fmt.Printf("success: %d failure: %d", totalSuccesses, totalErrors)
}
