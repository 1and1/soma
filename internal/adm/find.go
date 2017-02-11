/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package adm

import (
	"fmt"

	"github.com/1and1/soma/lib/proto"
)

// FindRepoPropSrcID fetches the source id of a locally set
// property on a repository
func FindRepoPropSrcID(pType, pName, view, repoID string,
	id *string) error {
	var (
		err  error
		res  *proto.Result
		repo proto.Repository
	)
	res, err = fetchObjList(fmt.Sprintf("/repository/%s", repoID))
	if err != nil {
		goto abort
	}

	if res.Repositories == nil || len(*res.Repositories) == 0 {
		err = fmt.Errorf("Resultset is empty")
		goto abort
	}

	repo = (*res.Repositories)[0]
	if repo.Properties == nil || len(*repo.Properties) == 0 {
		err = fmt.Errorf("Received no properties on repository")
		goto abort
	}

	return findPropSrcID(pType, pName, view, *repo.Properties, id)

abort:
	return fmt.Errorf("Failed to find source property: %s",
		err.Error())
}

// FindBucketPropSrcID fetches the source id of a locally set
// property on a bucket
func FindBucketPropSrcID(pType, pName, view, bucketID string,
	id *string) error {
	var (
		err    error
		res    *proto.Result
		bucket proto.Bucket
	)
	res, err = fetchObjList(fmt.Sprintf("/buckets/%s", bucketID))
	if err != nil {
		goto abort
	}

	if res.Buckets == nil || len(*res.Buckets) == 0 {
		err = fmt.Errorf("Resultset is empty")
		goto abort
	}

	bucket = (*res.Buckets)[0]
	if bucket.Properties == nil || len(*bucket.Properties) == 0 {
		err = fmt.Errorf("Received no properties on bucket")
		goto abort
	}

	return findPropSrcID(pType, pName, view, *bucket.Properties, id)

abort:
	return fmt.Errorf("Failed to find source property: %s",
		err.Error())
}

// FindGroupPropSrcID fetches the source id of a locally set
// property on a group
func FindGroupPropSrcID(pType, pName, view, groupID string,
	id *string) error {
	var (
		err   error
		res   *proto.Result
		group proto.Group
	)
	res, err = fetchObjList(fmt.Sprintf("/groups/%s", groupID))
	if err != nil {
		goto abort
	}

	if res.Groups == nil || len(*res.Groups) == 0 {
		err = fmt.Errorf("Resultset is empty")
		goto abort
	}

	group = (*res.Groups)[0]
	if group.Properties == nil || len(*group.Properties) == 0 {
		err = fmt.Errorf("Received no properties on group")
		goto abort
	}

	return findPropSrcID(pType, pName, view, *group.Properties, id)

abort:
	return fmt.Errorf("Failed to find source property: %s",
		err.Error())
}

// FindClusterPropSrcID fetches the source id of a locally set
// property on a cluster
func FindClusterPropSrcID(pType, pName, view, clusterID string,
	id *string) error {
	var (
		err     error
		res     *proto.Result
		cluster proto.Cluster
	)
	res, err = fetchObjList(fmt.Sprintf("/clusters/%s", clusterID))
	if err != nil {
		goto abort
	}

	if res.Clusters == nil || len(*res.Clusters) == 0 {
		err = fmt.Errorf("Resultset is empty")
		goto abort
	}

	cluster = (*res.Clusters)[0]
	if cluster.Properties == nil || len(*cluster.Properties) == 0 {
		err = fmt.Errorf("Received no properties on cluster")
		goto abort
	}

	return findPropSrcID(pType, pName, view, *cluster.Properties, id)

abort:
	return fmt.Errorf("Failed to find source property: %s",
		err.Error())
}

// FindNodePropSrcID fetches the source id of a locally set
// property on a node
func FindNodePropSrcID(pType, pName, view, nodeID string,
	id *string) error {
	var (
		err  error
		res  *proto.Result
		node proto.Node
	)
	res, err = fetchObjList(fmt.Sprintf("/nodes/%s", nodeID))
	if err != nil {
		goto abort
	}

	if res.Nodes == nil || len(*res.Nodes) == 0 {
		err = fmt.Errorf("Resultset is empty")
		goto abort
	}

	node = (*res.Nodes)[0]
	if node.Properties == nil || len(*node.Properties) == 0 {
		err = fmt.Errorf("Received no properties on node")
		goto abort
	}

	return findPropSrcID(pType, pName, view, *node.Properties, id)

abort:
	return fmt.Errorf("Failed to find source property: %s",
		err.Error())
}

// findPropSrcId browses through the provided slice of Properties
// and returns the source id of the requested one
func findPropSrcID(pType, pName, view string, props []proto.Property,
	id *string) error {

	for _, p := range props {
		// wrong Type
		if p.Type != pType {
			continue
		}

		// wrong view
		if p.View != view {
			continue
		}

		// inherited property
		if p.InstanceId != p.SourceInstanceId {
			continue
		}

		switch pType {
		case `system`:
			if p.System.Name == pName {
				*id = p.SourceInstanceId
				return nil
			}
		case `oncall`:
			if p.Oncall.Name == pName {
				*id = p.SourceInstanceId
				return nil
			}
		case `custom`:
			if p.Custom.Name == pName {
				*id = p.SourceInstanceId
				return nil
			}
		case `service`:
			if p.Service.Name == pName {
				*id = p.SourceInstanceId
				return nil
			}
		}
	}

	return fmt.Errorf("Failed to find source property.")
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
