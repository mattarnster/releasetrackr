package handlers

import "net/http"

// ManageHandler is for management of subscribed repos.
func ManageHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
}
