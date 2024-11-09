package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"flag"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	folder := flag.String("filename", "", "")
	flag.Parse()

	allFiles, err := os.ReadDir(*folder)
	if err != nil {
		fmt.Println("Cannot read this file path")
	}

	for _, file := range allFiles {
		if !file.IsDir() {
			filepath := filepath.Join(*folder, file.Name())
			md5Hash, err := calculateHash(filepath, md5.New)
			if err != nil {
				log.Println("Error calculating MD5 hash for", file, err)
			} else {
				fmt.Printf("MD5 hash of %s: %s\n", file, md5Hash)
			}
			sha1Hash, err := calculateHash(filepath, sha1.New)
			if err != nil {
				log.Println("Error calculating SHA-1 hash for", file, err)
			} else {
				fmt.Printf("SHA-1 hash of %s: %s\n", file, sha1Hash)
			}
			sh256Hash, err := calculateHash(filepath, sha256.New)
			if err != nil {
				log.Println("Error calculating SHA-256 hash for", file, err)
			} else {
				fmt.Printf("SHA-256 hash of %s: %s\n", file, sh256Hash)
			}
		}
	}
	if err != nil {
		log.Fatal(err)
	}
}

func calculateHash(filePath string, hashFunc func() hash.Hash) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := hashFunc()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}

	hashBytes := h.Sum(nil)
	return fmt.Sprintf("%x", hashBytes), nil
}
