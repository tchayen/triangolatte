package main

import (
	"io"
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
	url := "https://overpass-api.de/api/interpreter?data=[out:json];(way[building](50.0,19.85,50.105,20.13);relation[building](50.0,19.85,50.105,20.13););out body;>;out skel qt;"

	err := DownloadFile("cracow_tmp", url)
	if err != nil {
		log.Fatal("Error downloading file")
	}
}
