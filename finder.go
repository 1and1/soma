package tree

type Finder interface {
	Find(f FindRequest, b bool) Attacher
}

type FindRequest struct {
	ElementId   string
	ElementName string
	ElementType string
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
