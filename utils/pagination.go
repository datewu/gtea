package utils

// Pagination represents a pagination object.
type Pagination struct {
	page    int
	perPage int
}

// NewPagination returns a new Pagination object.
func NewPagination(page, perPage int) *Pagination {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	return &Pagination{
		page:    page,
		perPage: perPage,
	}
}

func (p *Pagination) SetPage(page int) {
	if page < 1 {
		p.page = 1
		return
	}
	p.page = page
}

func (p *Pagination) Limit() int {
	return p.perPage
}

func (p *Pagination) Offset() int {
	return (p.page - 1) * p.perPage
}
