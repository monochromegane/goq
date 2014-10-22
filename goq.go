package goq

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

func List(targetName string) {
	config := loadConfig()

	for name, target := range config.Targets {
		if name != targetName {
			continue
		}
		list, err := listQuery(target.Dir, target.Prefix)
		if err != nil {
			log.Fatal(err)
		}
		for _, q := range list {
			fmt.Printf("%s\n", q.name())
		}
	}
}

func Query(targetName, queryName string, args ...string) ([]string, [][]string) {

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

	return doQuery(db, q, args)
}

func listQuery(dir, prefix string) ([]query, error) {
	var list []query
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return list, err
	}
	for _, f := range files {
		if filepath.HasPrefix(f.Name(), prefix) {
			list = append(list, query{dir: dir, file: f.Name(), prefix: prefix})
		}
	}
	return list, nil
}

func findQuery(dir, prefix, name string) (string, error) {
	queries, err := listQuery(dir, prefix)
	if err != nil {
		return "", err
	}
	for _, q := range queries {
		if q.name() == name {
			return q.query()
		}
	}
	return "", fmt.Errorf("%s not found.", name)
}

func doQuery(db *sql.DB, q string, args []string) ([]string, [][]string) {
	rows, err := getRows(db, q, args)
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

func getRows(db *sql.DB, query string, args []string) (*sql.Rows, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	vals := make([]interface{}, len(args))
	for i, v := range args {
		vals[i] = v
	}
	rows, err := stmt.Query(vals...)
	if err != nil {
		return nil, err
	}
	return rows, err
}
