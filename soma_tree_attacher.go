package somatree

type SomaTreeAttacher interface {
	Attach(a AttachRequest)
	ReAttach(a AttachRequest)
	Destroy()
	setParent(p SomaTreeReceiver)
}

// implemented by: repository
type SomaTreeRootAttacher interface {
	SomaTreeAttacher
	GetName() string
	attachToRoot(a AttachRequest)
}

// implemented by: buckets
type SomaTreeRepositoryAttacher interface {
	SomaTreeAttacher
	GetName() string
	attachToRepository(a AttachRequest)
}

// implemented by: groups, clusters, nodes
type SomaTreeBucketAttacher interface {
	SomaTreeAttacher
	GetName() string
	attachToBucket(a AttachRequest)
}

// implemented by: groups, clusters, nodes
type SomaTreeGroupAttacher interface {
	SomaTreeAttacher
	GetName() string
	attachToGroup(a AttachRequest)
}

// implemented by: nodes
type SomaTreeClusterAttacher interface {
	SomaTreeAttacher
	GetName() string
	attachToCluster(a AttachRequest)
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
