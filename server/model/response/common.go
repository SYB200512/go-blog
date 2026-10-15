package response

type PageResult struct {
	List  interface{}
	Total int64 `json:"total"`
}
