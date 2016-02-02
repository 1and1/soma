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
	ReceiveRepository(r ReceiveRequest)
}

type SomaTreeRepositoryUnlinker interface {
	SomaTreeUnlinker
	UnlinkRepository(u UnlinkRequest)
}

// implemented by: repositories
type SomaTreeBucketReceiver interface {
	SomaTreeReceiver
	SomaTreeBucketUnlinker
	ReceiveBucket(r ReceiveRequest)
}

type SomaTreeBucketUnlinker interface {
	SomaTreeUnlinker
	UnlinkBucket(u UnlinkRequest)
}

// implemented by: buckets, groups
type SomaTreeGroupReceiver interface {
	SomaTreeReceiver
	SomaTreeUnlinker
	ReceiveGroup(r ReceiveRequest)
	UnlinkGroup(u UnlinkRequest)
}

// implemented by: buckets, groups
type SomaTreeClusterReceiver interface {
	SomaTreeReceiver
	SomaTreeUnlinker
	ReceiveCluster(r ReceiveRequest)
	UnlinkCluster(u UnlinkRequest)
}

// implemented by: buckets, groups, clusters
type SomaTreeNodeReceiver interface {
	SomaTreeReceiver
	SomaTreeUnlinker
	ReceiveNode(r ReceiveRequest)
	UnlinkNode(u UnlinkRequest)
}

type ReceiveRequest struct {
	ParentType string
	ParentId   string
	ParentName string
	ChildType  string
	Repository *SomaTreeElemRepository
	Bucket     *SomaTreeElemBucket
	Group      *SomaTreeElemGroup
	Cluster    *SomaTreeElemCluster
	Node       *SomaTreeElemNode
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
