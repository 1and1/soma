package somaproto

type ProtoRequestEnvironment struct {
  Environment   string  `json:"environment,omitempty"`
}

type ProtoResultEnvironment struct {
  Code          uint16  `json:"code,omitempty"`
  Status        string  `json:"status,omitempty"`
  Text          []string  `json:"text,omitempty"`
}

type ProtoResultEnvironmentList struct {
  Code        uint16  `json:"code,omitempty"`
  Status      string  `json:"status,omitempty"`
  Text        []string  `json:"text,omitempty"`
  Environments  []string  `json:"views,omitempty"`
}

type ProtoResultEnvironmentDetail struct {
  Code        uint16  `json:"code,omitempty"`
  Status      string  `json:"status,omitempty"`
  Text        []string  `json:"text,omitempty"`
  Details     ProtoEnvironmentDetails  `json:"details,omitempty"`
}

type ProtoEnvironmentDetails struct {
  Environment string  `json:"environment,omitempty"`
  CreatedAt   string  `json:"createdat,omitempty"`
  CreatedBy   string  `json:"createdby,omitempty"`
  UsedBy      []string  `json:"usedby,omitempty"`
}
