package helper

import (
	"go-crud-api/models"
	"strconv"
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
		Page:   page,
	}

	return pagination
}
