package somaproto

type Permission struct {
	Name     string             `json:"name, omitempty"`
	Category string             `json:"category, omitempty"`
	Grants   string             `json:"grants, omitempty"`
	Details  *PermissionDetails `json:"details, omitempty"`
}

type PermissionDetails struct {
	DetailsCreation
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
