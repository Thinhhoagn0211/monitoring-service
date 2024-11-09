package service

import (
	"context"
	"log"
	"training/grpc-verified/pb"
	"training/grpc-verified/serializer"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FileDownloadServer struct {
	pb.UnimplementedDownloadFileServer
}

func NewFileDownLoadServer() *FileDownloadServer {
	return &FileDownloadServer{}
}

func (server *FileDownloadServer) CreateDownload(ctx context.Context, req *pb.CreateDownloadRequest) (*pb.CreateDownloadResponse, error) {
	url := req.GetFileUrl()
	log.Printf("receive a file download request with path :%s\n", url)

	filePath, err := serializer.DownloadFile(url)

	if err := logError(err); err != nil {
		return nil, err
	}
	if err := contextError(ctx); err != nil {
		return nil, err
	}

	res := &pb.CreateDownloadResponse{
		FilePath: filePath,
	}
	return res, nil
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "request is canceled"))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "deadline is exceeded"))
	default:
		return nil
	}
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
