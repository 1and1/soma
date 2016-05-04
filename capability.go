package somaproto

type Capability struct {
	Id           string             `json:"id, omitempty"`
	Name         string             `json:"name, omitempty"`
	MonitoringId string             `json:"monitoringId, omitempty"`
	Metric       string             `json:"metric, omitempty"`
	View         string             `json:"view, omitempty"`
	Thresholds   uint64             `json:"thresholds, omitempty"`
	Details      *CapabilityDetails `json:"details, omitempty"`
}

type CapabilityFilter struct {
	MonitoringId   string `json:"monitoringId, omitempty"`
	MonitoringName string `json:"monitoringName, omitempty"`
	Metric         string `json:"metric, omitempty"`
	View           string `json:"view, omitempty"`
}

type CapabilityDetails struct {
	DetailsCreation
}

func (c *Capability) DeepCompare(a *Capability) bool {
	if c.Id != a.Id {
		return false
	}
	if c.Name != a.Name {
		return false
	}
	if c.MonitoringId != a.MonitoringId {
		return false
	}
	if c.Metric != a.Metric {
		return false
	}
	if c.View != a.View {
		return false
	}
	if c.Thresholds != a.Thresholds {
		return false
	}
	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
