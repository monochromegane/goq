package goq

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

type query struct {
	dir    string
	file   string
	prefix string
}

func (q query) name() string {
	return strings.TrimPrefix(q.file[:strings.LastIndex(q.file, ".")], q.prefix)
}

func (q query) query() (string, error) {
	bytes, err := ioutil.ReadFile(filepath.Join(q.dir, q.file))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
