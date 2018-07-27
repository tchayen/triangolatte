package triangolatte

import (
	"encoding/json"
	"io/ioutil"
	"math"
)

// polygonArea calculates real area of the polygon.
func polygonArea(data []Point) float64 {
	area := 0.0
	for i, j := 0, len(data)-1; i < len(data); i++ {
		area += data[i].X*data[j].Y - data[i].Y*data[j].X
		j = i
	}
	return math.Abs(area / 2)
}

// trianglesArea calculates summed area of all triangles.
func trianglesArea(t []float64) float64 {
	trianglesArea := 0.0
	for i := 0; i < len(t); i += 6 {
		trianglesArea += math.Abs((t[i]*(t[i+3]-t[i+5]) + t[i+2]*(t[i+5]-t[i+1]) + t[i+4]*(t[i+1]-t[i+3])) / 2)
	}
	return trianglesArea
}

// deviation calculates difference between real area and the one from triangulation.
func deviation(data []Point, t []float64) (real, triangles, deviation float64) {
	triangles = trianglesArea(t)
	real = polygonArea(data)
	deviation = math.Abs(triangles - real)
	return
}

// loadPointsFromFile takes file name and returns array of points.
func loadPointsFromFile(fileName string) ([][]Point, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	polygons := make([][][]float64, 0)
	json.Unmarshal([]byte(data), &polygons)

	points := make([][]Point, len(polygons))
	for i := range polygons {
		points[i] = make([]Point, len(polygons[i]))
		for j := range polygons[i] {
			points[i][j] = Point{polygons[i][j][0], polygons[i][j][1]}
		}
	}
	return points, nil
}

// Origin shift comes from the circumference of the Earth in meters (6378137).
const originShift = 2.0 * math.Pi * 6378137 / 2.0

// Converts longitude and latitude using WGS84 Geodetic Datum to meters with
// Spherical Mercator projection, known officially under EPSG:3857 codename.
//
// X is longitude, Y is latitude.
//
// Bounds: `[-180.0, -85.06, 180.0, 85.06]`.
func degreesToMeters(point Point) Point {
	return Point{
		point.X * originShift / 180.0,
		math.Log(math.Tan((90.0+point.Y)*math.Pi/360.0)) / (math.Pi / 180.0) * originShift / 180.0,
	}
}