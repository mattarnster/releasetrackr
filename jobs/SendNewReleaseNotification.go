package jobs

import (
	"log"

	"github.com/mattarnster/releasetrackr/helpers"
)

// SendNewReleaseNotification sends a new release notification
func SendNewReleaseNotification() {
	sess, err := helpers.GetDbSession()

	if err != nil {
		log.Panicf("Couldn't get DB session: %v", err.Error())
		return
	}

	sess.DB("releasetrackr").C("releases")
}
