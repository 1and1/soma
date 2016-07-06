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
	Environment   string         `json:"environment"`
	ReadOnly      bool           `json:"readonly,string"`
	OpenInstance  bool           `json:"open.door.policy,string"`
	LifeCycleTick uint64         `json:"lifecycle.tick.seconds,string"`
	PokePath      string         `json:"notify.path.element"`
	Database      SomaDbConfig   `json:"database"`
	Daemon        SomaDaemon     `json:"daemon"`
	Auth          SomaAuthConfig `json:"authentication"`
	Ldap          SomaLdapConfig `json:"ldap"`
}

type SomaDbConfig struct {
	Host    string `json:"host"`
	User    string `json:"user"`
	Name    string `json:"database"`
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
	Cert   string   `json:"cert.file"`
	Key    string   `json:"key.file"`
}

type SomaAuthConfig struct {
	KexExpirySeconds     uint64 `json:"kex.expiry,string"`
	TokenExpirySeconds   uint64 `json:"token.expiry,string"`
	CredentialExpiryDays uint64 `json:"credential.expiry,string"`
	Activation           string `json:"activation.mode"`
	// dd if=/dev/random bs=1M count=1 2>/dev/null | sha512
	TokenSeed string `json:"token.seed"`
	TokenKey  string `json:"token.key"`
}

type SomaLdapConfig struct {
	Attribute  string `json:"uid.attribute"`
	BaseDN     string `json:"base.dn"`
	UserDN     string `json:"user.dn"`
	Address    string `json:"address"`
	Port       uint64 `json:"port,string"`
	Tls        bool   `json:"tls,string"`
	Cert       string `json:"cert.file"`
	SkipVerify bool   `json:"insecure,string"`
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

	if c.Auth.Activation == `ldap` && !c.Ldap.Tls {
		log.Println(`Account activation via LDAP configured, but LDAP/TLS disabled!`)
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
