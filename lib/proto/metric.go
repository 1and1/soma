/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Metric struct {
	Path        string           `json:"path,omitempty"`
	Unit        string           `json:"unit,omitempty"`
	Description string           `json:"description,omitempty"`
	Packages    *[]MetricPackage `json:"packages,omitempty"`
	Details     *MetricDetails   `json:"details,omitempty"`
}

type MetricFilter struct {
	Unit     string `json:"unit,omitempty"`
	Provider string `json:"provider,omitempty"`
	Package  string `json:"package,omitempty"`
}

type MetricDetails struct {
	DetailsCreation
}

type MetricPackage struct {
	Provider string `json:"provider,omitempty"`
	Name     string `json:"name,omitempty"`
}

//
func (p *Metric) DeepCompare(a *Metric) bool {
	if p.Path != a.Path || p.Unit != a.Unit || p.Description != a.Description {
		return false
	}
packageloop:
	for _, pkg := range *p.Packages {
		if pkg.DeepCompareSlice(a.Packages) {
			continue packageloop
		}
		return false
	}
revpackageloop:
	for _, pkg := range *a.Packages {
		if pkg.DeepCompareSlice(p.Packages) {
			continue revpackageloop
		}
		return false
	}
	return true
}

func (p *MetricPackage) DeepCompare(a *MetricPackage) bool {
	if p.Provider != a.Provider || p.Name != a.Name {
		return false
	}
	return true
}

func (p *MetricPackage) DeepCompareSlice(a *[]MetricPackage) bool {
	if a == nil || *a == nil {
		return false
	}
	for _, pkg := range *a {
		if p.DeepCompare(&pkg) {
			return true
		}
	}
	return false
}

func NewMetricRequest() Request {
	return Request{
		Flags:  &Flags{},
		Metric: &Metric{},
	}
}

func NewMetricFilter() Request {
	return Request{
		Filter: &Filter{
			Metric: &MetricFilter{},
		},
	}
}

func NewMetricResult() Result {
	return Result{
		Errors:  &[]string{},
		Metrics: &[]Metric{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
