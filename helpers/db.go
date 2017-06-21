package helpers

import mgo "gopkg.in/mgo.v2"
import "os"

// DbConnection Will hold the global DB Connection
var DbConnection *mgo.Session

// GetDbSession returns the currently active DB connection (if not, then it creates one)
func GetDbSession() (*mgo.Session, error) {
	if DbConnection == nil {
		session, err := mgo.Dial(connectionString())
		if err != nil {
			return nil, err
		}
		DbConnection = session
	}
	return DbConnection, nil
}

// KillDbSession kills the currently active DB session
func KillDbSession() {
	if DbConnection != nil {
		DbConnection.Close()
	}
}

// Build the connection string
func connectionString() string {
	buildHost := os.Getenv("MONGO_HOST")
	buildPort := "27017"
	if os.Getenv("MONGO_PORT") != "" {
		buildPort = os.Getenv("MONGO_PORT")
	}

	return buildHost + ":" + string(buildPort)
}
