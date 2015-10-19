package main

import (
	"bytes"
	"encoding/json"
	"github.com/nahanni/go-ucl"
	"io/ioutil"
	"log"
)

type Config struct {
	Timeout string     `json:"timeout"`
	Api     string     `json:"api"`
	JobDb   string     `json:"jobdb"`
	Auth    AuthConfig `json:"auth"`
}

type AuthConfig struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

func (c *Config) populateFromFile(fname string) error {
	file, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}

	log.Printf("Loading configuration from %s", fname)

	// UCL parses into map[string]interface{}
	fileBytes := bytes.NewBuffer([]byte(file))
	parser := ucl.NewParser(fileBytes)
	uclData, err := parser.Ucl()
	if err != nil {
		log.Fatal("UCL error: ", err)
	}

	// take detour via JSON to load UCL into struct
	uclJson, err := json.Marshal(uclData)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(uclJson), &c)

	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
