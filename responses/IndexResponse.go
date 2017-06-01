package responses

// IndexResponse is used by IndexHandler to show information about the application
type IndexResponse struct {
	Name string `json:"name"`
	Ver  string `json:"version"`
}
