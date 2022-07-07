package util

import (
	"bookman/entity"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

const DefaultLimit = 20

func NewPaginationFromRequest(c *gin.Context) *entity.Pagination {
	offset := 0
	limit := DefaultLimit

	query := c.Request.URL.Query()
	for k, v := range query {
		queryValue := v[len(v)-1]
		switch k {
		case "offset":
			offset, _ = strconv.Atoi(queryValue)
			break
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
			break
		}
	}

	return &entity.Pagination{
		Offset: offset,
		Limit:  limit,
	}
}

func NewPaginationFromMessage(msg string) (*entity.Pagination, error) {
	var pagination entity.Pagination
	err := json.Unmarshal([]byte(msg), &pagination)
	if err != nil {
		return nil, err
	}

	if pagination.Limit == 0 {
		pagination.Limit = DefaultLimit
	}
	return &pagination, nil
}
