package responses

// TokenResponse is filled with a success code and token
// if the user was able to login successfully.
type TokenResponse struct {
	Code  int    `json:"code"`
	Token string `json:"token"`
}
