package somaproto

type Repository struct {
	Id         string             `json:"id,omitempty"`
	Name       string             `json:"name,omitempty"`
	TeamId     string             `json:"teamId,omitempty"`
	IsDeleted  bool               `json:"isDeleted,omitempty"`
	IsActive   bool               `json:"isActive,omitempty"`
	Details    *RepositoryDetails `json:"details,omitempty"`
	Properties *[]Property        `json:"properties,omitempty"`
}

type RepositoryFilter struct {
	Name      string `json:"name,omitempty"`
	TeamId    string `json:"teamId,omitempty"`
	IsDeleted bool   `json:"isDeleted,omitempty"`
	IsActive  bool   `json:"isActive,omitempty"`
}

type RepositoryDetails struct {
	DetailsCreation
}

func NewRepositoryRequest() Request {
	return Request{
		Repository: &Repository{},
	}
}

func NewRepositoryFilter() Request {
	return Request{
		Filter: &Filter{
			Repository: &RepositoryFilter{},
		},
	}
}

func NewRepositoryResult() Result {
	return Result{
		Errors:       &[]string{},
		Repositories: &[]Repository{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
