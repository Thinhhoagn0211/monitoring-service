package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	db "training/file-search/db/sqlc"
	"training/file-search/util"

	"training/file-index/pb"
	"training/file-search/token"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config             util.Config
	store              db.Store
	tokenMaker         token.Maker
	router             *gin.Engine
	fileSearcherClient pb.FileIndexClient
}

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

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	tlsCredential, err := loadTLSCredentials()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(config.GRPCServerAddress, grpc.WithTransportCredentials(tlsCredential))
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	server.fileSearcherClient = pb.NewFileIndexClient(conn)

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/login", server.loginUser)
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.GET("/file/:id", server.getFileSearcher)
	authRoutes.GET("/users", server.getUsers)
	authRoutes.POST("/users", server.createUser)
	authRoutes.PATCH("/users", server.updateUser)
	authRoutes.DELETE("/users/:id", server.deleteUser)
	authRoutes.GET("/users/:id", server.getUserById)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
