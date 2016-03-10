package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gopkg.in/resty.v0"

	"github.com/julienschmidt/httprouter"
)

// global variables
var conn *sql.DB
var MonsoonCfg MonsoonConfig

type notifyMessage struct {
	uuid string `json:"uuid"`
}

func main() {
	version := "0.0.1"
	log.Printf("Starting runtime config initialization, MonsoonCfg v%s", version)
	err := MonsoonCfg.readConfigFile("monsoon.conf")
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()

	router.POST("/deployment/", AddConfiguration)
}

func AddConfiguration(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	msg := notifyMessage{}
	if err := decoder.Decode(msg); err != nil {
		log.Fatal(err)
	}
	api := "http://127.0.0.1"

	client := resty.New().SetTimeout(500 * time.Millisecond)
	resp, err := client.R().Get(fmt.Sprintf("%s/deployments/id/%s", api, msg.uuid))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
