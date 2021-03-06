/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Property struct {
	Type             string           `json:"type"`
	RepositoryId     string           `json:"repositoryId,omitempty"`
	BucketId         string           `json:"bucketId,omitempty"`
	InstanceId       string           `json:"instanceId,omitempty"`
	View             string           `json:"view,omitempty"`
	Inheritance      bool             `json:"inheritance,omitempty"`
	ChildrenOnly     bool             `json:"childrenOnly,omitempty"`
	IsInherited      bool             `json:"isInherited,omitempty"`
	SourceInstanceId string           `json:"sourceInstanceId,omitempty"`
	SourceType       string           `json:"sourceType,omitempty"`
	InheritedFrom    string           `json:"inheritedFrom,omitempty"`
	Custom           *PropertyCustom  `json:"custom,omitempty"`
	System           *PropertySystem  `json:"system,omitempty"`
	Service          *PropertyService `json:"service,omitempty"`
	Native           *PropertyNative  `json:"native,omitempty"`
	Oncall           *PropertyOncall  `json:"oncall,omitempty"`
	Details          *PropertyDetails `json:"details,omitempty"`
}

type PropertyFilter struct {
	Name         string `json:"name,omitempty"`
	Type         string `json:"type,omitempty"`
	RepositoryId string `json:"repositoryId,omitempty"`
}

type PropertyDetails struct {
	DetailsCreation
}

type PropertyCustom struct {
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	RepositoryId string `json:"repositoryId,omitempty"`
	Value        string `json:"value,omitempty"`
}

func (t *PropertyCustom) DeepCompare(a *PropertyCustom) bool {
	if t.Id != a.Id || t.Name != a.Name || t.RepositoryId != a.RepositoryId || t.Value != a.Value {
		return false
	}
	return true
}

func (t *PropertyCustom) DeepCompareSlice(a *[]PropertyCustom) bool {
	if a == nil || *a == nil {
		return false
	}
	for _, cust := range *a {
		if t.DeepCompare(&cust) {
			return true
		}
	}
	return false
}

type PropertySystem struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

func (t *PropertySystem) DeepCompare(a *PropertySystem) bool {
	if t.Name != a.Name || t.Value != a.Value {
		return false
	}
	return true
}

func (t *PropertySystem) DeepCompareSlice(a *[]PropertySystem) bool {
	if a == nil || *a == nil {
		return false
	}
	for _, sys := range *a {
		if t.DeepCompare(&sys) {
			return true
		}
	}
	return false
}

type PropertyService struct {
	Name       string             `json:"name,omitempty"`
	TeamId     string             `json:"teamId,omitempty"`
	Attributes []ServiceAttribute `json:"attributes"`
}

func (t *PropertyService) DeepCompare(a *PropertyService) bool {
	if t.Name != a.Name || t.TeamId != a.TeamId {
		return false
	}
attrloop:
	for _, attr := range t.Attributes {
		if attr.DeepCompareSlice(&a.Attributes) {
			continue attrloop
		}
		return false
	}
revattrloop:
	for _, attr := range a.Attributes {
		if attr.DeepCompareSlice(&t.Attributes) {
			continue revattrloop
		}
		return false
	}
	return true
}

type PropertyNative struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

func (t *PropertyNative) DeepCompare(a *PropertyNative) bool {
	if t.Name != a.Name || t.Value != a.Value {
		return false
	}
	return true
}

type PropertyOncall struct {
	Id     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Number string `json:"number,omitempty"`
}

func (t *PropertyOncall) DeepCompare(a *PropertyOncall) bool {
	if t.Id != a.Id || t.Name != a.Name || t.Number != a.Number {
		return false
	}
	return true
}

type ServiceAttribute struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

func (t *ServiceAttribute) DeepCompare(a *ServiceAttribute) bool {
	if t.Name != a.Name || t.Value != a.Value {
		return false
	}
	return true
}

func (t *ServiceAttribute) DeepCompareSlice(a *[]ServiceAttribute) bool {
	if a == nil || *a == nil {
		return false
	}
	for _, attr := range *a {
		if t.DeepCompare(&attr) {
			return true
		}
	}
	return false
}

func NewPropertyRequest() Request {
	return Request{
		Flags:    &Flags{},
		Property: &Property{},
	}
}

func NewPropertyFilter() Request {
	return Request{
		Filter: &Filter{
			Property: &PropertyFilter{},
		},
	}
}

func NewPropertyResult() Result {
	return Result{
		Errors:     &[]string{},
		Properties: &[]Property{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
