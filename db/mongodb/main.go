package mongodb

import (
	"gopkg.in/mgo.v2"
)

func Connect() (*mgo.Session, error) {
	var session, err = mgo.Dial("localhost")
	if err != nil {
		return nil, err
	}
	return session, nil
}
