package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/nahanni/go-ucl"
)

type SomaConfig struct {
	Environment string         `json:"environment"`
	ReadOnly    bool           `json:"readonly,string"`
	Database    SomaDbConfig   `json:"database"`
	Daemon      SomaDaemon     `json:"daemon"`
	Auth        SomaAuthConfig `json:"authentication"`
}

type SomaDbConfig struct {
	Host    string `json:"host"`
	User    string `json:"user"`
	Name    string `json:"name"`
	Port    string `json:"port"`
	Pass    string `json:"password"`
	Timeout string `json:"timeout"`
	TlsMode string `json:"tlsmode"`
}

type SomaDaemon struct {
	url    *url.URL `json:"-"`
	Listen string   `json:"listen"`
	Port   string   `json:"port"`
	Tls    bool     `json:"tls,string"`
	Cert   string   `json:"cert-file"`
	Key    string   `json:"key-file"`
}

type SomaAuthConfig struct {
	KexExpirySeconds   uint64 `json:"kex_expiry,string"`
	TokenExpirySeconds uint64 `json:"token_expiry,string"`
	// dd if=/dev/random bs=1M count=1 2>/dev/null | sha512
	TokenSeed string `json:"token_seed"`
	TokenKey  string `json:"token_key"`
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
