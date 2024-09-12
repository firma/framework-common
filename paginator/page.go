package paginator

import (
	"math"
)

const PageDefaultSize = 10
const PageDefaultNum = 1
const PageMaxSize = 1000

type Page struct {
	Mark      int64 `json:"mark" form:"mark" uri:"mark"`
	Page      int64 `json:"page" form:"page"  uri:"page" `                    //页码
	PageSize  int64 `json:"page_size"  uri:"page_size"   form:"page_size"`    //分页行数
	PageTotal int64 `json:"page_total"   uri:"page_total"  form:"page_total"` //分页行总数
	TotalNum  int64 `json:"total_num" uri:"total_num"  form:"total_num"`      //查询总条数
}

func (p *Page) Offset() int64 {
	if p.PageSize == 0 {
		p.PageSize = PageDefaultSize
	}
	if p.Page < PageDefaultNum {
		p.Page = PageDefaultNum
	}

	return p.PageSize * (p.Page - 1)
}

func (p *Page) Limit() int64 { //限制
	if p.PageSize == 0 {
		p.PageSize = PageDefaultSize
	} else if p.PageSize > PageMaxSize {
		p.PageSize = PageMaxSize
	}

	return p.PageSize
}

// AllowMaxPageSize 允许最大页面大小
func (p *Page) AllowMaxPageSize(max int64) Page {
	if p.PageSize > max || p.PageSize == 0 {
		p.PageSize = max
	}

	return *p
}

// Calculate 计算分页行数
func (p *Page) Calculate(total int64) int64 {
	p.TotalNum = total
	p.PageTotal = int64(math.Ceil(float64(total) / float64(p.Limit())))
	return p.PageTotal
}

// Get 获取分页行数
func (p *Page) Get(total int64) int64 {
	return int64(math.Ceil(float64(total) / float64(p.Limit())))
}
