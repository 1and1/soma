package somatree

type SomaTreeReceiver interface {
	Receive(r ReceiveRequest)
}

type SomaTreeUnlinker interface {
	Unlink(u UnlinkRequest)
}

// implemented by: root
type SomaTreeRepositoryReceiver interface {
	SomaTreeReceiver
	SomaTreeRepositoryUnlinker

	receiveRepository(r ReceiveRequest)
}

type SomaTreeRepositoryUnlinker interface {
	SomaTreeUnlinker

	unlinkRepository(u UnlinkRequest)
}

// implemented by: repositories
type SomaTreeBucketReceiver interface {
	SomaTreeReceiver
	SomaTreeBucketUnlinker

	receiveBucket(r ReceiveRequest)
}

type SomaTreeBucketUnlinker interface {
	SomaTreeUnlinker

	unlinkBucket(u UnlinkRequest)
}

type SomaTreeFaultReceiver interface {
	SomaTreeReceiver
	SomaTreeFaultUnlinker

	receiveFault(r ReceiveRequest)
}

type SomaTreeFaultUnlinker interface {
	SomaTreeUnlinker

	unlinkFault(u UnlinkRequest)
}

// implemented by: buckets, groups
type SomaTreeGroupReceiver interface {
	SomaTreeReceiver
	SomaTreeGroupUnlinker

	receiveGroup(r ReceiveRequest)
}

type SomaTreeGroupUnlinker interface {
	SomaTreeUnlinker

	unlinkGroup(u UnlinkRequest)
}

// implemented by: buckets, groups
type SomaTreeClusterReceiver interface {
	SomaTreeReceiver
	SomaTreeClusterUnlinker

	receiveCluster(r ReceiveRequest)
}

type SomaTreeClusterUnlinker interface {
	SomaTreeUnlinker

	unlinkCluster(u UnlinkRequest)
}

// implemented by: buckets, groups, clusters
type SomaTreeNodeReceiver interface {
	SomaTreeReceiver
	SomaTreeNodeUnlinker

	receiveNode(r ReceiveRequest)
}

type SomaTreeNodeUnlinker interface {
	SomaTreeUnlinker

	unlinkNode(u UnlinkRequest)
}

//
type ReceiveRequest struct {
	ParentType string
	ParentId   string
	ParentName string
	ChildType  string
	Repository *SomaTreeElemRepository
	Bucket     *Bucket
	Group      *SomaTreeElemGroup
	Cluster    *SomaTreeElemCluster
	Node       *SomaTreeElemNode
	Fault      *SomaTreeElemFault
}

type UnlinkRequest struct {
	ParentType string
	ParentId   string
	ParentName string
	ChildType  string
	ChildName  string
	ChildId    string
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
