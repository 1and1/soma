package main

import (
	"bytes"
	"encoding/json"
	"github.com/nahanni/go-ucl"
	"io/ioutil"
	"log"
)

type SomaConfig struct {
	Environment string       `json:"environment"`
	Timeout     string       `json:"timeout"`
	TlsMode     string       `json:"tlsmode"`
	ReadOnly    bool         `json:"readonly,string"`
	Database    SomaDbConfig `json:"database"`
}

type SomaDbConfig struct {
	Host string `json:"host"`
	User string `json:"user"`
	Name string `json:"name"`
	Port string `json:"port"`
	Pass string `json:"password"`
}

func (c *SomaConfig) readConfigFile(fname string) error {
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
