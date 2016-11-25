/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package perm

// objectLookup is the cache data structure that allows lookup
// of basic information about the repository trees.
type objectLookup struct {
	//
	compactionCounter int64
	// repositoryID -> objectType -> []objectID
	byRepository map[string]map[string][]string
	// bucketID -> objectType -> []objectID
	byBucket map[string]map[string][]string
	// groupID -> objectType -> []objectID
	byGroup map[string]map[string][]string
	// clusterID -> objectType -> []objectID
	byCluster map[string]map[string][]string
	// nodeID -> objectType -> []objectID
	byNode map[string]map[string][]string
}

// addRepository adds a new repository to the cache
func (m *objectLookup) addRepository(repoID string) {
	if _, ok := m.byRepository[repoID]; ok {
		return
	}
	m.byRepository[repoID] = map[string][]string{}
}

// addBucket adds and links a new bucket to the cache
func (m *objectLookup) addBucket(repoID, bucketID string) {
	if _, ok := m.byBucket[bucketID]; ok {
		return
	}

	// insert new bucket
	m.byBucket[bucketID] = map[string][]string{
		`repository`: []string{repoID},
	}

	// add bucket to repository
	if _, ok := m.byRepository[repoID][`bucket`]; !ok {
		m.byRepository[repoID][`bucket`] = []string{}
	}
	m.byRepository[repoID][`bucket`] = append(
		m.byRepository[repoID][`bucket`], bucketID)
}

func (m *objectLookup) addGroup(bucketID, groupID string) {
	if _, ok := m.byGroup[groupID]; ok {
		return
	}

	// insert new group
	repoID := m.byBucket[bucketID][`repository`][0]
	m.byGroup[groupID] = map[string][]string{
		`repository`: []string{repoID},
		`bucket`:     []string{bucketID},
	}

	// add group to bucket
	if _, ok := m.byBucket[bucketID][`group`]; !ok {
		m.byBucket[bucketID][`group`] = []string{}
	}
	m.byBucket[bucketID][`group`] = append(
		m.byBucket[bucketID][`group`], groupID)

	// add group to repository
	if _, ok := m.byRepository[repoID][`group`]; !ok {
		m.byRepository[repoID][`group`] = []string{}
	}
	m.byRepository[repoID][`group`] = append(
		m.byRepository[repoID][`group`], groupID)
}

func (m *objectLookup) addCluster(bucketID, clusterID string) {
	if _, ok := m.byCluster[clusterID]; ok {
		return
	}

	// insert new cluster
	repoID := m.byBucket[bucketID][`repository`][0]
	m.byCluster[clusterID] = map[string][]string{
		`repository`: []string{repoID},
		`bucket`:     []string{bucketID},
	}

	// add cluster to bucket
	if _, ok := m.byBucket[bucketID][`cluster`]; !ok {
		m.byBucket[bucketID][`cluster`] = []string{}
	}
	m.byBucket[bucketID][`cluster`] = append(
		m.byBucket[bucketID][`cluster`], clusterID)

	// add cluster to repository
	if _, ok := m.byRepository[repoID][`cluster`]; !ok {
		m.byRepository[repoID][`cluster`] = []string{}
	}
	m.byRepository[repoID][`cluster`] = append(
		m.byRepository[repoID][`cluster`], clusterID)
}

func (m *objectLookup) addNode(bucketID, nodeID string) {
	if _, ok := m.byNode[nodeID]; ok {
		return
	}

	// insert new node
	repoID := m.byBucket[bucketID][`repository`][0]
	m.byNode[nodeID] = map[string][]string{
		`repository`: []string{repoID},
		`bucket`:     []string{bucketID},
	}

	// add node to bucket
	if _, ok := m.byBucket[bucketID][`node`]; !ok {
		m.byBucket[bucketID][`node`] = []string{}
	}
	m.byBucket[bucketID][`node`] = append(
		m.byBucket[bucketID][`node`], nodeID)

	// add node to repository
	if _, ok := m.byRepository[repoID][`node`]; !ok {
		m.byRepository[repoID][`node`] = []string{}
	}
	m.byRepository[repoID][`node`] = append(
		m.byRepository[repoID][`node`], nodeID)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix