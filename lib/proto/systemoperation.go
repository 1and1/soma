package proto

type SystemOperation struct {
	Request      string `json:"request,omitempty"`
	RepositoryId string `json:"repositoryId,omitempty"`
	RebuildLevel string `json:"rebuildLevel,omitempty"`
}

func NewSystemOperationRequest() Request {
	return Request{
		Flags:           &Flags{},
		SystemOperation: &SystemOperation{},
	}
}

func NewSystemOperationResult() Result {
	return Result{
		Errors:           &[]string{},
		SystemOperations: &[]SystemOperation{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
