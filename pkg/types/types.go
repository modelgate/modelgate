package types

type PageParam struct {
	Page     int64
	PageSize int64
	OrderBy  string
}

func NewPageParam(page, pageSize int64, orderBy string) *PageParam {
	return &PageParam{
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
	}
}
