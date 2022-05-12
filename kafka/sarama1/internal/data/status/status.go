package status

type Status struct {
	Message string `json:"msg"`
	Code    int    `json:"code"` // HTTP status code
}
