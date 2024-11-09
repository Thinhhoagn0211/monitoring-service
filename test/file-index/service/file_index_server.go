package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"training/file-index/pb"
	"training/file-index/serializer"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FileDiscoveryServer struct {
	pb.UnimplementedFileIndexServer
	fileStore FileStore
}
type FileInfoMap map[string]*pb.FileAttr

var fileInfoCache FileInfoMap

func NewFileDiscoveryServer(fileStore FileStore) *FileDiscoveryServer {
	return &FileDiscoveryServer{
		fileStore: fileStore,
	}
}

func (server *FileDiscoveryServer) GetCheckSumFiles(ctx context.Context, req *pb.CreateFileChecksumRequest) (*pb.CreateFileChecksumResponse, error) {
	filepaths := req.GetFilepath()
	var res = &pb.CreateFileChecksumResponse{}
	for _, filepath := range filepaths {
		filename := strings.Split(filepath, "/")
		log.Printf("receive a filepath checksum request with path :%s\n", filename[len(filename)-1])

		fileChecksum, err := serializer.CalculateHash(filepath, md5.New)
		if err != nil {
			return nil, err
		}

		if err := logError(err); err != nil {
			return nil, err
		}
		if err := contextError(ctx); err != nil {
			return nil, err
		}
		res = &pb.CreateFileChecksumResponse{
			Checksums: map[string]string{
				filename[len(filename)-1]: fileChecksum,
			},
		}
	}
	return res, nil
}

func (server *FileDiscoveryServer) ListFiles(req *pb.CreateFileDiscoverRequest, stream grpc.ServerStreamingServer[pb.CreateFileDiscoverResponse]) error {
	request := req.GetRequest()
	fmt.Printf("Receive request to list all files in computer %s\n", request)

	for {
		currentFiles := make(FileInfoMap)
		err := filepath.Walk("C:\\Users\\Raven\\Work", func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			// Lấy thông tin tập tin
			ext := filepath.Ext(path)
			stat := info.Sys().(*syscall.Win32FileAttributeData)

			createdAt := timestamppb.New(filetimeToTime(stat.CreationTime))
			modifiedAt := timestamppb.New(filetimeToTime(stat.LastWriteTime))
			accessedAt := timestamppb.New(filetimeToTime(stat.LastAccessTime))

			fileAttr := &pb.FileAttr{
				Id:         uuid.New().String(),
				Path:       path,
				Name:       info.Name(),
				Type:       ext,
				Size:       info.Size(),
				CreatedAt:  createdAt,
				ModifiedAt: modifiedAt,
				AccessedAt: accessedAt,
				Content:    "",
			}
			currentFiles[path] = fileAttr
			// Lưu vào CSDL
			if _, exists := fileInfoCache[path]; !exists {
				// Tập tin mới
				fmt.Println("Create file")
				if err := server.fileStore.Save(fileAttr); err != nil {
					return err
				}
			}

			res := &pb.CreateFileDiscoverResponse{
				Files: fileAttr,
			}
			if err := stream.Send(res); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}

		// Kiểm tra tập tin đã bị xóa
		for path, fileAttr := range fileInfoCache {
			if _, exists := currentFiles[path]; !exists {
				// Tập tin đã bị xóa
				fmt.Println("delete path ", fileAttr)
				if err := server.fileStore.Delete(fileAttr.Name); err != nil {
					return err
				}
			}
		}

		// Cập nhật cache
		fileInfoCache = currentFiles

		time.Sleep(5 * time.Second) // Quét lại sau 5 giây
	}
}

// Convert syscall.Filetime to time.Time
func filetimeToTime(ft syscall.Filetime) time.Time {
	return time.Unix(0, ft.Nanoseconds()).UTC()
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
