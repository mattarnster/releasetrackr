package requests

// TokenRequest is the request a user makes in order
// to get a token for login.
type TokenRequest struct {
	Email string `json:"email"`
}
