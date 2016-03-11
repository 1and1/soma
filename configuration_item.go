package main

import "github.com/satori/go.uuid"

type ConfigurationData struct {
	Configurations []ConfigurationItem `json:"configurations"`
}

type ConfigurationItem struct {
	ConfigurationItemId uuid.UUID                `json:"configuration_item_id"`
	Metric              string                   `json:"metric"`
	HostId              string                   `json:"host_id"`
	Tags                []string                 `json:"tags,omitempty"`
	Oncall              string                   `json:"oncall"`
	Interval            uint64                   `json:"interval"`
	Metadata            ConfigurationMetaData    `json:"metadata"`
	Thresholds          []ConfigurationThreshold `json:"thresholds"`
}

type ConfigurationMetaData struct {
	Monitoring string `json:"monitoring"`
	Team       string `json:"string"`
	Source     string `json:"source"`
	Targethost string `json:"targethost"`
}

type ConfigurationThreshold struct {
	Predicate string `json:"predicate"`
	Level     int    `json:"level"`
	Value     int    `json:"value"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
