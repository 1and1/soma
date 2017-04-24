/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"

	"github.com/1and1/soma/lib/proto"
)

// customProperties adds the custom properties to the node result
func (r *NodeRead) customProperties(node *proto.Node) error {
	var (
		rows                               *sql.Rows
		err                                error
		instanceID, sourceInstanceID, view string
		value, customProp, customID        string
	)

	if rows, err = r.stmtPropCustom.Query(
		node.Id,
	); err != nil {
		return err
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&sourceInstanceID,
			&view,
			&customID,
			&value,
			&customProp,
		); err != nil {
			rows.Close()
			return err
		}
		*node.Properties = append(*node.Properties, proto.Property{
			Type:             `custom`,
			RepositoryId:     node.Config.RepositoryId,
			BucketId:         node.Config.BucketId,
			InstanceId:       instanceID,
			SourceInstanceId: sourceInstanceID,
			View:             view,
			Custom: &proto.PropertyCustom{
				Id:    customID,
				Name:  customProp,
				Value: value,
			},
		})
	}
	err = rows.Err()
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
