package responses

// ErrorResponse is used by the Handlers and gives a
// generic code and error back to the user
type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}
