package responses

// SuccessResponse is a generic success response which the user will receive
type SuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
