package somatree

type SomaTreeBucketeer interface {
	GetBucket() SomaTreeReceiver
	GetEnvironment() string
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
