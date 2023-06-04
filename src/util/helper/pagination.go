package helper

import (
	"backend/src/constant"
	"strconv"

	"github.com/pkg/errors"
)

type Pagination struct {
	PageStr  string `form:"page"`
	LimitStr string `form:"limit"`

	Page   int
	Limit  int
	Offset int
}

type PaginationResponse struct {
	Page        int `json:"page"`
	DataPerPage int `json:"data_per_page"`
	TotalData   int `json:"total_data"`
	TotalPage   int `json:"total_page"`
}

func (pg *Pagination) ValidatePagination() (err error) {
	if pg.PageStr != "" {
		pg.Page, err = strconv.Atoi(pg.PageStr)
		if err != nil {
			err = errors.Wrap(err, "validate pagination: page")
			return
		}
	}
	if pg.LimitStr != "" {
		pg.Limit, err = strconv.Atoi(pg.LimitStr)
		if err != nil {
			err = errors.Wrap(err, "validate pagination: limit")
			return
		}
	}

	if pg.Limit < constant.MinLimit || pg.Limit > constant.MaxLimit || pg.Limit == 0 {
		pg.Limit = constant.LimitDefault
	}
	if pg.Page < constant.MinPage || pg.Page == 0 {
		pg.Page = constant.PageDefault
	}
	pg.Offset = (pg.Page * pg.Limit) - pg.Limit

	return
}
