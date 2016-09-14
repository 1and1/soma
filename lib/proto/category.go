package proto

type Category struct {
	Name    string           `json:"name,omitempty"`
	Details *CategoryDetails `json:"details,omitempty"`
}

type CategoryDetails struct {
	CreatedAt string `json:"createdAt,omitempty"`
	CreatedBy string `json:"createdBy,omitempty"`
}

func NewCategoryRequest() Request {
	return Request{
		Flags:    &Flags{},
		Category: &Category{},
	}
}

func NewCategoryResult() Result {
	return Result{
		Errors:     &[]string{},
		Categories: &[]Category{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
