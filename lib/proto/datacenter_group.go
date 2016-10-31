/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type DatacenterGroup struct {
	Name    string                  `json:"name,omitempty"`
	Members *[]Datacenter           `json:"members,omitempty"`
	Details *DatacenterGroupDetails `json:"details,omitempty"`
}

type DatacenterGroupDetails struct {
	DetailsCreation
}

func NewDatacenterGroupRequest() Request {
	return Request{
		Flags:           &Flags{},
		DatacenterGroup: &DatacenterGroup{},
	}
}

func NewDatacenterGroupResult() Result {
	return Result{
		Errors:           &[]string{},
		DatacenterGroups: &[]DatacenterGroup{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
