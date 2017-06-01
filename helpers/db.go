package helpers

import mgo "gopkg.in/mgo.v2"

// DbConnection Will hold the global DB Connection
var DbConnection *mgo.Session

// GetDbSession returns the currently active DB connection (if not, then it creates one)
func GetDbSession() (*mgo.Session, error) {
	if DbConnection == nil {
		session, err := mgo.Dial("localhost:32768")
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
