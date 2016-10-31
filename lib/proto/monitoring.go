/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Monitoring struct {
	Id       string             `json:"id,omitempty"`
	Name     string             `json:"name,omitempty"`
	Mode     string             `json:"mode,omitempty"`
	Contact  string             `json:"contact,omitempty"`
	TeamId   string             `json:"teamId,omitempty"`
	Callback string             `json:"callback,omitempty"`
	Details  *MonitoringDetails `json:"details,omitempty"`
}

type MonitoringFilter struct {
	Name    string `json:"name,omitempty"`
	Mode    string `json:"mode,omitempty"`
	Contact string `json:"contact,omitempty"`
	TeamId  string `json:"teamId,omitempty"`
}

type MonitoringDetails struct {
	DetailsCreation
	Checks    uint64 `json:"checks,omitempty"`
	Instances uint64 `json:"instances,omitempty"`
}

func (p *Monitoring) DeepCompare(a *Monitoring) bool {
	if p.Id != a.Id || p.Name != a.Name || p.Mode != a.Mode ||
		p.Contact != a.Contact || p.TeamId != a.TeamId || p.Callback != a.Callback {
		return false
	}
	return true
}

func NewMonitoringRequest() Request {
	return Request{
		Flags:      &Flags{},
		Monitoring: &Monitoring{},
	}
}

func NewMonitoringFilter() Request {
	return Request{
		Filter: &Filter{
			Monitoring: &MonitoringFilter{},
		},
	}
}

func NewMonitoringResult() Result {
	return Result{
		Errors:      &[]string{},
		Monitorings: &[]Monitoring{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
