package somaproto

type HostDeployment struct {
	CheckInstanceId            string   `json:"checkInstanceId"`
	DeleteInstance             bool     `json:"deleteInstance"`
	CurrentCheckInstanceIdList []string `json:"currentCheckInstanceIdList, omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
