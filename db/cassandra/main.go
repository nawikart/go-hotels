package cassandra

import (
	"fmt"
	"github.com/gocql/gocql"
)

var Session *gocql.Session

func init() {
	var err error

	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "mydb"
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("cassandra init done")
}

// func Connect() *gocql.Session{
// 	var err error

// 	cluster := gocql.NewCluster("127.0.0.1")
// 	cluster.Keyspace = "demo"
// 	Session, err = cluster.CreateSession()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("cassandra init done")
// 	return Session
// }
