package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	. "triangolatte"
)

// GPX is the root element in the XML file.
type GPX struct {
	Trk  Trk    `xml:"trk"`
	Time string `xml:"metadata>time"`
}

// Trk represents a track - an ordered list of points describing a path.
type Trk struct {
	Name   string   `xml:"name"`
	Trkseg []Trkseg `xml:"trkseg"`
}

// Trkseg is a Track Segment - it holds a list of Track Points which are
// logically connected in order. To represent a single GPS track where GPS
// reception was lost, or the GPS receiver was turned off, start a new Track
// Segment for each continuous span of track data.
type Trkseg struct {
	Trkpt []Trkpt `xml:"trkpt"`
}

// Trkpt is a Track Point - geographic point with optional elevation and time.
type Trkpt struct {
	Latitude  float64 `xml:"lat,attr"`
	Longitude float64 `xml:"lon,attr"`
	Elevation float64 `xml:"ele"`
	Time      string  `xml:"time"`
}

// XMLToGPX takes byte array from *.xml file and returns parsed GPX.
func XMLToGPX(data []byte) (gpx GPX, err error) {
	err = xml.Unmarshal(data, &gpx)
	return
}

// GPXToPoints takes parsed GPX data and returns array of arrays of points
// (divided into segments as in GPX source).
func GPXToPoints(gpx GPX) [][]Point {
	segmentPoints := make([][]Point, len(gpx.Trk.Trkseg))
	for i := range gpx.Trk.Trkseg {
		segmentPoints[i] = make([]Point, len(gpx.Trk.Trkseg[i].Trkpt))

		for j := range gpx.Trk.Trkseg[i].Trkpt {
			trackPoint := gpx.Trk.Trkseg[i].Trkpt[j]
			lon, lat := trackPoint.Longitude, trackPoint.Latitude
			point := Point{X: lon, Y: lat}

			segmentPoints[i][j] = DegreesToMeters(point)
		}
	}
	return segmentPoints
}

// TriangulatePoints takes array of arrays of points and triangulates them.
func TriangulatePoints(points [][]Point) [][]float64 {
	triangles := make([][]float64, len(points))
	for i := range points {
		triangles[i] = Miter(points[i], 2)
	}
	return triangles
}

func main() {
	// Load data from file.
	data, err := ioutil.ReadFile("assets/gpx_tmp")

	if err != nil {
		log.Fatal("Could not read file")
	}

	gpx, err := XMLToGPX(data)

	if err != nil {
		log.Fatal("Failed to parse GPX file")
	}

	segmentPoints := GPXToPoints(gpx)
	triangles := TriangulatePoints(segmentPoints)

	fmt.Println(triangles)
}
