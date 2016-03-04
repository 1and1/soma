package somatree

type SomaTreeAttacher interface {
	Attach(a AttachRequest)
	Destroy()
	Detach()

	SomaTreePropertier
	Checker

	clearParent()
	setFault(f *SomaTreeElemFault)
	setParent(p SomaTreeReceiver)
	updateFaultRecursive(f *SomaTreeElemFault)
	updateParentRecursive(p SomaTreeReceiver)
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
	CloneRepository() SomaTreeRepositoryAttacher
	setActionDeep(c chan *Action)
}

// implemented by: groups, clusters, nodes
type SomaTreeBucketAttacher interface {
	SomaTreeAttacher
	GetName() string
	attachToBucket(a AttachRequest)
	CloneBucket() SomaTreeBucketAttacher
	ReAttach(a AttachRequest)
	setActionDeep(c chan *Action)
}

// implemented by: groups, clusters, nodes
type SomaTreeGroupAttacher interface {
	SomaTreeAttacher
	GetName() string
	attachToGroup(a AttachRequest)
	CloneGroup() SomaTreeGroupAttacher
	ReAttach(a AttachRequest)
	setActionDeep(c chan *Action)
}

// implemented by: nodes
type SomaTreeClusterAttacher interface {
	SomaTreeAttacher
	GetName() string
	attachToCluster(a AttachRequest)
	CloneCluster() SomaTreeClusterAttacher
	ReAttach(a AttachRequest)
	setActionDeep(c chan *Action)
}

type AttachRequest struct {
	Root       SomaTreeReceiver
	ParentType string
	ParentId   string
	ParentName string
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
