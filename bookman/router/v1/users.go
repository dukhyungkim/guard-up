package v1

import (
	"bookman/common"
	"bookman/entity"
	"bookman/service"
	"bookman/util"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const keyUserID = "userId"

type UsersRouter struct {
	userService service.UserService
}

func NewUserRouter(userService service.UserService) *UsersRouter {
	return &UsersRouter{
		userService: userService,
	}
}

func (r *UsersRouter) SetupRouter(router *gin.Engine) {
	usersGroup := router.Group("/v1/users")
	usersGroup.POST("", r.createUser)
	usersGroup.GET("", r.listUsers)
	usersGroup.GET(":"+keyUserID, r.getUser)
	usersGroup.PUT(":"+keyUserID, r.updateUser)
	usersGroup.DELETE(":"+keyUserID, r.deleteUser)
}

func (r *UsersRouter) createUser(c *gin.Context) {
	book, err := util.ParseBody[entity.User](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	newUser, err := r.userService.SaveNewUser(book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrInternal(err))
		return
	}

	result := entity.Response[*entity.User]{Data: newUser}
	c.JSON(http.StatusCreated, result)
}

func (r *UsersRouter) listUsers(c *gin.Context) {
	pagination := util.NewPaginationFromRequest(c)

	users, err := r.userService.ListUsers(pagination)
	if err != nil {
		return
	}

	result := entity.PaginatedResponse[*entity.User]{
		Pagination: pagination,
		Data:       users,
	}
	c.JSON(http.StatusOK, result)
}

func (r *UsersRouter) getUser(c *gin.Context) {
	userID, err := util.ParseID[int](c, keyUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user, err := r.userService.GetUser(&entity.User{ID: userID})
	if err != nil {
		var customErr *common.Err
		if errors.As(err, &customErr) {
			c.JSON(http.StatusNotFound, customErr)
			return
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

func (r *UsersRouter) updateUser(c *gin.Context) {
	userID, err := util.ParseID[int](c, keyUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user, err := util.ParseBody[entity.User](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user.ID = userID
	updateUser, err := r.userService.UpdateUser(user)
	if err != nil {
		var customErr *common.Err
		if errors.As(err, &customErr) {
			c.JSON(http.StatusNotFound, customErr)
			return
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, updateUser)
}

func (r *UsersRouter) deleteUser(c *gin.Context) {
	userID, err := util.ParseID[int](c, keyUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user := &entity.User{ID: userID}
	err = r.userService.DeleteUser(user)
	if err != nil {
		var customErr *common.Err
		if errors.As(err, &customErr) {
			c.JSON(http.StatusNotFound, customErr)
			return
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
