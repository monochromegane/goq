package goq

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func Query(targetName, queryName string) ([]string, [][]string) {

	config := loadConfig()
	target, err := config.Find(targetName)
	if err != nil {
		log.Fatal(err)
	}

	q, err := findQuery(target.Dir, target.Prefix, queryName)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open(target.Driver, target.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return query(db, q)
}

func findQuery(dir, prefix, name string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, f := range files {
		if f.Name()[:strings.LastIndex(f.Name(), ".")] == prefix+name {
			bytes, err := ioutil.ReadFile(filepath.Join(dir, f.Name()))
			if err != nil {
				return "", err
			}
			return string(bytes), nil
		}
	}
	return "", fmt.Errorf("%s not found.", name)
}

func query(db *sql.DB, q string) ([]string, [][]string) {
	rows, err := getRows(db, q)
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

func getRows(db *sql.DB, query string) (*sql.Rows, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	return rows, err
}
