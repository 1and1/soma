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
	version := "0.9.2"
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
	listen := url.URL{}
	listen.Host = fmt.Sprintf("%s:%s", Eye.Daemon.Listen, Eye.Daemon.Port)
	if Eye.Daemon.Tls {
		listen.Scheme = "https"
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
		listen.Scheme = "http"
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
			listen.String(),
			Eye.Daemon.Cert,
			Eye.Daemon.Key,
			router))
	} else {
		log.Fatal(http.ListenAndServe(listen.String(), router))
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
