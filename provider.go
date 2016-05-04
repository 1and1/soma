package somaproto

type Provider struct {
	Name    string           `json:"name, omitempty"`
	Details *ProviderDetails `json:"details, omitempty"`
}

type ProviderFilter struct {
	Name string `json:"name, omitempty"`
}

type ProviderDetails struct {
	DetailsCreation
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
