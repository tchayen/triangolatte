package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	. "triangolatte"
)

// Building collects meta data about building and its points.
type Building struct {
	Properties map[string]string
	Points     [][]Point
}

type Triangulated struct {
	Properties map[string]string `json:"properties"`
	Triangles  []float64         `json:"triangles"`
}

// parseData takes JSON naively to map[string]interface{} and returns more
// organized []Building array.
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
				buildings[i].Points = append(buildings[i].Points, []Point{})

				for _, p := range polygon.([]interface{}) {
					// Cast from interface{} to []interface{}.
					pointArray := p.([]interface{})

					point := Point{
						X: pointArray[0].(float64),
						Y: pointArray[1].(float64),
					}

					// Convert coordinates.
					pointInMeters := DegreesToMeters(point)
					buildings[i].Points[j] = append(buildings[i].Points[j], pointInMeters)
				}
			}
		case "LineString":
		case "Point":
		}
	}
	return
}

// findMinMax takes array of points and finds min and max coordinates.
func findMinMax(points []Point) (xMin, yMin, xMax, yMax float64) {
	xMin, yMin, xMax, yMax = math.MaxFloat64, math.MaxFloat64, 0.0, 0.0
	for _, p := range points {
		if p.X < xMin {
			xMin = p.X
		}

		if p.X > xMax {
			xMax = p.X
		}

		if p.Y < yMin {
			yMin = p.Y
		}

		if p.Y > yMax {
			yMax = p.Y
		}
	}
	return
}

// normalizeCoordinates takes building coordinates and changes them to fit in
// range [0.0, 1.0].
func normalizeCoordinates(buildings []Building) {
	for _, b := range buildings {
		if len(b.Points) == 0 {
			continue
		}

		xMin, yMin, xMax, yMax := findMinMax(b.Points[0])
		for i := range b.Points {
			for j := range b.Points[i] {
				v := b.Points[i][j]
				b.Points[i][j] = Point{
					X: (v.X - xMin) / (xMax - xMin),
					Y: (v.Y - yMin) / (yMax - yMin),
				}
			}
		}
	}
}

// triangulate takes building coordinates and triangulates them resulting in
// array of floats and sums of total errors and successes as a side effect.
func triangulate(buildings []Building) (
	triangulated []Triangulated,
	totalSuccesses int,
	totalErrors int,
) {
	triangulated = make([]Triangulated, len(buildings))

	for i, b := range buildings {
		if len(b.Points) == 0 {
			continue
		}

		errorHappened := false
		cleaned, err := JoinHoles(b.Points)

		if err != nil {
			errorHappened = true
		}

		t, err := EarCut(cleaned)

		if err != nil {
			errorHappened = true
		}

		var h [][]Point
		if len(b.Points) > 1 {
			h = b.Points[1:]
		} else {
			h = [][]Point{}
		}
		_, _, deviation := Deviation(b.Points[0], h, t)

		triangulated[i] = Triangulated{b.Properties, t}

		// Chosen arbitrarily as a frontier between low and high error rate.
		if deviation > 1e-15 {
			errorHappened = true
		}

		if errorHappened {
			totalErrors++
		} else {
			totalSuccesses++
		}
	}
	return
}

func main() {
	data, err := ioutil.ReadFile("../../assets/cracow_tmp")

	if err != nil {
		log.Fatal("Could not read file")
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)
	buildings := parseData(m)
	normalizeCoordinates(buildings)
	triangulated, successes, errors := triangulate(buildings)
	converted, err := json.Marshal(triangulated)

	if err != nil {
		log.Fatalf("Could not marshal to JSON: %s", err)
	}

	err = ioutil.WriteFile("../../assets/json_tmp", converted, 0644)

	if err != nil {
		log.Fatalf("Could not save file: %s", err)
	}

	fmt.Printf("success: %d failure: %d", successes, errors)
}
