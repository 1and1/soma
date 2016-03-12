package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/julienschmidt/httprouter"
)

// global variables
var conn *sql.DB
var Eye EyeConfig

func main() {
	version := "0.0.7"
	log.Printf("Starting runtime config initialization, Eye v%s", version)
	err := Eye.readConfigFile("eye.conf")
	if err != nil {
		log.Fatal(err)
	}

	connectToDatabase()
	go pingDatabase()

	router := httprouter.New()

	//router.GET("/api/v1/configuration/:lookup", RetrieveConfigurationItems)
	router.GET("/api/v1/item/", ListConfigurationItems)
	router.POST("/api/v1/item/", AddConfigurationItem)
	//router.GET("/api/v1/item/:item", GetConfigurationItem)
	router.PUT("/api/v1/item/:item", UpdateConfigurationItem)
	router.DELETE("/api/v1/item/:item", DeleteConfigurationItem)
	router.POST("/api/v1/notify/", FetchConfigurationItems)

	j, _ := json.Marshal(Eye)
	fmt.Println(string(j))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
