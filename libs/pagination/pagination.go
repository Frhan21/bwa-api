package pagination

import (
	"bwa-api/core/domain/entity"
	"math"
)

type PaginationInterface interface {
	AddPagination(totalData, page, perPage int) (*entity.Pages, error)
}

type Options struct {
	TotalData int
	Page      int
	PerPage   int
}

// AddPagination implements PaginationInterface.
func (o *Options) AddPagination(totalData int, page int, perPage int) (*entity.Pages, error) {
	newPage := page
	if newPage == 0 {
		return nil, ErrorPage
	}

	limit := 10
	if perPage > 0 {
		limit = perPage
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(limit)))
	lasPage := (newPage * limit)
	firstPage := lasPage - limit

	if totalData < lasPage {
		lasPage = totalData
	}

	zeroPage := &entity.Pages{PageCount: 1, Page: newPage}
	if totalData == 0 && newPage == 1 {
		return zeroPage, nil
	}

	if newPage > totalPage {
		return nil, ErrorPage
	}

	pages := &entity.Pages{Page: newPage, PageCount: totalPage, First: firstPage, Last: lasPage}
	return pages, nil
}

func Pagination() PaginationInterface {
	pagination := new(Options)

	return pagination
}
