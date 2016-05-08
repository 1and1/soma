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

func NewPermissionRequest() Request {
	return Request{
		Permission: &Permission{},
	}
}

func NewPermissionResult() Result {
	return Result{
		Errors:      &[]string{},
		Permissions: &[]Permission{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
