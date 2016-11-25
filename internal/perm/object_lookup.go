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

// addRepository inserts a new repository into the cache
func (m *objectLookup) addRepository(repoID string) {
	if _, ok := m.byRepository[repoID]; ok {
		return
	}
	m.byRepository[repoID] = map[string][]string{}
}

// addBucket inserts a new bucket into the cache
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

// addGroup inserts a new group into the cache
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

// addCluster inserts a new cluster into the cache
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

// addNode inserts a new node into the cache
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

// rmRepository removes a repository from the cache
func (m *objectLookup) rmRepository(repoID string) {
	if _, ok := m.byRepository[repoID]; !ok {
		return
	}

	// make a local copy of the buckets to iterate over; deleting
	// the buckets will clean all groups/clusters/nodes in them
	buckets := make([]string, len(m.byRepository[repoID][`bucket`]))
	copy(buckets, m.byRepository[repoID][`bucket`])
	for _, b := range buckets {
		m.rmBucket(b)
	}

	// remove repository
	delete(m.byRepository, repoID)
}

// rmBucket removes a bucket from the cache
func (m *objectLookup) rmBucket(bucketID string) {
	if _, ok := m.byBucket[bucketID]; !ok {
		return
	}

	// make local copies to iterate over instead of
	// iterating over the slice we modify
	groups := make([]string, len(m.byBucket[bucketID][`group`]))
	copy(groups, m.byBucket[bucketID][`group`])
	for _, g := range groups {
		m.rmGroup(g)
	}

	clusters := make([]string, len(m.byBucket[bucketID][`cluster`]))
	copy(clusters, m.byBucket[bucketID][`cluster`])
	for _, c := range clusters {
		m.rmCluster(c)
	}

	nodes := make([]string, len(m.byBucket[bucketID][`node`]))
	copy(nodes, m.byBucket[bucketID][`node`])
	for _, n := range nodes {
		m.rmNode(n)
	}

	// remove bucket from repository
	repoID := m.byBucket[bucketID][`repository`][0]
	for i := range m.byRepository[repoID][`bucket`] {
		if bucketID != m.byRepository[repoID][`bucket`][i] {
			continue
		}
		m.byRepository[repoID][`bucket`] = append(
			m.byRepository[repoID][`bucket`][:i],
			m.byRepository[repoID][`bucket`][i+1:]...)
		m.compactionCounter++
	}

	// remove bucket
	delete(m.byBucket, bucketID)
}

// rmGroup removes a group from the cache
func (m *objectLookup) rmGroup(groupID string) {
	if _, ok := m.byGroup[groupID]; !ok {
		return
	}

	// get repoID/bucketID for group
	repoID := m.byGroup[groupID][`repository`][0]
	bucketID := m.byGroup[groupID][`bucket`][0]

	// remove group from bucket
	for i := range m.byBucket[bucketID][`group`] {
		if groupID != m.byBucket[bucketID][`group`][i] {
			continue
		}
		m.byBucket[bucketID][`group`] = append(
			m.byBucket[bucketID][`group`][:i],
			m.byBucket[bucketID][`group`][i+1:]...)
		m.compactionCounter++
	}

	// remove group from repository
	for i := range m.byRepository[repoID][`group`] {
		if groupID != m.byRepository[repoID][`group`][i] {
			continue
		}
		m.byRepository[repoID][`group`] = append(
			m.byRepository[repoID][`group`][:i],
			m.byRepository[repoID][`group`][i+1:]...)
		m.compactionCounter++
	}

	// remove group
	delete(m.byGroup, groupID)
}

// rmCluster removes a cluster from the cache
func (m *objectLookup) rmCluster(clusterID string) {
	if _, ok := m.byCluster[clusterID]; !ok {
		return
	}

	// get repoID/clusterID for cluster
	repoID := m.byCluster[clusterID][`repository`][0]
	bucketID := m.byCluster[clusterID][`bucket`][0]

	// remove cluster from bucket
	for i := range m.byBucket[bucketID][`cluster`] {
		if clusterID != m.byBucket[bucketID][`cluster`][i] {
			continue
		}
		m.byBucket[bucketID][`cluster`] = append(
			m.byBucket[bucketID][`cluster`][:i],
			m.byBucket[bucketID][`cluster`][i+1:]...)
		m.compactionCounter++
	}

	// remove cluster from repository
	for i := range m.byRepository[repoID][`cluster`] {
		if clusterID != m.byRepository[repoID][`cluster`][i] {
			continue
		}
		m.byRepository[repoID][`cluster`] = append(
			m.byRepository[repoID][`cluster`][:i],
			m.byRepository[repoID][`cluster`][i+1:]...)
		m.compactionCounter++
	}

	// remove cluster
	delete(m.byCluster, clusterID)
}

// rmNode removes a node from the cache
func (m *objectLookup) rmNode(nodeID string) {
	if _, ok := m.byNode[nodeID]; !ok {
		return
	}

	// get repoId for bucket (can only have one repository)
	repoID := m.byNode[nodeID][`repository`][0]
	bucketID := m.byNode[nodeID][`bucket`][0]

	// remove node from bucket
	for i := range m.byBucket[bucketID][`node`] {
		if nodeID != m.byBucket[bucketID][`node`][i] {
			continue
		}
		m.byBucket[bucketID][`node`] = append(
			m.byBucket[bucketID][`node`][:i],
			m.byBucket[bucketID][`node`][i+1:]...)
		m.compactionCounter++
	}

	// remove node from repository
	for i := range m.byRepository[repoID][`node`] {
		if nodeID != m.byRepository[repoID][`node`][i] {
			continue
		}
		m.byRepository[repoID][`node`] = append(
			m.byRepository[repoID][`node`][:i],
			m.byRepository[repoID][`node`][i+1:]...)
		m.compactionCounter++
	}

	// remove node
	delete(m.byNode, nodeID)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
