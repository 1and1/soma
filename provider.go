package proto

type Provider struct {
	Name    string           `json:"name,omitempty"`
	Details *ProviderDetails `json:"details,omitempty"`
}

type ProviderFilter struct {
	Name string `json:"name,omitempty"`
}

type ProviderDetails struct {
	DetailsCreation
}

func NewProviderRequest() Request {
	return Request{
		Provider: &Provider{},
	}
}

func NewProviderFilter() Request {
	return Request{
		Filter: &Filter{
			Provider: &ProviderFilter{},
		},
	}
}

func NewProviderResult() Result {
	return Result{
		Errors:    &[]string{},
		Providers: &[]Provider{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
