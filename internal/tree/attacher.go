/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import log "github.com/Sirupsen/logrus"

type Attacher interface {
	Propertier
	Checker

	Attach(a AttachRequest)
	Destroy()
	Detach()
	GetName() string
	ComputeCheckInstances()
	ClearLoadInfo()

	setActionDeep(c chan *Action)
	setLoggerDeep(l *log.Logger)

	clearParent()
	setFault(f *Fault)
	setParent(p Receiver)
	updateFaultRecursive(f *Fault)
	updateParentRecursive(p Receiver)
}

// implemented by: repository
type RootAttacher interface {
	Attacher

	attachToRoot(a AttachRequest)
	setLog(l *log.Logger)
}

// implemented by: buckets
type RepositoryAttacher interface {
	Attacher

	CloneRepository() RepositoryAttacher

	attachToRepository(a AttachRequest)
}

// implemented by: groups, clusters, nodes
type BucketAttacher interface {
	Attacher

	CloneBucket() BucketAttacher
	ReAttach(a AttachRequest)

	attachToBucket(a AttachRequest)
}

// implemented by: groups, clusters, nodes
type GroupAttacher interface {
	Attacher

	CloneGroup() GroupAttacher
	ReAttach(a AttachRequest)

	attachToGroup(a AttachRequest)
}

// implemented by: nodes
type ClusterAttacher interface {
	Attacher

	CloneCluster() ClusterAttacher
	ReAttach(a AttachRequest)

	attachToCluster(a AttachRequest)
}

type AttachRequest struct {
	Root       Receiver
	ParentType string
	ParentId   string
	ParentName string
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
