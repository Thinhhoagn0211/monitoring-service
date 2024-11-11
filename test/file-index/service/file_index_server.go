package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"training/file-index/pb"
	"training/file-index/serializer"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/fumiama/go-docx"
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
		// Adjust the directory path as needed (for example, from `request` if specified)
		err := filepath.Walk("/home/thinh/Desktop/test", func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Skip directories
			if info.IsDir() {
				return nil
			}
			ext := filepath.Ext(path)

			// Get file timestamps
			createdAt, modifiedAt, accessedAt := getFileTimes(info)

			var content string
			if ext == ".docx" {
				content, _ = readDocxContent(path)
				fmt.Println("content", content)
			} else if ext == ".xlsx" {
				content, _ = readExcelContent(path)
				fmt.Println("content", content)
			} else if ext == ".pptx" {
				content, _ = readPPTXContent(path)
				fmt.Println("content", content)
			}

			fileAttr := &pb.FileAttr{
				Path:       path,
				Name:       info.Name(),
				Type:       ext,
				Size:       info.Size(),
				CreatedAt:  timestamppb.New(createdAt),
				ModifiedAt: timestamppb.New(modifiedAt),
				AccessedAt: timestamppb.New(accessedAt),
				Content:    content,
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

func readDocxContent(filePath string) (string, error) {
	readFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	fileinfo, err := readFile.Stat()
	if err != nil {
		panic(err)
	}
	var content string
	size := fileinfo.Size()
	doc, err := docx.Parse(readFile, size)
	if err != nil {
		panic(err)
	}
	fmt.Println("Plain text:")
	for _, it := range doc.Document.Body.Items {
		switch it.(type) {
		case *docx.Paragraph, *docx.Table: // printable
			switch v := it.(type) {
			case *docx.Paragraph:
				content += v.String()
			case *docx.Table:
				content += v.String()
			}
		}
	}
	return content, nil
}

func readExcelContent(filePath string) (string, error) {
	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}

	// Extract content from the first sheet
	var content string
	for _, sheetName := range f.GetSheetMap() {
		rows := f.GetRows(sheetName)
		for _, row := range rows {
			for _, cell := range row {
				content += cell + " "
			}
			content += "\n"
		}
	}

	return content, nil
}

func readPPTXContent(filePath string) (string, error) {
	// Open the PowerPoint file
	var content string

	return content, nil
}

// Helper function to get created time, modified time, and accessed time for a file
func getFileTimes(info os.FileInfo) (createdAt, modifiedAt, accessedAt time.Time) {
	modifiedAt = info.ModTime() // Modified time from FileInfo

	// Assume createdAt and accessedAt are set to modified time as placeholders
	createdAt, accessedAt = modifiedAt, modifiedAt

	// You may use OS-specific logic here to retrieve exact created and accessed times
	return createdAt, modifiedAt, accessedAt
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
