package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const chunkSize = 2140

func main() {
	var mode string
	flag.Func("mode", "choose mode to encrypt or decrypt", func(s string) error {
		if s == "encrypt" || s == "decrypt" {
			mode = s
			return nil
		} else {
			return fmt.Errorf("unsupported mode")
		}
	})
	inputPath := flag.String("input", "", "")
	outputPath := flag.String("output", "", "")
	flag.Parse()

	privateKeyBytes, err := readKeyFromFile("private.key")
	if err != nil {
		fmt.Println("Error reading private key:", err)
		return
	}

	publicKeyBytes, err := readKeyFromFile("public.key")
	if err != nil {
		fmt.Println("Error reading public key:", err)
		return
	}

	// Parse the private key
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		fmt.Println("failed to parse private key")
		return
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("failed to parse private key:", err)
		return
	}

	// Parse the public key
	block, _ = pem.Decode(publicKeyBytes)
	if block == nil {
		fmt.Println("failed to parse public key")
		return
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("failed to parse public key:", err)
		return
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		fmt.Println("public key is not RSA")
		return
	}

	privateKeyType, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		fmt.Println("hello")
	}

	if mode == "encrypt" {
		err := EncryptWithPublicKey(publicKey, *inputPath, *outputPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := DecryptWithPrivateKey(privateKeyType, *inputPath, *outputPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func readKeyFromFile(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(pub *rsa.PublicKey, inputPath string, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open file %s: %v", inputPath, err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot create file %s: %v", outputPath, err)
	}
	defer outputFile.Close()

	hash := sha512.New()
	buffer := make([]byte, chunkSize)

	for {
		n, err := inputFile.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("error reading file %s: %v", inputPath, err)
			}
			break
		}

		ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, buffer[:n], nil)
		if err != nil {
			return fmt.Errorf("encryption error: %v", err)
		}

		if _, err := outputFile.Write(ciphertext); err != nil {
			return fmt.Errorf("cannot write encrypted data into file %s: %v", outputPath, err)
		}
	}

	return nil
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(priv *rsa.PrivateKey, inputPath string, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open file %s: %v", inputPath, err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot create file %s: %v", outputPath, err)
	}
	defer outputFile.Close()

	hash := sha512.New()
	buffer := make([]byte, 256) // Assuming a 2048-bit RSA key, the buffer size should be 256 bytes

	for {
		n, err := inputFile.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("error reading file %s: %v", inputPath, err)
			}
			break // End of file
		}

		// Attempt to decrypt the entire buffer read
		plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, buffer[:n], nil)
		if err != nil {
			return fmt.Errorf("decryption error: %v", err)
		}

		if _, err := outputFile.Write(plaintext); err != nil {
			return fmt.Errorf("cannot write decrypted data into file %s: %v", outputPath, err)
		}
	}

	return nil
}
