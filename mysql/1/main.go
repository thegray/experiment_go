package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	fmt.Println("Using gomysql")

	// Open up our database connection.
	// The database is called testDb
	db, err := sql.Open("mysql", "root:rootroot@tcp(127.0.0.1:3306)/testdb")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// selecting username from the users table
	results, err := db.Query("Select username FROM users")
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		var user User
		err = results.Scan(&user.Username)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(user.Username)
	}

}
