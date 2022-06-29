package util

import (
	"bookman/common"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseID(c *gin.Context, key string) (int, error) {
	id, err := strconv.Atoi(c.Param(key))
	if err != nil {
		return 0, common.ErrInvalidParam(err)
	}
	return id, nil
}

func ParseBody[T any](c *gin.Context) (*T, error) {
	var data T
	if err := c.ShouldBind(&data); err != nil {
		return nil, common.ErrInvalidRequestBody(err)
	}
	return &data, nil
}
