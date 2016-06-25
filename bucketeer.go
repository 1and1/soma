package tree

type Bucketeer interface {
	GetBucket() Receiver
	GetEnvironment() string
	GetRepository() string
	GetRepositoryName() string
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
