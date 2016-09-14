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
