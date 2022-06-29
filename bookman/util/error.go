package util

import (
	"bookman/common"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseError(c *gin.Context, err error) {
	if err != nil {
		var customErr *common.Err
		if errors.As(err, &customErr) {
			c.JSON(customErr.HTTPStatus, customErr)
			return
		}
		c.JSON(http.StatusInternalServerError, common.ErrInternal(err))
		return
	}
}
