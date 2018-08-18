package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"

	. "github.com/tchayen/triangolatte"
)

// TODO
//
// backend:
// - read *.geojson file
// - triangulate its content
// - join the arrays
// - save as one huge binary blob array buffer whatever
// - serve it as a file
//
// frontend:
// - read array buffer file
// - display it

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

// loadFileToJSON
func loadFileToJSON(fileName string) map[string]interface{} {
	file, err := ioutil.ReadFile("./data.geojson")
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

// parseData
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

// normalize
func normalize(shapes []Shape) []Shape {
	return shapes
}

// triangulate
func triangulate(shapes []Shape) [][]float64 {
	return [][]float64{}
}

// flatten
func flatten(triangles [][]float64) []float64 {
	return []float64{}
}

// saveToFile
func saveToFile(vertices []float64) {

}

// startServer
func startServer() {
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
	check(err)
}

func main() {
	jsonMap := loadFileToJSON("./data.geojson")
	shapes := parseData(jsonMap)
	normalized := normalize(shapes)
	triangulated := triangulate(normalized)
	vertices := flatten(triangulated)
	saveToFile(vertices)
	startServer()
}
