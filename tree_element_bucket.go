package somatree

import (
	"fmt"

	"github.com/satori/go.uuid"
)

type SomaTreeElemBucket struct {
	Id          uuid.UUID
	Name        string
	Environment string
	Parent      SomaTreeBucketReceiver `json:"-"`
	Children    map[string]*SomaTreeBucketAttacher
	//PropertyOncall  map[string]*SomaTreePropertyOncall
	//PropertyService map[string]*SomaTreePropertyService
	//PropertySystem  map[string]*SomaTreePropertySystem
	//PropertyCustom  map[string]*SomaTreePropertyCustom
	//Checks          map[string]*SomaTreeCheck
}

func NewBucket(name string, environment string, id string) *SomaTreeElemBucket {
	teb := new(SomaTreeElemBucket)
	if id == "" {
		teb.Id = uuid.NewV4()
	} else {
		teb.Id, _ = uuid.FromString(id)
	}
	teb.Name = name
	teb.Environment = environment
	teb.Children = make(map[string]*SomaTreeBucketAttacher)
	//teb.PropertyOncall = make(map[string]*SomaTreePropertyOncall)
	//teb.PropertyService = make(map[string]*SomaTreePropertyService)
	//teb.PropertySystem = make(map[string]*SomaTreePropertySystem)
	//teb.PropertyCustom = make(map[string]*SomaTreePropertyCustom)
	//teb.Checks = make(map[string]*SomaTreeCheck)

	return teb
}

func (teb *SomaTreeElemBucket) GetID() string {
	return teb.Id.String()
}

func (teb *SomaTreeElemBucket) GetName() string {
	return teb.Name
}

func (teb *SomaTreeElemBucket) SetParent(p SomaTreeReceiver) {
	switch p.(type) {
	case SomaTreeBucketReceiver:
		teb.SetBucketParent(p.(SomaTreeBucketReceiver))
	default:
		panic("not allowed")
	}
}

func (teb *SomaTreeElemBucket) SetBucketParent(p SomaTreeBucketReceiver) {
	teb.Parent = p
}

// Interface: SomaTreeAttacher
func (teb *SomaTreeElemBucket) Attach(a AttachRequest) {
	switch {
	case a.ParentType == "repository" &&
		a.ChildType == "bucket" &&
		a.ChildName == teb.Name:
		teb.AttachToRepository(a)
	}
}

func (teb *SomaTreeElemBucket) ReAttach(a AttachRequest) {
	fmt.Println("bl")
}

// Interface: SomaTreeRepositoryAttacher
func (teb *SomaTreeElemBucket) AttachToRepository(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  "bucket",
		Bucket:     teb,
	})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
