package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"training/grpc-verified/pb"

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

func main() {
	serverAddress := flag.String("address", "", "the server address")
	fileUrl := flag.String("url", "", "url to download")
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

	fileDownloaderClient := pb.NewDownloadFileClient(conn)
	req := &pb.CreateDownloadRequest{
		FileUrl: *fileUrl,
	}

	res, err := fileDownloaderClient.CreateDownload(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Print("file already exists")
		} else {
			log.Fatal("cannot download file: ", err)
		}
	}
	log.Printf("download file from url %s and name file is %s", *fileUrl, res.FilePath)
}
