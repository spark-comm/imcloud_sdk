package sdk_params_callback

type PageResult[T any] struct {
	Total int64 `json:"total"`
	List  []T   `json:"list"`
}

func NewPageResult[T any](list []T, total int64) *PageResult[T] {
	return &PageResult[T]{
		Total: total,
		List:  list,
	}
}
