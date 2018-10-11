package pagination

import (
	"github.com/jasonsoft/napnap"
)

type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalCount int `json:"total_count"`
	TotalPage  int `json:"total_page"`
}

func (p *Pagination) SetTotalCountAndPage(total int) {
	p.TotalCount = total
	p.TotalPage = p.TotalCount / p.PerPage
	if p.TotalCount > 0 && p.TotalPage == 0 {
		p.TotalPage = 1
	} else {
		mod := p.TotalCount % p.PerPage
		if mod > 0 {
			p.TotalPage++
		}
	}
}

func (p *Pagination) Skip() int {
	return (p.Page - 1) * p.PerPage
}

type ApiPagiationResult struct {
	Pagination Pagination  `json:"meta"`
	Data       interface{} `json:"data"`
}

type ApiCollectionResult struct {
	TotalCount int         `json:"total_count"`
	Result     interface{} `json:"result"`
}

func FromContext(c *napnap.Context) Pagination {
	page, err := c.QueryInt("page")
	if err != nil || page <= 0 {
		page = 1
	}
	perPage, err := c.QueryInt("per_page")
	if err != nil || perPage <= 0 {
		perPage = 25
	}

	pagination := Pagination{
		Page:    page,
		PerPage: perPage,
	}
	return pagination
}
