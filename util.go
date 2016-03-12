package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"

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

func Itemize(details *somaproto.DeploymentDetails) (string, *ConfigurationItem, error) {
	var (
		fqdn, dns_zone string
		err            error
	)
	lookupID := CalculateLookupId(details.Node.AssetId, details.Metric.Metric)

	item := &ConfigurationItem{
		Metric:   details.Metric.Metric,
		Interval: details.CheckConfiguration.Interval,
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

	// set oncall duty if available
	if details.Oncall.Id != "" {
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
	if details.Service != nil {
		item.Metadata.Source = details.Service.Name
	} else {
		item.Metadata.Source = fmt.Sprintf("System (%s)", details.Node.Name)
	}

	// slurp all thresholds
	for _, thr := range details.CheckConfiguration.Thresholds {
		t := ConfigurationThreshold{
			Predicate: thr.Predicate.Predicate,
			Level:     thr.Level.Numeric,
			Value:     thr.Value,
		}
		item.Thresholds = append(item.Thresholds, t)
	}

	govalidator.SetFieldsRequiredByDefault(true)
	if ok, err := govalidator.ValidateStruct(err); !ok {
		log.Println(err)
		return "", nil, err
	}
	return lookupID, item, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
