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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/1and1/soma/lib/proto"
	"github.com/asaskevich/govalidator"
	"github.com/satori/go.uuid"
)

func CalculateLookupId(id uint64, metric string) string {
	asset := strconv.FormatUint(id, 10)
	hash := sha256.New()
	hash.Write([]byte(asset))
	hash.Write([]byte(metric))

	return hex.EncodeToString(hash.Sum(nil))
}

func Itemize(details *proto.Deployment) (string, *ConfigurationItem, error) {
	var (
		fqdn, dns_zone string
		err            error
	)
	lookupID := CalculateLookupId(details.Node.AssetId, details.Metric.Path)

	item := &ConfigurationItem{
		Metric:   details.Metric.Path,
		Interval: details.CheckConfig.Interval,
		HostId:   strconv.FormatUint(details.Node.AssetId, 10),
		Metadata: ConfigurationMetaData{
			Monitoring: details.Monitoring.Name,
			Team:       details.Team.Name,
		},
		Thresholds: []ConfigurationThreshold{},
	}
	if item.ConfigurationItemId, err = uuid.FromString(details.CheckInstance.InstanceId); err != nil {
		return "", nil, err
	}

	switch item.Metric {
	case `disk.write.per.second`:
		fallthrough
	case `disk.read.per.second`:
		fallthrough
	case `disk.free`:
		fallthrough
	case `disk.usage.percent`:
		mpt := GetServiceAttributeValue(details, `filesystem`)
		if mpt == `` {
			return ``, nil, fmt.Errorf(`Disk metric is missing filesystem service attribute`)
		}
		item.Metric = fmt.Sprintf("%s:%s", item.Metric, mpt)
		// recalculate lookupID
		lookupID = CalculateLookupId(details.Node.AssetId, item.Metric)
	}

	// set oncall duty if available
	if details.Oncall != nil && details.Oncall.Id != "" {
		item.Oncall = fmt.Sprintf("%s (%s)", details.Oncall.Name, details.Oncall.Number)
	}

	// construct item.Metadata.Targethost with help of system properties
	if details.Properties != nil {
		for _, prop := range *details.Properties {
			switch prop.Name {
			case "fqdn":
				fqdn = prop.Value
			case "dns_zone":
				dns_zone = prop.Value
			}
		}
	}
	switch {
	case len(fqdn) > 0:
		item.Metadata.Targethost = fqdn
	case len(dns_zone) > 0:
		item.Metadata.Targethost = fmt.Sprintf("%s.%s", details.Node.Name, dns_zone)
	default:
		item.Metadata.Targethost = details.Node.Name
	}

	// construct item.Metadata.Source
	if details.Service != nil && details.Service != `` {
		item.Metadata.Source = fmt.Sprintf("%s, %s", details.Service.Name, details.CheckConfig.Name)
	} else {
		item.Metadata.Source = fmt.Sprintf("System (%s), %s", details.Node.Name, details.CheckConfig.Name)
	}

	// slurp all thresholds
	for _, thr := range details.CheckConfig.Thresholds {
		t := ConfigurationThreshold{
			Predicate: thr.Predicate.Symbol,
			Level:     thr.Level.Numeric,
			Value:     thr.Value,
		}
		item.Thresholds = append(item.Thresholds, t)
	}

	govalidator.SetFieldsRequiredByDefault(true)
	if ok, err := govalidator.ValidateStruct(item); !ok {
		log.Println(err)
		return "", nil, err
	}
	return lookupID, item, nil
}

func GetServiceAttributeValue(details *proto.Deployment, attribute string) string {
	if details.Service == nil {
		return ``
	}
	if len(details.Service.Attributes) == 0 {
		return ``
	}
	for _, attr := range details.Service.Attributes {
		if attr.Name == attribute {
			return attr.Value
		}
	}
	return ``
}

func abortOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// 200
func dispatchJsonOK(w *http.ResponseWriter, jsonb *[]byte) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(*jsonb)
}

// 204
func dispatchNoContent(w *http.ResponseWriter) {
	(*w).WriteHeader(http.StatusNoContent)
	(*w).Write(nil)
}

// 400
func dispatchBadRequest(w *http.ResponseWriter, err string) {
	http.Error(*w, err, http.StatusBadRequest)
	log.Println(err)
}

// 404
func dispatchNotFound(w *http.ResponseWriter) {
	http.Error(*w, "No items found", http.StatusNotFound)
	log.Println("No items found")
}

// 410
func dispatchGone(w *http.ResponseWriter, err string) {
	http.Error(*w, err, http.StatusGone)
	log.Println(err)
}

// 412
func dispatchPrecondition(w *http.ResponseWriter, err string) {
	http.Error(*w, err, http.StatusPreconditionFailed)
	log.Println(err)
}

// 422
func dispatchUnprocessable(w *http.ResponseWriter, err string) {
	http.Error(*w, err, 422)
	log.Println(err)
}

// 500
func dispatchInternalServerError(w *http.ResponseWriter, err string) {
	http.Error(*w, err, http.StatusInternalServerError)
	if Eye.Volatile {
		log.Fatal(err)
	}
	log.Println(err)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
