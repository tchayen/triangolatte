package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	// data, err := ioutil.ReadFile("assets/cracow_tmp")

	// if err != nil {
	// 	log.Fatal("Could not read file")
	// }

	data := []byte("{\"features\":[{\"geometry\":{\"coordinates\":[[[0,100],[100,100]]]}}]}")
	type Geometry struct {
		coordinates [][][]float64
	}

	type Feature struct {
		geometry Geometry
	}

	type FeatureCollection struct {
		features []Feature
	}

	d := FeatureCollection{[]Feature{{Geometry{[][][]float64{}}}}}
	a, _ := json.Marshal(&d)
	fmt.Println(string(a))

	obj := FeatureCollection{}
	parsed := json.Unmarshal(data, &obj)
	log.Println(parsed)
}
