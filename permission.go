package proto

type Permission struct {
	Id       string             `json:"id,omitempty"`
	Name     string             `json:"name,omitempty"`
	Category string             `json:"category,omitempty"`
	Grants   string             `json:"grants,omitempty"`
	Details  *PermissionDetails `json:"details,omitempty"`
}

type PermissionDetails struct {
	CreatedAt string `json:"createdAt,omitempty"`
	CreatedBy string `json:"createdBy,omitempty"`
}

type PermissionFilter struct {
	Name string `json:"name,omitempty"`
}

func NewPermissionRequest() Request {
	return Request{
		Flags:      &Flags{},
		Permission: &Permission{},
	}
}

func NewPermissionFilter() Request {
	return Request{
		Filter: &Filter{
			Permission: &PermissionFilter{},
		},
	}
}

func NewPermissionResult() Result {
	return Result{
		Errors:      &[]string{},
		Permissions: &[]Permission{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
