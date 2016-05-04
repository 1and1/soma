package somaproto

type Metric struct {
	Path        string           `json:"path, omitempty"`
	Unit        string           `json:"unit, omitempty"`
	Description string           `json:"description, omitempty"`
	Packages    *[]MetricPackage `json:"packages, omitempty"`
	Details     *MetricDetails   `json:"details, omitempty"`
}

type MetricFilter struct {
	Unit     string `json:"unit, omitempty"`
	Provider string `json:"provider, omitempty"`
	Package  string `json:"package, omitempty"`
}

type MetricDetails struct {
	DetailsCreation
}

type MetricPackage struct {
	Provider string `json:"provider, omitempty"`
	Name     string `json:"name, omitempty"`
}

//
func (p *Metric) DeepCompare(a *Metric) bool {
	if p.Path != a.Path || p.Unit != a.Unit || p.Description != a.Description {
		return false
	}
packageloop:
	for _, pkg := range *p.Packages {
		if pkg.DeepCompareSlice(a.Packages) {
			continue packageloop
		}
		return false
	}
	return true
}

func (p *MetricPackage) DeepCompare(a *MetricPackage) bool {
	if p.Provider != a.Provider || p.Name != a.Name {
		return false
	}
	return true
}

func (p *MetricPackage) DeepCompareSlice(a *[]MetricPackage) bool {
	for _, pkg := range *a {
		if p.DeepCompare(&pkg) {
			return true
		}
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
