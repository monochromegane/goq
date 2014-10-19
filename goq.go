package goq

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Query(target, query string) {

	config := loadConfig()
	t := config.find(target)
	db, err := sql.Open(t.Driver, t.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var num, name string
		rows.Scan(&num, &name)
		fmt.Printf("%s:%s\n", num, name)
	}

}
