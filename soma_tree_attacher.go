package somatree

type SomaTreeAttacher interface {
	Attach(a AttachRequest)
	ReAttach(a AttachRequest)
	SetParent(p SomaTreeReceiver)
}

// implemented by: repository
type SomaTreeRootAttacher interface {
	SomaTreeAttacher
	GetName() string
	AttachToRoot(a AttachRequest)
}

// implemented by: buckets
type SomaTreeRepositoryAttacher interface {
	SomaTreeAttacher
	GetName() string
	AttachToRepository(a AttachRequest)
}

// implemented by: groups, clusters, nodes
type SomaTreeBucketAttacher interface {
	SomaTreeAttacher
	AttachToBucket(a AttachRequest)
}

// implemented by: groups, clusters, nodes
type SomaTreeGroupAttacher interface {
	SomaTreeAttacher
	AttachToGroup(a AttachRequest)
}

// implemented by: nodes
type SomaTreeClusterAttacher interface {
	SomaTreeAttacher
	AttachToCluster(a AttachRequest)
}

type AttachRequest struct {
	Root       SomaTreeReceiver
	ParentType string
	ParentId   string
	ParentName string
	ChildType  string
	ChildName  string
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
