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

func Query(target, queryName string) ([]string, [][]string) {

	config := loadConfig()
	t := config.find(target)

	queryFile := queryFile{dir: t.Dir, prefix: t.Prefix}

	db, err := sql.Open(t.Driver, t.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	q, err := queryFile.find(queryName)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := query(db, q)
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

type queryFile struct {
	dir    string
	prefix string
}

func (q queryFile) find(name string) (string, error) {
	files, err := ioutil.ReadDir(q.dir)
	if err != nil {
		return "", err
	}
	for _, f := range files {
		if f.Name()[:strings.LastIndex(f.Name(), ".")] == q.prefix+name {
			bytes, err := ioutil.ReadFile(filepath.Join(q.dir, f.Name()))
			if err != nil {
				return "", err
			}
			return string(bytes), nil
		}
	}
	return "", fmt.Errorf("%d not found.", name)
}

func query(db *sql.DB, query string) (*sql.Rows, error) {
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
