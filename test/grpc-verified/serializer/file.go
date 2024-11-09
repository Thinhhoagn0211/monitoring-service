package serializer

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func DownloadFile(url string) (string, error) {
	filePath := strings.Split(url, "/")
	// Create the file
	out, err := os.Create(filePath[len(filePath)-1])
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return filePath[len(filePath)-1], nil
}
