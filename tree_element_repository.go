package somatree

import (
	"log"

	"github.com/satori/go.uuid"
)

type SomaTreeElemRepository struct {
	Id      uuid.UUID
	Name    string
	Team    uuid.UUID
	Deleted bool
	Active  bool
	Type    string
	Parent  SomaTreeRepositoryReceiver `json:"-"`
	//Fault    SomaTreeAttacher `json:"-"`
	Children map[string]SomaTreeRepositoryAttacher
	//PropertyOncall  map[string]*SomaTreePropertyOncall
	//PropertyService map[string]*SomaTreePropertyService
	//PropertySystem  map[string]*SomaTreePropertySystem
	//PropertyCustom  map[string]*SomaTreePropertyCustom
	//Checks          map[string]*SomaTreeCheck
}

type RepositorySpec struct {
	Id      uuid.UUID
	Name    string
	Team    uuid.UUID
	Deleted bool
	Active  bool
}

//
// NEW
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

//
// Interface: SomaTreeBuilder
func (ter *SomaTreeElemRepository) GetID() string {
	return ter.Id.String()
}

func (ter *SomaTreeElemRepository) GetName() string {
	return ter.Name
}

func (ter *SomaTreeElemRepository) GetType() string {
	return ter.Type
}

//
// Interface: SomaTreeAttacher
func (ter *SomaTreeElemRepository) Attach(a AttachRequest) {
	switch {
	case a.ParentType == "root" &&
		a.ChildType == "repository" &&
		a.ChildName == ter.Name:
		ter.attachToRoot(a)
	}
}

func (ter *SomaTreeElemRepository) ReAttach(a AttachRequest) {
	log.Fatal("Not implemented")
}

func (ter *SomaTreeElemRepository) setParent(p SomaTreeRepositoryReceiver) {
	ter.Parent = p
}

func (ter *SomaTreeElemRepository) Destroy() {
	ter.Parent.Unlink(UnlinkRequest{
		ParentType: ter.Parent.(SomaTreeBuilder).GetType(),
		ParentId:   ter.Parent.(SomaTreeBuilder).GetID(),
		ParentName: ter.Parent.(SomaTreeBuilder).GetName(),
		ChildType:  ter.GetType(),
		ChildName:  ter.GetName(),
		ChildId:    ter.GetID(),
	},
	)
}

// Interface: SomaTreeRootAttacher
func (ter *SomaTreeElemRepository) attachToRoot(a AttachRequest) {
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
	if receiveRequestCheck(r, ter) {
		switch r.ChildType {
		case "bucket":
			ter.receiveBucket(r)
		default:
			panic(`SomaTreeElemRepository.Receive`)
		}
		return
	}
	for _, child := range ter.Children {
		child.(SomaTreeReceiver).Receive(r)
	}
}

// Interface: SomaTreeUnlinker
func (ter *SomaTreeElemRepository) Unlink(u UnlinkRequest) {
	if unlinkRequestCheck(u, ter) {
		switch u.ChildType {
		case "bucket":
			ter.unlinkBucket(u)
		default:
			panic(`SomaTreeElemRepository.Unlink`)
		}
		return
	}
	for _, child := range ter.Children {
		child.(SomaTreeUnlinker).Unlink(u)
	}
}

// Interface: SomaTreeBucketReceiver
func (ter *SomaTreeElemRepository) receiveBucket(r ReceiveRequest) {
	if receiveRequestCheck(r, ter) {
		switch r.ChildType {
		case "bucket":
			ter.Children[r.Bucket.GetID()] = r.Bucket
			r.Bucket.setParent(ter)
		default:
			panic(`SomaTreeElemRepository.receiveBucket`)
		}
	}
}

// Interface: SomaTreeBucketUnlinker
func (ter *SomaTreeElemRepository) unlinkBucket(u UnlinkRequest) {
	if unlinkRequestCheck(u, ter) {
		switch u.ChildType {
		case "bucket":
			if _, ok := ter.Children[u.ChildId]; ok {
				if u.ChildName == ter.Children[u.ChildId].GetName() {
					delete(ter.Children, u.ChildId)
				}
			}
		default:
			panic(`SomaTreeElemRepository.unlinkBucket`)
		}
		return
	}
	panic(`SomaTreeElemRepository.unlinkBucket`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
