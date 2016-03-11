package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/asaskevich/govalidator"
	"github.com/nahanni/go-ucl"
)

type EyeConfig struct {
	Environment string     `json:"environment" valid:"alpha"`
	Timeout     string     `json:"timeout" valid:"numeric"`
	TlsMode     string     `json:"tlsmode" valid:"alpha"`
	ReadOnly    bool       `json:"readonly,string" valid:"-"`
	Database    DbConfig   `json:"database" valid:"required"`
	Soma        SomaConfig `json:"soma" valid:"required"`
	conn        *sql.DB    `json:"-" valid:"-"`
}

type DbConfig struct {
	Host string `json:"host" valid:"dns"`
	User string `json:"user" valid:"alphanum"`
	Name string `json:"name" valid:"alphanum"`
	Port string `json:"port" valid:"port"`
	Pass string `json:"password" valid:"-"`
}

type SomaConfig struct {
	url     *url.URL `json:"-"`
	Address string   `json:"address" valid:"requrl"`
}

func (c *EyeConfig) readConfigFile(fname string) error {
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

	govalidator.SetFieldsRequiredByDefault(true)
	if ok, err := govalidator.ValidateStruct(c); !ok {
		return err
	}
	c.Soma.url, _ = url.Parse(c.Soma.Address)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
