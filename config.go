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
	ReadOnly    bool       `json:"readonly,string" valid:"-"`
	Daemon      EyeDaemon  `json:"daemon" valid:"required"`
	Database    DbConfig   `json:"database" valid:"required"`
	Soma        SomaConfig `json:"soma" valid:"required"`
	run         EyeRuntime `json:"-" valid:"-"`
}

type DbConfig struct {
	Host    string `json:"host" valid:"dns"`
	User    string `json:"user" valid:"alphanum"`
	Name    string `json:"name" valid:"alphanum"`
	Port    string `json:"port" valid:"port"`
	Pass    string `json:"password" valid:"-"`
	Timeout string `json:"timeout" valid:"numeric"`
	TlsMode string `json:"tlsmode" valid:"alpha"`
}

type SomaConfig struct {
	url     *url.URL `json:"-"`
	Address string   `json:"address" valid:"requrl"`
}

type EyeDaemon struct {
	Listen string `json:"listen" valid:"ip"`
	Port   string `json:"port" valid:"port"`
	Tls    bool   `json:"tls,string" valid:"-"`
	Cert   string `json:"cert-file" valid:"optional"`
	Key    string `json:"key-file" valid:"optional"`
}

type EyeRuntime struct {
	conn          *sql.DB   `json:"-" valid:"-"`
	check_item    *sql.Stmt `json:"-" valid:"-"`
	update_item   *sql.Stmt `json:"-" valid:"-"`
	check_lookup  *sql.Stmt `json:"-" valid:"-"`
	insert_lookup *sql.Stmt `json:"-" valid:"-"`
	insert_item   *sql.Stmt `json:"-" valid:"-"`
	delete_item   *sql.Stmt `json:"-" valid:"-"`
	delete_lookup *sql.Stmt `json:"-" valid:"-"`
	get_lookup    *sql.Stmt `json:"-" valid:"-"`
	item_count    *sql.Stmt `json:"-" valid:"-"`
	get_config    *sql.Stmt `json:"-" valid:"-"`
	get_items     *sql.Stmt `json:"-" valid:"-"`
	retrieve      *sql.Stmt `json:"-" valid:"-"`
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
