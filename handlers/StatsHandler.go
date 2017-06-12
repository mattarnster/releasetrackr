package handlers

import "net/http"

// StatsHandler deals with showing stats of how many people
// have subscribed for notifications and stuff.
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
}
