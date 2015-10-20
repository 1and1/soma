package somaproto

type ProtoRequestPermission struct {
	Permission     string               `json:"permission,omitempty"`
	PermissionType string               `json:"permissiontype,omitempty"`
	GrantEnabled   bool                 `json:"grantenabled,omitempty"`
	Grant          ProtoPermissionGrant `json:"grant,omitempty"`
}

type ProtoResultPermission struct {
	Code           uint16                 `json:"code,omitempty"`
	Status         string                 `json:"status,omitempty"`
	Text           []string               `json:"text,omitempty"`
	Details        ProtoPermissionDetails `json:"details,omitempty"`
	PermissionList []string               `json:"permissionlist,omitempty"`
	UserList       []string               `json:"userlist,omitempty"`
}

type ProtoPermissionDetails struct {
	CreatedAt      string `json:"createdat,omitempty"`
	CreatedBy      string `json:"createdby,omitempty"`
	PermissionType string `json:"permissiontype,omitempty"`
}

type ProtoPermissionGrant struct {
	GrantType  string `json:"granttype,omitempty"`
	Repository string `json:"repository,omitempty"`
	Bucket     string `json:"bucket,omitempty"`
	Group      string `json:"group,omitempty"`
	Cluster    string `json:"cluster,omitempty"`
}
