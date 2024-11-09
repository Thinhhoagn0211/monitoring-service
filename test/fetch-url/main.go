package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Custom type to handle multiple URLs
type stringSlice []string

// Implement the String method for the flag
func (s *stringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

// Implement the Set method to handle adding URLs
func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

var urls stringSlice

func init() {
	flag.Var(&urls, "url", "A flag that accepts multiple values")
}

func main() {

	flag.Parse()

	if len(urls) == 0 {
		fmt.Println("Errors: No URLs provided")
		os.Exit(1)
	}

	fmt.Println("Fetching the following URLs:", urls)
	for _, url := range urls {
		fmt.Printf(" - %s\n", url)
		err := fetchContentUrl(url)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func fetchContentUrl(url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("can not read response ", err)
	}
	urlName := strings.Split(url, "/")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Status is not be OK")
	}
	out, err := os.Create(urlName[len(urlName)-1] + ".txt")
	if err != nil {
		return fmt.Errorf("Can not create a file")
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Cannot copy content into file %s", urlName)
	}
	return nil
}
