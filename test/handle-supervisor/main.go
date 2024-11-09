package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	filename := "text.txt"
	content := "hello1"
	var f *os.File
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Create a File")
		f, err = os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		contentReader, err := os.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", contentReader)
	}
	defer f.Close()

	os.WriteFile(filename, []byte(content), 0644)
}
