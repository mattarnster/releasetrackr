package requests

// TrackRequest is used by the TrackHandler - this is what the user
// POSTs to the server to create a new "Track"
type TrackRequest struct {
	Repo              string `json:"repo"`
	Email             string `json:"email"`
	RecaptchaResponse string `json:"recaptcha_response"`
}
