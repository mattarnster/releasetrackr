package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mattarnster/releasetrackr/responses"
)

// IndexHandler lives at / and shows information about the application
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	json, _ := json.Marshal(&responses.IndexResponse{
		Name: "releasetrackr",
		Ver:  "0.1",
	})
	w.Write(json)
}
