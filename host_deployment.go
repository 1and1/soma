package somaproto

type HostDeployment struct {
	CheckInstanceId            string   `json:"checkInstanceId"`
	DeleteInstance             bool     `json:"deleteInstance"`
	CurrentCheckInstanceIdList []string `json:"currentCheckInstanceIdList, omitempty"`
}

func NewHostDeploymentRequest() Request {
	return Request{
		HostDeployment: &HostDeployment{},
	}
}

func NewHostDeploymentResult() Result {
	return Result{
		Errors:          &[]string{},
		HostDeployments: &[]HostDeployment{},
		Deployments:     &[]Deployment{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
