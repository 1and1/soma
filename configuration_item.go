package main

import "github.com/satori/go.uuid"

type ConfigurationData struct {
	Configurations []ConfigurationItem `json:"configurations"`
}

type ConfigurationList struct {
	ConfigurationItemIdList []string `json:"configuration_item_id_list"`
}

type ConfigurationItem struct {
	ConfigurationItemId uuid.UUID                `json:"configuration_item_id" valid:"-"`
	Metric              string                   `json:"metric" valid:"printableascii"`
	HostId              string                   `json:"host_id" valid:"numeric"`
	Tags                []string                 `json:"tags,omitempty" valid:"-"`
	Oncall              string                   `json:"oncall" valid:"required"`
	Interval            uint64                   `json:"interval" valid:"numeric"`
	Metadata            ConfigurationMetaData    `json:"metadata" valid:"required"`
	Thresholds          []ConfigurationThreshold `json:"thresholds" valid:"required"`
}

type ConfigurationMetaData struct {
	Monitoring string `json:"monitoring" valid:"required"`
	Team       string `json:"string" valid:"required"`
	Source     string `json:"source" valid:"required"`
	Targethost string `json:"targethost" valid:"required"`
}

type ConfigurationThreshold struct {
	Predicate string `json:"predicate" valid:"required"`
	Level     uint16 `json:"level" valid:"required"`
	Value     int64  `json:"value" valid:"required"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
