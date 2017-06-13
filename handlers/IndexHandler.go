package handlers

import (
	"html/template"
	"net/http"

	"github.com/mattarnster/releasetrackr/responses"
)

// IndexHandler lives at / and shows information about the application
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	resp := &responses.IndexResponse{
		Name: "releasetrackr",
		Ver:  "1.0",
	}

	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, resp)
}
