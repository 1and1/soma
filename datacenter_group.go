package somaproto

type DatacenterGroup struct {
	Name    string                  `json:"name, omitempty"`
	Members *[]Datacenter           `json:"members, omitempty"`
	Details *DatacenterGroupDetails `json:"details, omitempty"`
}

type DatacenterGroupDetails struct {
	DetailsCreation
}

func NewDatacenterGroupRequest() Request {
	return Request{
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
