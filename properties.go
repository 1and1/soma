package somaproto

type ProtoRequestProperty struct {
	Custom  *ProtoPropertyCustom  `json:"custom,omitempty"`
	System  *ProtoPropertySystem  `json:"system,omitempty"`
	Service *ProtoPropertyService `json:"service,omitempty"`
	Native  *ProtoPropertyNative  `json:"native,omitempty"`
	Filter  *ProtoPropertyFilter  `json:"filter,omitempty"`
}

type ProtoResultProperty struct {
	Code    uint16                 `json:"code,omitempty"`
	Status  string                 `json:"status,omitempty"`
	Text    []string               `json:"text,omitempty"`
	Custom  []ProtoPropertyCustom  `json:"custom,omitempty"`
	System  []ProtoPropertySystem  `json:"system,omitempty"`
	Service []ProtoPropertyService `json:"service,omitempty"`
	Native  []ProtoPropertyNative  `json:"native,omitempty"`
	JobId   string                 `json:"jobid,omitempty"`
}

type ProtoPropertyCustom struct {
	Id         string `json:"id,omitempty"`
	Repository string `json:"repository,omitempty"`
	Property   string `json:"property,omitempty"`
	Value      string `json:"value,omitempty"`
}

type ProtoPropertySystem struct {
	Property string `json:"property,omitempty"`
	Value    string `json:"value,omitempty"`
}

type ProtoPropertyService struct {
	Property   string                 `json:"property,omitempty"`
	Team       string                 `json:"team,omitempty"`
	Attributes ProtoServiceAttributes `json:"attributes,omitempty"`
}

type ProtoPropertyNative struct {
	Property string `json:"property,omitempty"`
}

type ProtoServiceAttributes struct {
	ProtoTransport   []string `json:"proto_transport,omitempty"`
	ProtoApplication []string `json:"proto_application,omitempty"`
	Port             []string `json:"port,omitempty"`
	ProcessComm      []string `json:"process_comm,omitempty"`
	ProcessArgs      []string `json:"process_args,omitempty"`
	FilePath         []string `json:"file_path,omitempty"`
	DirectoryPath    []string `json:"directory_path,omitempty"`
	UnixSocketPath   []string `json:"unix_socket_path,omitempty"`
	Uid              []string `json:"uid,omitempty"`
	Tls              []string `json:"tls,omitempty"`
	SoftwareProvider []string `json:"software_provider,omitempty"`
}

type ProtoPropertyFilter struct {
	Property string `json:"property,omitempty"`
	Type     string `json:"type,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
