package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"


	"gopkg.in/resty.v0"

	"github.com/asaskevich/govalidator"
	"github.com/julienschmidt/httprouter"
)

func FetchConfigurationItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		dec    *json.Decoder
		msg    notifyMessage
		err    error
		soma   *url.URL
		client *resty.Client
		resp   *resty.Response
		res    somaproto.DeploymentDetailsResult
	)
	dec = json.NewDecoder(r.Body)
	if err = dec.Decode(msg); err != nil {
		os.Exit(1)
	}
	govalidator.SetFieldsRequiredByDefault(true)
	govalidator.TagMap["abspath"] = govalidator.Validator(func(str string) bool {
		return filepath.IsAbs(str)
	})
	if ok, err := govalidator.ValidateStruct(msg); !ok {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	soma, _ = url.Parse(Eye.Soma.url.String())
	soma.Path = fmt.Sprintf("%s/%s", msg.Path, msg.Uuid)
	client = resty.New().SetTimeout(500 * time.Millisecond)
	if resp, err = client.R().Get(soma.String()); err != nil || resp.StatusCode() > 299 {
		if err == nil {
			err = fmt.Errorf(resp.Status())
		}
		log.Println(err)
		http.Error(w, err.Error(), http.StatusPreconditionFailed)
		return
	}
	if err = json.Unmarshal(resp.Body(), res); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 422)
		return
	}
	if res.Code != 200 {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusGone)
		return
	}
	if len(res.Deployments) != 1 {
		http.Error(w, err.Error(), http.StatusPreconditionFailed)
		return
	}
	if err = CheckUpdateOrInsert(&res.Deployments[0]); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
