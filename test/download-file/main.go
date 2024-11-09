package main

import (
	"flag"
	"io"
	"net/http"
	"os"
)

func main() {
	url := flag.String("url", "", "")
	filename := flag.String("filename", "", "")
	flag.Parse()
	err := DownloadFile(*url, *filename)
	if err != nil {
		panic(err)
	}
}

/*
DownloadFile to create a filepath with filename and get file from url and copy into filepath
*/
func DownloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
