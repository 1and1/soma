package somatree

import (
	"testing"

	"github.com/satori/go.uuid"
)

func TestSetProperty(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`

	propUUID, _ := uuid.FromString(propId)

	// create tree
	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	// set property on repository
	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalue`,
	})

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}

	if len(actionChan) != 14 {
		t.Error(
			`Expected 14 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`repository`, `property_new`},
		[]string{`bucket`, `create`},
		[]string{`bucket`, `property_new`},
		[]string{`group`, `create`},
		[]string{`group`, `property_new`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `property_new`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`node`, `property_new`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	if len(sTree.Child.PropertySystem) != 1 {
		t.Error(
			`Exptected 1 system property, got `,
			len(sTree.Child.PropertySystem),
		)
	}

	if _, ok := sTree.Child.PropertySystem[propId]; !ok {
		t.Error(
			`Could not find property under requested id`,
		)
	}

	if len(sTree.Child.Children[buckId].(*Bucket).PropertySystem) != 1 {
		t.Error(
			`Expected bucket to have 1 system property, found`,
			len(sTree.Child.Children[buckId].(*Bucket).PropertySystem),
		)
	}

	for _, p := range sTree.Child.Children[buckId].(*Bucket).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalue` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*SomaTreeElemGroup).PropertySystem) != 1 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*SomaTreeElemCluster).PropertySystem) != 1 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*SomaTreeElemNode).PropertySystem) != 1 {
		t.Error(
			`Node has wrong system property count`,
		)
	}
}

func TestUpdateProperty(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`

	propUUID, _ := uuid.FromString(propId)

	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalue`,
	})

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).UpdateProperty(&PropertySystem{
		SourceId:     propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalueUPDATED`,
	})
	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}

	if len(actionChan) != 19 {
		t.Error(
			`Expected 19 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`repository`, `property_new`},
		[]string{`bucket`, `create`},
		[]string{`bucket`, `property_new`},
		[]string{`group`, `create`},
		[]string{`group`, `property_new`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `property_new`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`node`, `property_new`},
		[]string{`repository`, `property_update`},
		[]string{`bucket`, `property_update`},
		[]string{`group`, `property_update`},
		[]string{`cluster`, `property_update`},
		[]string{`node`, `property_update`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	if len(sTree.Child.PropertySystem) != 1 {
		t.Error(
			`Exptected repository to have 1 system property, found`,
			len(sTree.Child.PropertySystem),
		)
	}

	if _, ok := sTree.Child.PropertySystem[propId]; !ok {
		t.Error(
			`Could not find property under requested id`,
		)
	}

	if len(sTree.Child.Children[buckId].(*Bucket).PropertySystem) != 1 {
		t.Error(
			`Expected bucket to have 1 system property, found`,
			len(sTree.Child.Children[buckId].(*Bucket).PropertySystem),
		)
	}

	for _, p := range sTree.Child.Children[buckId].(*Bucket).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalueUPDATED` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*SomaTreeElemGroup).PropertySystem) != 1 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*SomaTreeElemCluster).PropertySystem) != 1 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*SomaTreeElemNode).PropertySystem) != 1 {
		t.Error(
			`Node has wrong system property count`,
		)
	}
}

func TestDeleteProperty(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`

	propUUID, _ := uuid.FromString(propId)

	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalue`,
	})

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).UpdateProperty(&PropertySystem{
		SourceId:     propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalueUPDATED`,
	})

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).DeleteProperty(&PropertySystem{
		SourceId: propUUID,
	})

	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}

	if len(actionChan) != 24 {
		t.Error(
			`Expected 24 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`repository`, `property_new`},
		[]string{`bucket`, `create`},
		[]string{`bucket`, `property_new`},
		[]string{`group`, `create`},
		[]string{`group`, `property_new`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `property_new`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`node`, `property_new`},
		[]string{`repository`, `property_update`},
		[]string{`bucket`, `property_update`},
		[]string{`group`, `property_update`},
		[]string{`cluster`, `property_update`},
		[]string{`node`, `property_update`},
		[]string{`repository`, `property_delete`},
		[]string{`bucket`, `property_delete`},
		[]string{`group`, `property_delete`},
		[]string{`cluster`, `property_delete`},
		[]string{`node`, `property_delete`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	if len(sTree.Child.PropertySystem) != 0 {
		t.Error(
			`Exptected repository to have 0 system property, found`,
			len(sTree.Child.PropertySystem),
		)
	}

	if len(sTree.Child.Children[buckId].(*Bucket).PropertySystem) != 0 {
		t.Error(
			`Expected bucket to have 0 system property, found`,
			len(sTree.Child.Children[buckId].(*Bucket).PropertySystem),
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*SomaTreeElemGroup).PropertySystem) != 0 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*SomaTreeElemCluster).PropertySystem) != 0 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*SomaTreeElemNode).PropertySystem) != 0 {
		t.Error(
			`Node has wrong system property count`,
		)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
