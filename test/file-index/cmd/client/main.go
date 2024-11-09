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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
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

type arrayFlags []string

// String is an implementation of the flag.Value interface
func (i *arrayFlags) String() string {
	return fmt.Sprintf("%v", *i)
}

// Set is an implementation of the flag.Value interface
func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var fileUrl arrayFlags

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Var(&fileUrl, "paths", "Some description for this param.")
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
	req := &pb.CreateFileChecksumRequest{
		Filepath: fileUrl,
	}
	go checksumFile(fileSearcherClient, req)

	searchFile(fileSearcherClient)

	// log.Printf("checksum file from url %s and name file is %s", fileUrl, res.Checksums)
}

func checksumFile(fileSearcherClient pb.FileIndexClient, req *pb.CreateFileChecksumRequest) (*pb.CreateFileChecksumResponse, error) {
	res, err := fileSearcherClient.GetCheckSumFiles(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Print("file already exists")
		} else {
			log.Fatal("cannot checksum file: ", err)
		}
	}
	return res, err
}

func searchFile(fileClient pb.FileIndexClient) {
	ctx := context.Background()

	req := &pb.CreateFileDiscoverRequest{
		Request: "hello",
	}
	stream, err := fileClient.ListFiles(ctx, req)
	if err != nil {
		log.Fatal("cannot search laptop: ", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}

		if err != nil {
			log.Fatal("cannot receive response: %v", err)
		}

		res.GetFiles()
		// log.Printf("- found: %s", file)
	}
}
