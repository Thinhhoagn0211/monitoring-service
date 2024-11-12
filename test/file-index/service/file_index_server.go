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
	"syscall"
	"time"
	"training/file-index/pb"
	"training/file-index/serializer"
	"training/file-index/util"

	"github.com/fumiama/go-docx"
	"github.com/xuri/excelize/v2"
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
		drives := util.GetAvailableDrives()
		for _, drive := range drives {
			currentFiles := make(FileInfoMap)
			// Adjust the directory path as needed (for example, from `request` if specified)
			err := filepath.Walk(drive, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					if os.IsPermission(err) {
						fmt.Printf("Permission denied: %s\n", path)
						return nil
					}
					return err
				}
				// Skip directories
				if info.IsDir() {
					return nil
				}
				ext := filepath.Ext(path)

				// Get file timestamps
				createdAt, modifiedAt, accessedAt, _ := getFileTimes(path)

				var content string
				if ext == ".docx" {
					content, err = readDocxContent(path)
					if err != nil {
						fmt.Printf("Failed to read .docx file: %s, error: %v\n", path, err)
						return nil // Skip to the next file
					}
					fmt.Println("content", content)
				} else if ext == ".xlsx" {
					content, err = readExcelContent(path)
					if err != nil {
						fmt.Printf("Failed to read .xlsx file: %s, error: %v\n", path, err)
						return nil // Skip to the next file
					}
					fmt.Println("content", content)
				} else if ext == ".pptx" {
					content, err = readPPTXContent(path)
					if err != nil {
						fmt.Printf("Failed to read .pptx file: %s, error: %v\n", path, err)
						return nil // Skip to the next file
					}
					fmt.Println("content", content)
				}
				var attribute string
				flagAttributes, _ := IsHiddenFile(path)
				if flagAttributes {
					attribute = "Hidden"
				} else {
					attribute = "Read Only"
				}
				fileAttr := &pb.FileAttr{
					Path:       path,
					Name:       info.Name(),
					Type:       ext,
					Size:       info.Size(),
					CreatedAt:  timestamppb.New(createdAt),
					ModifiedAt: timestamppb.New(modifiedAt),
					AccessedAt: timestamppb.New(accessedAt),
					Attributes: attribute,
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
		}

		time.Sleep(5 * time.Second) // Quét lại sau 5 giây
	}
}

func IsHiddenFile(filename string) (bool, error) {
	pointer, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return false, err
	}
	attributes, err := syscall.GetFileAttributes(pointer)
	if err != nil {
		return false, err
	}
	return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0, nil
}

func readDocxContent(filePath string) (string, error) {
	readFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	fileinfo, err := readFile.Stat()
	if err != nil {
		return "", err
	}

	var content string
	size := fileinfo.Size()
	doc, err := docx.Parse(readFile, size)
	if err != nil {
		return "", err
	}
	for _, it := range doc.Document.Body.Items {
		switch it.(type) {
		case *docx.Paragraph, *docx.Table: // printable
			switch v := it.(type) {
			case *docx.Paragraph:
				content += v.String() + "\n"
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
	defer f.Close()
	var contentBuilder strings.Builder
	for _, sheetName := range f.GetSheetMap() {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return "", fmt.Errorf("failed to get rows: %w", err)
		}
		for _, row := range rows {
			for _, col := range row {
				contentBuilder.WriteString(col)
				contentBuilder.WriteString("\t")
			}
			contentBuilder.WriteString("\n")
		}
	}

	content := contentBuilder.String()
	return content, nil
}

func readPPTXContent(filePath string) (string, error) {
	// Convert .pptx file to text
	// res, err := docconv.ConvertPath(filePath)
	// if err != nil {
	// 	return "", err
	// }
	// content := res.Body
	return "content", nil
}

func getFileTimes(path string) (createdAt, modifiedAt, accessedAt time.Time, err error) {
	// Get file information
	fileInfo, err := os.Stat(path)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, err
	}
	stat := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	createdAt = time.Unix(0, stat.CreationTime.Nanoseconds())
	accessedAt = time.Unix(0, stat.LastAccessTime.Nanoseconds())
	modifiedAt = fileInfo.ModTime()

	return createdAt, modifiedAt, accessedAt, nil
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
