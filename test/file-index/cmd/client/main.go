package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"training/file-index/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair("cert/client-cert.pem", "cert/client-key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("Dial server %s", *serverAddress)

	tlsCredential, err := loadTLSCredentials()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(tlsCredential))
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	fileSearcherClient := pb.NewFileIndexClient(conn)

	searchFile(fileSearcherClient)
}

func searchFile(fileClient pb.FileIndexClient) {
	ctx := context.Background()

	req := &pb.CreateFileDiscoverRequest{
		Request: "Start search file of system",
	}
	stream, err := fileClient.ListFiles(ctx, req)
	if err != nil {
		log.Fatal("cannot search file: ", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}

		if err != nil {
			log.Fatal("cannot receive response: ", err)
		}

		file := res.GetFiles()
		log.Printf("- found: %s", file)
	}
}
