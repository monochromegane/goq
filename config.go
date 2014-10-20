package goq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	Targets map[string]Target `json:"targets"`
}

type Target struct {
	Driver string `json:"driver"`
	Dsn    string `json:"dsn"`
	Dir    string `json:"query_dir"`
	Prefix string `json:"prefix"`
}

func loadConfig() Config {
	var config Config

	bytes, err := ioutil.ReadFile(filepath.Join(getHomeDir(), ".goq/config.json"))
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func (c Config) Find(name string) (Target, error) {
	if t, exist := c.Targets[name]; exist {
		return t, nil
	}
	return Target{}, fmt.Errorf("%s not found.", name)
}

func getHomeDir() string {
	usr, err := user.Current()
	var homeDir string
	if err == nil {
		homeDir = usr.HomeDir
	} else {
		// Maybe it's cross compilation without cgo support. (darwin, unix)
		homeDir = os.Getenv("HOME")
	}
	return homeDir
}
