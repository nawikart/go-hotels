package mysql

import (
	_ "../../../github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/demo")
	if err != nil {
		fmt.Println(err.Error())
	} 
	// else {
	// 	fmt.Println("db is connected")
	// }
	//defer db.Close()
	// make sure connection is available
	err = db.Ping()
	// fmt.Println(err)
	if err != nil {
		fmt.Println("db is not connected")
		fmt.Println(err.Error())
	}
	return db
}