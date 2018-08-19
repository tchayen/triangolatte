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

const (
	BUILDINGS int = 0
	PARKS     int = 1
	ROADS     int = 2
)

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
func parseData(m map[string]interface{}) [][]Shape {
	shapes := make([][]Shape, 3)
	shapes[BUILDINGS] = make([]Shape, 0)
	shapes[PARKS] = make([]Shape, 0)
	shapes[ROADS] = make([]Shape, 0)

	for _, f := range m["features"].([]interface{}) {
		feature := f.(map[string]interface{})
		geometry := feature["geometry"].(map[string]interface{})

		switch geometry["type"] {
		case "Polygon":
			var s *Shape
			properties := feature["properties"].(map[string]interface{})

			if properties["building"] == "yes" {
				shapes[BUILDINGS] = append(shapes[BUILDINGS], Shape{})
				s = &shapes[BUILDINGS][len(shapes[BUILDINGS])-1]
			} else if properties["leisure"] == "park" {
				shapes[PARKS] = append(shapes[PARKS], Shape{})
				s = &shapes[PARKS][len(shapes[PARKS])-1]
			} else {
				continue
			}

			s.Type = "Polygon"

			for j, polygon := range geometry["coordinates"].([]interface{}) {
				s.Points = append(s.Points, []Point{})
				for _, p := range polygon.([]interface{}) {
					// Cast from interface{} to []interface{} is necessary.
					pointArray := p.([]interface{})

					point := Point{
						X: pointArray[0].(float64),
						Y: pointArray[1].(float64),
					}

					// Convert coordinates.
					pointInMeters := degreesToMeters(point)
					s.Points[j] = append(s.Points[j], pointInMeters)
				}
			}
		case "LineString":
			var s *Shape
			properties := feature["properties"].(map[string]interface{})

			if properties["highway"] != nil {
				shapes[ROADS] = append(shapes[ROADS], Shape{})
				s = &shapes[ROADS][len(shapes[ROADS])-1]
			} else {
				continue
			}

			s.Type = "LineString"
			s.Points = make([][]Point, 1)
			for _, p := range geometry["coordinates"].([]interface{}) {
				// Cast from interface{} to []interface{} is necessary.
				pointArray := p.([]interface{})

				point := Point{
					X: pointArray[0].(float64),
					Y: pointArray[1].(float64),
				}

				// Convert coordinates.
				pointInMeters := degreesToMeters(point)
				s.Points[0] = append(s.Points[0], pointInMeters)
			}
		}
	}
	return shapes
}

// findGlobalBounds finds min and max coordinates.
func findGlobalBounds(shapes [][]Shape) (xMin, yMin, xMax, yMax float64) {
	xMin, yMin, xMax, yMax = math.MaxFloat64, math.MaxFloat64, 0.0, 0.0
	for _, group := range shapes {
		for _, s := range group {
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
	}
	return
}

// normalize puts all points in range [-1, 1], where the longer axis (either X
// or Y) is stretched to the whole range and the other is translated in a way
// that keeps proportions.
func normalize(shapes [][]Shape) [][]Shape {
	xMin, yMin, xMax, yMax := findGlobalBounds(shapes)
	upper := math.Max(xMax-xMin, yMax-yMin)
	for _, group := range shapes {
		for k := range group {
			for i := range group[k].Points {
				for j := range group[k].Points[i] {
					p := group[k].Points[i][j]
					group[k].Points[i][j] = Point{
						X: (p.X - xMin) / upper,
						Y: (p.Y - yMin) / upper,
					}
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
			triangles[i], err = Line(Miter, s.Points[0], 0.003)
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
func saveToFile(fileName string, vertices []float32) {
	f, err := os.Create(fmt.Sprintf("%s_tmp", fileName))
	check(err)
	w := bufio.NewWriter(f)
	err = binary.Write(w, binary.LittleEndian, vertices)
	w.Flush()
	f.Close()
}

// processShapes is for the last phase when each set of shapes can be finished
// separately.
func processShapes(fileName string, shapes []Shape) {
	triangulated := triangulate(shapes)
	vertices := flatten(triangulated)
	converted := toFloat32(vertices)
	saveToFile(fileName, converted)
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
	jsonMap := loadFileToJSON("./export.geojson")
	shapes := parseData(jsonMap)
	normalized := normalize(shapes)

	names := []string{"buildings", "parks", "roads"}
	for i := 0; i < 3; i++ {
		processShapes(names[i], normalized[i])
	}

	startServer()
}
