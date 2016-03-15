/*
Copyright (c) 2016, Jörg Pernfuß <code.jpe@gmail.com>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/asaskevich/govalidator"
	"github.com/julienschmidt/httprouter"
)

// global variables
var conn *sql.DB
var Eye EyeConfig

func main() {
	version := "0.9.5"
	/*
	 * Read configuration file
	 */
	log.Printf("Starting runtime config initialization, Eye v%s", version)
	err := Eye.readConfigFile("eye.conf")
	if err != nil {
		log.Fatal(err)
	}

	/*
	 * Construct listen address
	 */
	Eye.Daemon.url = &url.URL{}
	Eye.Daemon.url.Host = fmt.Sprintf("%s:%s", Eye.Daemon.Listen, Eye.Daemon.Port)
	if Eye.Daemon.Tls {
		Eye.Daemon.url.Scheme = "https"
		if ok, ptype := govalidator.IsFilePath(Eye.Daemon.Cert); !ok {
			log.Fatal("Missing required certificate configuration config/daemon/cert-file")
		} else {
			if ptype != govalidator.Unix {
				log.Fatal("config/daemon/cert-File: valid Windows paths are not helpful")
			}
		}
		if ok, ptype := govalidator.IsFilePath(Eye.Daemon.Key); !ok {
			log.Fatal("Missing required key configuration config/daemon/key-file")
		} else {
			if ptype != govalidator.Unix {
				log.Fatal("config/daemon/key-file: valid Windows paths are not helpful")
			}
		}
	} else {
		Eye.Daemon.url.Scheme = "http"
	}

	/*
	 * Initialize database
	 */
	connectToDatabase()
	prepareStatements()
	// Close() must be deferred here since it triggers on function exit
	defer Eye.run.check_item.Close()
	defer Eye.run.check_lookup.Close()
	defer Eye.run.delete_item.Close()
	defer Eye.run.delete_lookup.Close()
	defer Eye.run.get_config.Close()
	defer Eye.run.get_items.Close()
	defer Eye.run.get_lookup.Close()
	defer Eye.run.insert_item.Close()
	defer Eye.run.insert_lookup.Close()
	defer Eye.run.item_count.Close()
	defer Eye.run.retrieve.Close()
	defer Eye.run.update_item.Close()
	go pingDatabase()

	/*
	 * Register http handlers
	 */
	router := httprouter.New()
	router.GET("/api/v1/configuration/:lookup", RetrieveConfigurationItems)
	router.GET("/api/v1/item/", ListConfigurationItems)
	router.POST("/api/v1/item/", AddConfigurationItem)
	router.GET("/api/v1/item/:item", GetConfigurationItem)
	router.PUT("/api/v1/item/:item", UpdateConfigurationItem)
	router.DELETE("/api/v1/item/:item", DeleteConfigurationItem)
	router.POST("/api/v1/notify/", FetchConfigurationItems)

	if Eye.Daemon.Tls {
		log.Fatal(http.ListenAndServeTLS(
			Eye.Daemon.url.Host,
			Eye.Daemon.Cert,
			Eye.Daemon.Key,
			router))
	} else {
		log.Fatal(http.ListenAndServe(Eye.Daemon.url.Host, router))
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
