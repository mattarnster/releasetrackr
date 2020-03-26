package requests

// LoginRequest is used by the APIHandler - this is what the user
// POSTs to the server to login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
