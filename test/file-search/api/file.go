package api

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"
	"training/file-index/pb"
	db "training/file-search/db/sqlc"
	"training/file-search/util"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type newDataFile struct {
	FilePath string `json:"filepath"`
	CheckSum string `json:"checksum"`
}
type Metadata struct {
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
type newFileResponse struct {
	Status int           `json:"status"`
	Data   []newDataFile `json:"data"`
	Meta   Metadata      `json:"meta"`
}

func (server *Server) getFileSearcher(ctx *gin.Context) {
	// Retrieve query parameters with defaults
	name := ctx.Param("id")
	extension := ctx.DefaultQuery("extension", "")
	sizeMin, _ := strconv.Atoi(ctx.DefaultQuery("size_min", "0"))
	sizeMax, _ := strconv.Atoi(ctx.DefaultQuery("size_max", "1000000000000"))
	createdAfter, _ := time.Parse(time.RFC3339, ctx.DefaultQuery("created_after", ""))
	createdBefore, _ := time.Parse(time.RFC3339, ctx.DefaultQuery("created_before", "1920652419"))
	// modifiedAfter, _ := time.Parse(time.RFC3339, ctx.DefaultQuery("modified_after", ""))
	// modifiedBefore, _ := time.Parse(time.RFC3339, ctx.DefaultQuery("modified_before", "1920652419"))
	// accessedAfter, _ := time.Parse(time.RFC3339, ctx.DefaultQuery("accessed_after", ""))
	// accessedBefore, _ := time.Parse(time.RFC3339, ctx.DefaultQuery("accessed_before", "1920652419"))
	// content := ctx.DefaultQuery("content", "")

	// Prepare the database parameters
	arg := db.GetFileParams{
		Column1:   util.NullableString(name),
		Extension: extension,
		Size:      int64(sizeMin),
		Size_2:    int64(sizeMax),
		// CreatedAt:   createdAfter,
		// CreatedAt_2: createdBefore,
		// ModifiedAt:   modifiedAfter,
		// ModifiedAt_2: modifiedBefore,
		// AccessedAt:   accessedAfter,
		// AccessedAt_2: accessedBefore,
		// Column11:     util.NullableString(content),
	}
	// fmt.Println("params ", arg)

	// Call the database to retrieve the files
	paths, err := server.store.GetFile(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	req := &pb.CreateFileChecksumRequest{
		Filepath: paths,
	}
	var data []newDataFile
	for _, path := range paths {
		res, err := checksumFile(server.fileSearcherClient, req)
		if err != nil {
			log.Fatal("error: ", err)
		}
		data = append(data, newDataFile{
			FilePath: path,
			CheckSum: res.Checksums[path],
		})
	}
	res := newFileResponse{
		Status: http.StatusOK,
		Data:   data,
		Meta: Metadata{
			Total:  len(paths),
			Offset: 0,
			Limit:  10,
		},
	}
	// Return the search results
	ctx.JSON(http.StatusOK, res)
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
