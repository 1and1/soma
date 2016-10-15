package main

// Populate the node structure with data, overwriting the client
// submitted values.
func (g *guidePost) fillNode(q *treeRequest) error {
	var (
		err                      error
		ndName, ndTeam, ndServer string
		ndAsset                  int64
		ndOnline, ndDeleted      bool
	)
	if err = g.node_stmt.QueryRow(q.Node.Node.Id).Scan(
		&ndAsset,
		&ndName,
		&ndTeam,
		&ndServer,
		&ndOnline,
		&ndDeleted,
	); err != nil {
		return err
	}
	q.Node.Node.AssetId = uint64(ndAsset)
	q.Node.Node.Name = ndName
	q.Node.Node.TeamId = ndTeam
	q.Node.Node.ServerId = ndServer
	q.Node.Node.IsOnline = ndOnline
	q.Node.Node.IsDeleted = ndDeleted
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
