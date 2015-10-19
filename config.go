package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nahanni/go-ucl"
	"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Timeout string     `json:"timeout"`
	Api     string     `json:"api"`
	JobDb   string     `json:"jobdb"`
	LogDir  string     `json:"logdir"`
	Auth    AuthConfig `json:"auth"`
	Run     RunTimeConfig
}

type AuthConfig struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type RunTimeConfig struct {
	LevelDB     *leveldb.DB
	PathLevelDB string
	PathLogs    string
	Logger      *log.Logger
}

func (c *Config) populateFromFile(fname string) error {
	file, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}

	// UCL parses into map[string]interface{}
	fileBytes := bytes.NewBuffer([]byte(file))
	parser := ucl.NewParser(fileBytes)
	uclData, err := parser.Ucl()
	if err != nil {
		fmt.Fprintf(os.Stderr, "UCL error: %s\n", err.Error())
		os.Exit(1)
	}

	// take detour via JSON to load UCL into struct
	uclJson, err := json.Marshal(uclData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON marshal error: %s\n", err.Error())
		os.Exit(1)
	}
	json.Unmarshal([]byte(uclJson), &c)

	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
