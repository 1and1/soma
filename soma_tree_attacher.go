package somatree

type SomaTreeAttacher interface {
	SomaTreePropertier
	Checker

	Attach(a AttachRequest)
	Destroy()
	Detach()
	GetName() string
	ComputeCheckInstances()
	ClearLoadInfo()

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

	CloneRepository() SomaTreeRepositoryAttacher

	attachToRepository(a AttachRequest)
}

// implemented by: groups, clusters, nodes
type SomaTreeBucketAttacher interface {
	SomaTreeAttacher

	CloneBucket() SomaTreeBucketAttacher
	ReAttach(a AttachRequest)

	attachToBucket(a AttachRequest)
}

// implemented by: groups, clusters, nodes
type SomaTreeGroupAttacher interface {
	SomaTreeAttacher

	CloneGroup() SomaTreeGroupAttacher
	ReAttach(a AttachRequest)

	attachToGroup(a AttachRequest)
}

// implemented by: nodes
type SomaTreeClusterAttacher interface {
	SomaTreeAttacher

	CloneCluster() SomaTreeClusterAttacher
	ReAttach(a AttachRequest)

	attachToCluster(a AttachRequest)
}

type AttachRequest struct {
	Root       SomaTreeReceiver
	ParentType string
	ParentId   string
	ParentName string
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
