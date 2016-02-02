package somatree

import (
	"fmt"

	"github.com/satori/go.uuid"
)

type SomaTreeElemRepository struct {
	Id       uuid.UUID
	Name     string
	Type     string
	Team     uuid.UUID
	Parent   SomaTreeRepositoryReceiver `json:"-"`
	Children map[string]SomaTreeRepositoryAttacher
	//PropertyOncall  map[string]*SomaTreePropertyOncall
	//PropertyService map[string]*SomaTreePropertyService
	//PropertySystem  map[string]*SomaTreePropertySystem
	//PropertyCustom  map[string]*SomaTreePropertyCustom
	//Checks          map[string]*SomaTreeCheck
}

func NewRepository(name string) *SomaTreeElemRepository {
	ter := new(SomaTreeElemRepository)
	ter.Id = uuid.NewV4()
	ter.Name = name
	ter.Type = "repository"
	ter.Children = make(map[string]SomaTreeRepositoryAttacher)
	//ter.PropertyOncall = make(map[string]*SomaTreePropertyOncall)
	//ter.PropertyService = make(map[string]*SomaTreePropertyService)
	//ter.PropertySystem = make(map[string]*SomaTreePropertySystem)
	//ter.PropertyCustom = make(map[string]*SomaTreePropertyCustom)
	//ter.Checks = make(map[string]*SomaTreeCheck)

	return ter
}

func (ter *SomaTreeElemRepository) GetID() string {
	return ter.Id.String()
}

func (ter *SomaTreeElemRepository) GetName() string {
	return ter.Name
}

// Interface: SomaTreeAttacher
func (ter *SomaTreeElemRepository) Attach(a AttachRequest) {
	switch {
	case a.ParentType == "root" &&
		a.ChildType == "repository" &&
		a.ChildName == ter.Name:
		ter.AttachToRoot(a)
	}
}

func (ter *SomaTreeElemRepository) ReAttach(a AttachRequest) {
	fmt.Println("bl")
}

func (ter *SomaTreeElemRepository) SetParent(p SomaTreeRepositoryReceiver) {
	ter.Parent = p
}

// Interface: SomaTreeRootAttacher
func (ter *SomaTreeElemRepository) AttachToRoot(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  "repository",
		Repository: ter,
	})
}

// Interface: SomaTreeReceiver
func (ter *SomaTreeElemRepository) Receive(r ReceiveRequest) {
	switch {
	case r.ParentType == ter.Type &&
		(r.ParentId == ter.Id.String() ||
			r.ParentName == ter.Name) &&
		r.ChildType == "bucket":
		ter.ReceiveBucket(r)
	default:
		for _, child := range ter.Children {
			child.(SomaTreeReceiver).Receive(r)
		}
	}
}

// Interface: SomaTreeUnlinker
func (ter *SomaTreeElemRepository) Unlink(u UnlinkRequest) {
	switch {
	case u.ParentType == ter.Type &&
		(u.ParentId == ter.Id.String() ||
			u.ParentName == ter.Name) &&
		u.ChildType == "bucket":
		ter.UnlinkBucket(u)
	default:
		for _, child := range ter.Children {
			child.(SomaTreeUnlinker).Unlink(u)
		}
	}
}

// Interface: SomaTreeBucketReceiver
func (ter *SomaTreeElemRepository) ReceiveBucket(r ReceiveRequest) {
	switch {
	case r.ParentType == ter.Type &&
		(r.ParentId == ter.Id.String() ||
			r.ParentName == ter.Name) &&
		r.ChildType == "bucket":
		ter.Children[r.Bucket.GetID()] = r.Bucket
		r.Bucket.SetParent(ter)
	default:
		panic("not allowed")
	}
}

// Interface: SomaTreeBucketUnlinker
func (ter *SomaTreeElemRepository) UnlinkBucket(u UnlinkRequest) {
	switch {
	case u.ParentType == ter.Type &&
		(u.ParentId == ter.Id.String() ||
			u.ParentName == ter.Name):
		if _, ok := ter.Children[u.ChildId]; ok {
			if u.ChildName == ter.Children[u.ChildId].GetName() {
				delete(ter.Children, u.ChildId)
			}
		}
	default:
		panic("not allowed")
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
