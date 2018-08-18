package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"

	. "github.com/tchayen/triangolatte"
)

// Shape joins type info about geometry feature and its points.
type Shape struct {
	Type   string
	Points [][]Point
}

// check does log.Fatal() with the error message (if any).
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// loadFileToJSON takes file name and returns JSON as a map[string]interface{}.
func loadFileToJSON(fileName string) map[string]interface{} {
	file, err := ioutil.ReadFile(fileName)
	check(err)

	var m map[string]interface{}
	json.Unmarshal(file, &m)

	return m
}

// Origin shift comes from the circumference of the Earth in meters (6378137).
const originShift = 2.0 * math.Pi * 6378137 / 2.0

// degreesToMeters converts longitude and latitude using WGS84 Geodetic Datum to
// meters with Spherical Mercator projection, known officially under EPSG:3857
// codename.
//
// X is longitude, Y is latitude.
func degreesToMeters(point Point) Point {
	return Point{
		X: point.X * originShift / 180.0,
		Y: math.Log(math.Tan((90.0+point.Y)*math.Pi/360.0)) / (math.Pi / 180.0) * originShift / 180.0,
	}
}

// parseData takes map created from JSON file and returns array of shapes.
func parseData(m map[string]interface{}) []Shape {
	shapes := make([]Shape, 0)
	for i, f := range m["features"].([]interface{}) {
		shapes = append(shapes, Shape{})
		feature := f.(map[string]interface{})
		geometry := feature["geometry"].(map[string]interface{})

		switch geometry["type"] {
		case "Polygon":
			shapes[i].Type = "Polygon"

			for j, polygon := range geometry["coordinates"].([]interface{}) {
				shapes[i].Points = append(shapes[i].Points, []Point{})

				for _, p := range polygon.([]interface{}) {
					// Cast from interface{} to []interface{} is necessary.
					pointArray := p.([]interface{})

					point := Point{
						X: pointArray[0].(float64),
						Y: pointArray[1].(float64),
					}

					// Convert coordinates.
					pointInMeters := degreesToMeters(point)
					shapes[i].Points[j] = append(shapes[i].Points[j], pointInMeters)
				}
			}
		case "LineString":
			shapes[i].Type = "LineString"
		case "Point":
			shapes[i].Type = "Point"
		}
	}
	return shapes
}

// findGlobalBounds takes array of shapes and finds min and max coordinates.
func findGlobalBounds(shapes []Shape) (xMin, yMin, xMax, yMax float64) {
	xMin, yMin, xMax, yMax = math.MaxFloat64, math.MaxFloat64, 0.0, 0.0
	for _, s := range shapes {
		for i := range s.Points {
			for _, p := range s.Points[i] {
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
		}
	}
	return
}

// normalize puts all points in range [-1, 1], where the longer axis (either X
// or Y) is stretched to the whole range and the other is translated in a way
// that keeps proportions.
func normalize(shapes []Shape) []Shape {
	for k := range shapes {
		xMin, yMin, xMax, yMax := findGlobalBounds(shapes)
		min, max := math.Min(xMin, yMin), math.Max(xMax, yMax)
		for i := range shapes[k].Points {
			for j := range shapes[k].Points[i] {
				p := shapes[k].Points[i][j]
				shapes[k].Points[i][j] = Point{
					X: (p.X - min) / (max - min),
					Y: (p.Y - min) / (max - min),
				}
			}
		}
	}
	return shapes
}

// triangulate takes array of shapes and returns them triangulated based on
// their type.
func triangulate(shapes []Shape) [][]float64 {
	triangles := make([][]float64, len(shapes))
	var err error

	for i, s := range shapes {
		switch s.Type {
		case "Polygon":
			joined, err := JoinHoles(s.Points)
			check(err)
			triangles[i], err = Polygon(joined)
		case "LineString":
			triangles[i], err = Line(Miter, s.Points[0], 2)
		}
		check(err)
	}
	return triangles
}

// flatten joins array of []float64 into one array.
func flatten(triangles [][]float64) []float64 {
	size := 0
	for i := range triangles {
		size += len(triangles[i])
	}

	flattened := make([]float64, size)
	offset := 0
	for i := range triangles {
		for j := range triangles[i] {
			flattened[offset+j] = triangles[i][j]
		}
		offset += len(triangles[i])
	}

	return flattened
}

// toFloat32 takes array of float64 elements and changes them to float32.
func toFloat32(array []float64) []float32 {
	converted := make([]float32, len(array))
	for i, v := range array {
		converted[i] = float32(v)
	}
	return converted
}

// saveToFile takes simple []float64 array and saves it to a binary file.
func saveToFile(vertices []float32) {
	f, err := os.Create("data_tmp")
	check(err)
	w := bufio.NewWriter(f)
	err = binary.Write(w, binary.LittleEndian, vertices)
	w.Flush()
	f.Close()
}

// startServer starts HTTP server responding to API calls.
func startServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fileName := r.RequestURI[1:]

		if r.RequestURI == "/" {
			fileName = "index.html"
		}

		log.Printf("opening %s", fileName)
		http.ServeFile(w, r, fileName)
	})

	port := 3010
	log.Printf("Listening on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	check(err)
}

func main() {
	// jsonMap := loadFileToJSON("data.geojson")
	// shapes := parseData(jsonMap)
	shapes := []Shape{{
		Type: "Polygon",
		Points: [][]Point{
			{{X: 50, Y: 110}, {X: 150, Y: 30}, {X: 240, Y: 115}, {X: 320, Y: 65}, {X: 395, Y: 170}, {X: 305, Y: 160}, {X: 265, Y: 240}, {X: 190, Y: 100}, {X: 95, Y: 125}, {X: 100, Y: 215}},
		},
	}}
	normalized := normalize(shapes)
	triangulated := triangulate(normalized)
	vertices := flatten(triangulated)
	converted := toFloat32(vertices)
	saveToFile(converted)
	startServer()
}
