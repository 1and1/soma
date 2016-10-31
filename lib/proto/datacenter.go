/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Datacenter struct {
	Locode  string             `json:"locode,omitempty"`
	Details *DatacenterDetails `json:"details,omitempty"`
}

type DatacenterDetails struct {
	DetailsCreation
}

func NewDatacenterRequest() Request {
	return Request{
		Flags:      &Flags{},
		Datacenter: &Datacenter{},
	}
}

func NewDatacenterResult() Result {
	return Result{
		Errors:      &[]string{},
		Datacenters: &[]Datacenter{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
