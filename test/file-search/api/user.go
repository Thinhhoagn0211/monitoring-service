package api

import (
	db "filesearch/db/sqlc"
	"filesearch/model"
	"filesearch/util"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type updateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type Meta struct {
	Total  int `json:"metadata"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type Response struct {
	Status int       `json:"status"`
	Data   []db.User `json:"data"`
	Meta   Meta      `json:"meta"`
}
type userResponse struct {
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		ID:           int32(util.GenerateRandomID()),
		Email:        req.Email,
		Username:     req.Username,
		Password:     req.Password,
		PasswordHash: hashedPassword,
		Phone:        req.Phone,
		Fullname:     req.FullName,
		Avatar:       req.Avatar,
		State:        int64(model.UserStateActive),
		Role:         req.Role,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Println(user)
	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:  user.Username,
		FullName:  user.Fullname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (server *Server) getUsers(ctx *gin.Context) {
	searchTerm := ctx.DefaultQuery("search", "")               // Get search term from query string
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))  // Default limit to 10
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0")) // Default offset to 0
	state, _ := strconv.Atoi(ctx.DefaultQuery("state", ""))
	orderby := ctx.DefaultQuery("orderby", "id")
	order := ctx.DefaultQuery("order", "asc")
	var usersTotal []db.User
	if order == "asc" {
		// Fetch users with the search term and state filter
		users, err := server.store.GetUsersAsc(ctx, db.GetUsersAscParams{
			Column1: &searchTerm,
			Limit:   int32(limit),
			Offset:  int32(offset),
			Column4: orderby,
			State:   int64(state),
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		usersTotal = users
	} else {
		// Fetch users with the search term and state filter
		users, err := server.store.GetUsersDesc(ctx, db.GetUsersDescParams{
			Column1: &searchTerm,
			Limit:   int32(limit),
			Offset:  int32(offset),
			Column4: orderby,
			State:   int64(state),
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		usersTotal = users
	}
	rsp := Response{
		Status: http.StatusOK,
		Data:   usersTotal,
		Meta: Meta{
			Total:  len(usersTotal),
			Offset: offset,
			Limit:  limit,
		},
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getUserById(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUserById(ctx, int32(idInt))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
}
func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.UpdateUserParams{
		Fullname:     &req.FullName,
		Email:        &req.Email,
		Phone:        &req.Phone,
		Password:     &req.Password,
		Avatar:       &req.Avatar,
		Role:         &req.Role,
		PasswordHash: &hashedPassword,
	}
	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
}
func (server *Server) deleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.DeleteUser(ctx, db.DeleteUserParams{
		State: int64(model.UserStateDeleted),
		ID:    int32(idInt),
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required"`
}

type DataObject struct {
	AccessToken string `json:"access_token"`
}
type loginUserResponse struct {
	Status int        `json:"status"`
	Error  string     `json:"error"`
	Errors []string   `json:"errors"`
	Data   DataObject `json:"data"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	var rsp loginUserResponse
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.Username != "admin" || req.Password != "admin" {
		rsp = loginUserResponse{
			Status: 404,
			Error:  "username and password is invalid",
			Errors: []string{"username and password is invalid"},
			Data:   DataObject{},
		}
		ctx.JSON(http.StatusUnauthorized, rsp)
		return
	}
	accessToken, _, err := server.tokenMaker.CreateToken(
		req.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp = loginUserResponse{
		Status: 200,
		Error:  "",
		Errors: nil,
		Data: DataObject{
			AccessToken: accessToken,
		},
	}
	ctx.JSON(http.StatusOK, rsp)
}
