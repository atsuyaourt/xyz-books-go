package util

type PaginatedList[T any] struct {
	CurrentPage int32 `json:"current_page"`
	PerPage     int32 `json:"per_page"`
	TotalPages  int32 `json:"total_pages"`
	NextPage    int32 `json:"next_page"`
	PrevPage    int32 `json:"prev_page"`
	TotalItems  int32 `json:"total_items"`
	Items       []T   `json:"items"`
}

func NewPaginatedList[T any](curPage, limit, totalItems int32, items []T) PaginatedList[T] {
	p := PaginatedList[T]{
		CurrentPage: curPage,
		PerPage:     limit,
		TotalItems:  totalItems,
		Items:       items,
	}

	p.TotalPages = p.calculateTotalPages()
	p.setNextPage()
	p.setPrevPage()

	return p
}

func (p PaginatedList[T]) calculateTotalPages() int32 {
	if p.PerPage <= 0 {
		return 0
	}
	return (p.TotalItems + p.PerPage - 1) / p.PerPage
}

func (p *PaginatedList[T]) setNextPage() {
	nextPage := p.CurrentPage + 1
	if nextPage <= p.TotalPages {
		p.NextPage = nextPage
	}
}

func (p *PaginatedList[T]) setPrevPage() {
	prevPage := p.CurrentPage - 1
	if prevPage >= 1 {
		p.PrevPage = prevPage
	}
}
