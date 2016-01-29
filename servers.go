package somaproto

type ProtoRequestServer struct {
	Server  *ProtoServer       `json:"server,omitempty"`
	Filter  *ProtoServerFilter `json:"filter,omitempty"`
	Purge   bool               `json:"purge,omitempty"`
	Restore bool               `json:"restore,omitempty"`
}

type ProtoResultServer struct {
	Code    uint16        `json:"code,omitempty"`
	Status  string        `json:"status,omitempty"`
	Text    []string      `json:"text,omitempty"`
	Servers []ProtoServer `json:"servers,omitempty"`
	JobId   string        `json:"jobid,omitempty"`
}

type ProtoServer struct {
	Id         string              `json:"id,omitempty"`
	AssetId    uint64              `json:"assetid,omitempty"`
	Datacenter string              `json:"datacenter,omitempty"`
	Location   string              `json:"location,omitempty"`
	Name       string              `json:"name,omitempty"`
	IsOnline   bool                `json:"online,omitempty"`
	IsDeleted  bool                `json:"deleted,omitempty"`
	Details    *ProtoServerDetails `json:"details,omitempty"`
}

type ProtoServerDetails struct {
	CreatedAt string   `json:"createdat,omitempty"`
	CreatedBy string   `json:"createdby,omitempty"`
	Nodes     []string `json:"nodes,omitempty"`
}

type ProtoServerFilter struct {
	Online     bool   `json:"online,omitempty"`
	Deleted    bool   `json:"deleted,omitempty"`
	Datacenter string `json:"datacenter,omitempty"`
	Name       string `json:"name,omitempty"`
}

//
func (p *ProtoResultServer) ErrorMark(err error, imp bool, found bool, length int) bool {
	if p.markError(err) {
		return true
	}
	if p.markImplemented(imp) {
		return true
	}
	if p.markFound(found, length) {
		return true
	}
	return false
}

func (p *ProtoResultServer) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *ProtoResultServer) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *ProtoResultServer) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
