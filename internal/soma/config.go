/*-
 * Copyright (c) 2015-2017, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"path/filepath"

	"golang.org/x/sys/unix"

	log "github.com/Sirupsen/logrus"
	"github.com/nahanni/go-ucl"
)

type Config struct {
	Environment   string     `json:"environment"`
	ReadOnly      bool       `json:"readonly,string"`
	OpenInstance  bool       `json:"open.door.policy,string"`
	LifeCycleTick uint64     `json:"lifecycle.tick.seconds,string"`
	PokePath      string     `json:"notify.path.element"`
	PokeBatchSize uint64     `json:"notify.batch.size,string"`
	PokeTimeout   uint64     `json:"notify.timeout.ms,string"`
	Observer      bool       `json:"observer,string"`
	ObserverRepo  string     `json:"-"`
	NoPoke        bool       `json:"no.poke,string"`
	PrintChannels bool       `json:"startup.print.channel.errors,string"`
	ShutdownDelay uint64     `json:"shutdown.delay.seconds,string"`
	InstanceName  string     `json:"instance.name"`
	LogPath       string     `json:"log.path"`
	Database      DbConfig   `json:"database"`
	Daemon        Daemon     `json:"daemon"`
	Auth          AuthConfig `json:"authentication"`
	Ldap          LdapConfig `json:"ldap"`
}

type DbConfig struct {
	Host    string `json:"host"`
	User    string `json:"user"`
	Name    string `json:"database"`
	Port    string `json:"port"`
	Pass    string `json:"password"`
	Timeout string `json:"timeout"`
	TlsMode string `json:"tlsmode"`
}

type Daemon struct {
	Url    *url.URL `json:"-"`
	Listen string   `json:"listen"`
	Port   string   `json:"port"`
	Tls    bool     `json:"tls,string"`
	Cert   string   `json:"cert.file"`
	Key    string   `json:"key.file"`
}

type AuthConfig struct {
	KexExpirySeconds     uint64 `json:"kex.expiry,string"`
	TokenExpirySeconds   uint64 `json:"token.expiry,string"`
	CredentialExpiryDays uint64 `json:"credential.expiry,string"`
	Activation           string `json:"activation.mode"`
	// dd if=/dev/random bs=1M count=1 2>/dev/null | sha512
	TokenSeed string `json:"token.seed"`
	TokenKey  string `json:"token.key"`
}

type LdapConfig struct {
	Attribute  string `json:"uid.attribute"`
	BaseDN     string `json:"base.dn"`
	UserDN     string `json:"user.dn"`
	Address    string `json:"address"`
	Port       uint64 `json:"port,string"`
	Tls        bool   `json:"tls,string"`
	Cert       string `json:"cert.file"`
	SkipVerify bool   `json:"insecure,string"`
}

func (c *Config) ReadConfigFile(fname string) error {
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

	if c.InstanceName == `` {
		log.Println(`Setting default value for instance.name: huxley`)
		c.InstanceName = `huxley`
	}

	if c.LogPath == `` {
		c.LogPath = filepath.Join(`/srv/soma`, c.InstanceName, `log`)
		log.Printf("Setting default value for log.path: %s\n",
			c.LogPath)
	}
	for _, p := range []string{
		c.LogPath,
		filepath.Join(c.LogPath, `job`),
		filepath.Join(c.LogPath, `repository`),
	} {
		if err := c.verifyPathWritable(p); err != nil {
			log.Fatal(`Log directory missing or not writable:`,
				p, `Error:`, err)
		}
	}

	if c.Environment == `` {
		log.Println(`Setting default value for environment: production`)
		c.Environment = `production`
	}

	if c.LifeCycleTick == 0 {
		log.Println(`Setting default value for lifecycle.tick.seconds: 60`)
		c.LifeCycleTick = 60
	}

	if c.PokeBatchSize != 0 {
		log.Println(`Configuration value notify.batch.size is deprecated and no longer has any effect.`)
	}

	if c.PokeTimeout == 0 {
		log.Println(`Setting default value for notify.timeout.ms: 1000`)
		c.PokeTimeout = 1000
	}

	if c.PokePath == `` {
		log.Println(`Setting default value for notify.path.element: /deployments/id`)
		c.PokePath = `/deployments/id`
	}

	if c.Auth.Activation == `ldap` && !c.Ldap.Tls {
		log.Println(`Account activation via LDAP configured, but LDAP/TLS disabled!`)
	}
	if c.ShutdownDelay == 0 {
		log.Println(`Setting default value for shutdown.delay.seconds: 5`)
		c.ShutdownDelay = 5
	}

	return nil
}

func (c *Config) verifyPathWritable(path string) error {
	return unix.Access(path, unix.W_OK)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
