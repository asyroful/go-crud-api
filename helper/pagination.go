package helper

import (
	"strconv"
	"go-crud-api/models"
)

func SetPaginationFromQuery(queryLimit string, queryPage string) models.QueryPagination {

	limit := -1
	page := -1
	offset := -1

	if queryLimit != "" {
		limit, _ = strconv.Atoi(queryLimit)
	}

	if queryPage != "" {
		page, _ = strconv.Atoi(queryPage)
	}

	if limit != -1 || page != -1 {
		offset = (page - 1) * limit
	}

	pagination := models.QueryPagination{
		Limit:  limit,
		Offset: offset,
	}

	return pagination
}
