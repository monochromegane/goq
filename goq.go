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

	return query(db, q, args)
}

type q struct {
	file   string
	prefix string
}

func (query q) name() string {
	return strings.TrimPrefix(query.file[:strings.LastIndex(query.file, ".")], query.prefix)
}

func listQuery(dir, prefix string) ([]q, error) {
	var list []q
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return list, err
	}
	for _, f := range files {
		if filepath.HasPrefix(f.Name(), prefix) {
			list = append(list, q{file: f.Name(), prefix: prefix})
		}
	}
	return list, nil
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

func query(db *sql.DB, q string, args []string) ([]string, [][]string) {
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
