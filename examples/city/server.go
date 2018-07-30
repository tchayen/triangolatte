package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func DownloadFile(filePath string, url string) error {
	// Create the file.
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data.
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file.
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	data, err := ioutil.ReadFile("assets/cracow_tmp")

	if err != nil {
		log.Fatal("Could not read file")
	}

	type Geometry struct {
		coordinates [][][]float64
	}

	type Feature struct {
		geometry Geometry
	}

	type FeatureCollection struct {
		features []Feature
	}

	obj := FeatureCollection{}
	parsed := json.Unmarshal([]byte(data), &obj)

	log.Println(parsed)
}
