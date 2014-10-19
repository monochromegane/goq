package goq

import (
	"database/sql"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func Query(target, query string) ([]string, [][]string) {

	config := loadConfig()
	t := config.find(target)
	db, err := sql.Open(t.Driver, t.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(t.Dir)
	if err != nil {
		log.Fatal(err)
	}
	var q string
	for _, f := range files {
		if f.Name()[:strings.LastIndex(f.Name(), ".")] == t.Prefix+query {
			bytes, err := ioutil.ReadFile(filepath.Join(t.Dir, f.Name()))
			if err != nil {
				log.Fatal(err)
			}
			q = string(bytes)
		}
	}

	stmt, err := db.Prepare(q)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}

	columns, _ := rows.Columns()
	var values [][]string
	for rows.Next() {
		vals := make([]string, len(columns))
		ptrs := make([]interface{}, len(columns))
		for i, _ := range vals {
			ptrs[i] = &vals[i]
		}
		rows.Scan(ptrs...)
		values = append(values, vals)
	}
	rows.Close()

	return columns, values

}
