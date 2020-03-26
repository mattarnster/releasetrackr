package responses

import "releasetrackr/models"

// APIUserResponse returns User struct
type APIUserResponse struct {
	Code int         `json:"code"`
	User models.User `json:"user"`
}
