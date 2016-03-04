package somatree

type SomaTreeAttacher interface {
	Attach(a AttachRequest)
	Destroy()
	Detach()

	SomaTreePropertier
	Checker

	GetName() string
	setActionDeep(c chan *Action)

	clearParent()
	setFault(f *SomaTreeElemFault)
	setParent(p SomaTreeReceiver)
	updateFaultRecursive(f *SomaTreeElemFault)
	updateParentRecursive(p SomaTreeReceiver)
}

// implemented by: repository
type SomaTreeRootAttacher interface {
	SomaTreeAttacher
	attachToRoot(a AttachRequest)
}

// implemented by: buckets
type SomaTreeRepositoryAttacher interface {
	SomaTreeAttacher
	attachToRepository(a AttachRequest)
	CloneRepository() SomaTreeRepositoryAttacher
}

// implemented by: groups, clusters, nodes
type SomaTreeBucketAttacher interface {
	SomaTreeAttacher
	attachToBucket(a AttachRequest)
	CloneBucket() SomaTreeBucketAttacher
	ReAttach(a AttachRequest)
}

// implemented by: groups, clusters, nodes
type SomaTreeGroupAttacher interface {
	SomaTreeAttacher
	attachToGroup(a AttachRequest)
	CloneGroup() SomaTreeGroupAttacher
	ReAttach(a AttachRequest)
}

// implemented by: nodes
type SomaTreeClusterAttacher interface {
	SomaTreeAttacher
	attachToCluster(a AttachRequest)
	CloneCluster() SomaTreeClusterAttacher
	ReAttach(a AttachRequest)
}

type AttachRequest struct {
	Root       SomaTreeReceiver
	ParentType string
	ParentId   string
	ParentName string
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
