package somatree

type Receiver interface {
	Receive(r ReceiveRequest)
}

type Unlinker interface {
	Unlink(u UnlinkRequest)
}

// implemented by: root
type RepositoryReceiver interface {
	Receiver
	RepositoryUnlinker

	receiveRepository(r ReceiveRequest)
}

type RepositoryUnlinker interface {
	Unlinker

	unlinkRepository(u UnlinkRequest)
}

// implemented by: repositories
type BucketReceiver interface {
	Receiver
	BucketUnlinker

	receiveBucket(r ReceiveRequest)
}

type BucketUnlinker interface {
	Unlinker

	unlinkBucket(u UnlinkRequest)
}

type FaultReceiver interface {
	Receiver
	FaultUnlinker

	receiveFault(r ReceiveRequest)
}

type FaultUnlinker interface {
	Unlinker

	unlinkFault(u UnlinkRequest)
}

// implemented by: buckets, groups
type GroupReceiver interface {
	Receiver
	GroupUnlinker

	receiveGroup(r ReceiveRequest)
}

type GroupUnlinker interface {
	Unlinker

	unlinkGroup(u UnlinkRequest)
}

// implemented by: buckets, groups
type ClusterReceiver interface {
	Receiver
	ClusterUnlinker

	receiveCluster(r ReceiveRequest)
}

type ClusterUnlinker interface {
	Unlinker

	unlinkCluster(u UnlinkRequest)
}

// implemented by: buckets, groups, clusters
type NodeReceiver interface {
	Receiver
	NodeUnlinker

	receiveNode(r ReceiveRequest)
}

type NodeUnlinker interface {
	Unlinker

	unlinkNode(u UnlinkRequest)
}

//
type ReceiveRequest struct {
	ParentType string
	ParentId   string
	ParentName string
	ChildType  string
	Repository *Repository
	Bucket     *Bucket
	Group      *Group
	Cluster    *Cluster
	Node       *Node
	Fault      *Fault
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
