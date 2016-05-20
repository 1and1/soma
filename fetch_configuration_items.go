/*
 * Copyright (c) 2016, 1&1 Internet SE
 * Written by Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved.
 */

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"time"


	"gopkg.in/resty.v0"

	"github.com/asaskevich/govalidator"
	"github.com/julienschmidt/httprouter"
)

type NotifyMessage struct {
	Uuid string `json:"uuid" valid:"uuidv4"`
	Path string `json:"path" valid:"abspath"`
}

func FetchConfigurationItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		dec    *json.Decoder
		msg    NotifyMessage
		err    error
		soma   *url.URL
		client *resty.Client
		resp   *resty.Response
		res    proto.Result
	)
	dec = json.NewDecoder(r.Body)
	if err = dec.Decode(msg); err != nil {
		dispatchBadRequest(&w, err.Error())
		return
	}
	govalidator.SetFieldsRequiredByDefault(true)
	govalidator.TagMap["abspath"] = govalidator.Validator(func(str string) bool {
		return filepath.IsAbs(str)
	})
	if ok, err := govalidator.ValidateStruct(msg); !ok {
		dispatchBadRequest(&w, err.Error())
		return
	}

	soma, _ = url.Parse(Eye.Soma.url.String())
	soma.Path = fmt.Sprintf("%s/%s", msg.Path, msg.Uuid)
	client = resty.New().SetTimeout(500 * time.Millisecond)
	if resp, err = client.R().Get(soma.String()); err != nil || resp.StatusCode() > 299 {
		if err == nil {
			err = fmt.Errorf(resp.Status())
		}
		dispatchPrecondition(&w, err.Error())
		return
	}
	if err = json.Unmarshal(resp.Body(), res); err != nil {
		dispatchUnprocessable(&w, err.Error())
		return
	}
	if res.StatusCode != 200 {
		dispatchGone(&w, err.Error())
		return
	}
	if len(*res.Deployments) != 1 {
		dispatchPrecondition(&w, err.Error())
		return
	}
	if err = CheckUpdateOrInsertOrDelete(&(*res.Deployments)[0]); err != nil {
		dispatchInternalServerError(&w, err.Error())
		return
	}
	dispatchNoContent(&w)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
