package somaproto

type HostDeploymentRequest struct {
	IdList []string `json:"id_list,omitempty"`
}

type HostDeploymentResult struct {
	Code        uint16              `json:"code,omitempty"`
	Status      string              `json:"status,omitempty"`
	Deployments []DeploymentDetails `json:"deployments,omitempty"`
	Delete      []string            `json:"delete,omitempty"`
	JobId       string              `json:"jobid,omitempty"`
}

func (hd *HostDeploymentResult) ErrorMark(err error, imp bool,
	found bool, length int, jobid string) bool {
	if hd.markError(err) {
		return true
	}
	if hd.markImplemented(imp) {
		return true
	}
	if hd.markFound(found, length) {
		return true
	}
	if hd.hasJobId(jobid) {
		return hd.markAccepted()
	}
	return hd.markOk()
}

func (hd *HostDeploymentResult) markError(err error) bool {
	if err != nil {
		hd.Code = 500
		hd.Status = "ERROR"
		return true
	}
	return false
}

func (hd *HostDeploymentResult) markImplemented(f bool) bool {
	if f {
		hd.Code = 501
		hd.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (hd *HostDeploymentResult) markFound(f bool, i int) bool {
	if f || i == 0 {
		hd.Code = 404
		hd.Status = "NOT FOUND"
		return true
	}
	return false
}

func (hd *HostDeploymentResult) markOk() bool {
	hd.Code = 200
	hd.Status = "OK"
	return false
}

func (hd *HostDeploymentResult) hasJobId(s string) bool {
	if s != "" {
		hd.JobId = s
		return true
	}
	return false
}

func (hd *HostDeploymentResult) markAccepted() bool {
	hd.Code = 202
	hd.Status = "ACCEPTED"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
