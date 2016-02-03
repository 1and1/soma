package somatree

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/satori/go.uuid"
)

type SomaTreeElemRepository struct {
	Id       uuid.UUID
	Name     string
	Team     uuid.UUID
	Deleted  bool
	Active   bool
	Type     string
	State    string
	Parent   SomaTreeRepositoryReceiver `json:"-"`
	Fault    *SomaTreeElemFault         `json:"-"`
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
	ter.State = "floating"
	ter.Parent = nil
	ter.Children = make(map[string]SomaTreeRepositoryAttacher)
	//ter.PropertyOncall = make(map[string]*SomaTreePropertyOncall)
	//ter.PropertyService = make(map[string]*SomaTreePropertyService)
	//ter.PropertySystem = make(map[string]*SomaTreePropertySystem)
	//ter.PropertyCustom = make(map[string]*SomaTreePropertyCustom)
	//ter.Checks = make(map[string]*SomaTreeCheck)

	// return new repository with attached fault handler
	newFault().Attach(
		AttachRequest{
			Root:       ter,
			ParentType: ter.Type,
			ParentName: ter.Name,
		},
	)
	return ter
}

func (ter SomaTreeElemRepository) Clone() SomaTreeElemRepository {
	f := make(map[string]SomaTreeRepositoryAttacher)
	for k, child := range ter.Children {
		f[k] = child.CloneRepository()
	}
	ter.Children = f
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
	if ter.Parent != nil {
		panic(`SomaTreeElemRepository.Attach: already attached`)
	}
	switch {
	case a.ParentType == "root" &&
		a.ChildType == "repository" &&
		a.ChildName == ter.Name:
		ter.attachToRoot(a)
	}
}

func (ter *SomaTreeElemRepository) setParent(p SomaTreeReceiver) {
	switch p.(type) {
	case SomaTreeRepositoryReceiver:
		ter.setRepositoryParent(p.(SomaTreeRepositoryReceiver))
		ter.State = "attached"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`SomaTreeElemBucket.setParent`)
	}
}

func (ter *SomaTreeElemRepository) setRepositoryParent(p SomaTreeRepositoryReceiver) {
	ter.Parent = p
}

func (ter *SomaTreeElemRepository) updateParentRecursive(p SomaTreeReceiver) {
	ter.setParent(p.(SomaTreeRepositoryReceiver))
	var wg sync.WaitGroup
	for child, _ := range ter.Children {
		wg.Add(1)
		go func(str SomaTreeReceiver) {
			defer wg.Done()
			ter.Children[child].updateParentRecursive(str)
		}(ter)
	}
	wg.Wait()
}

func (ter *SomaTreeElemRepository) clearParent() {
	ter.Parent = nil
	ter.State = "floating"
}

func (ter *SomaTreeElemRepository) setFault(f *SomaTreeElemFault) {
	ter.Fault = f
}

func (ter *SomaTreeElemRepository) updateFaultRecursive(f *SomaTreeElemFault) {
	ter.setFault(f)
	var wg sync.WaitGroup
	for child, _ := range ter.Children {
		wg.Add(1)
		go func(ptr *SomaTreeElemFault) {
			defer wg.Done()
			ter.Children[child].updateFaultRecursive(ptr)
		}(f)
	}
	wg.Wait()
}

func (ter *SomaTreeElemRepository) Destroy() {
	if ter.Parent == nil {
		panic(`SomaTreeElemRepository.Destroy called without Parent to unlink from`)
	}

	// the Destroy handler of SomaTreeElemFault calls
	// updateFaultRecursive(nil) on us
	ter.Fault.Destroy()

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

func (ter *SomaTreeElemRepository) Detach() {
	ter.Destroy()
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
		case "fault":
			ter.receiveFault(r)
		default:
			panic(`SomaTreeElemRepository.Receive`)
		}
		return
	}
	for child, _ := range ter.Children {
		ter.Children[child].(SomaTreeReceiver).Receive(r)
	}
}

// Interface: SomaTreeUnlinker
func (ter *SomaTreeElemRepository) Unlink(u UnlinkRequest) {
	if unlinkRequestCheck(u, ter) {
		switch u.ChildType {
		case "bucket":
			ter.unlinkBucket(u)
		case "fault":
			ter.unlinkFault(u)
		default:
			panic(`SomaTreeElemRepository.Unlink`)
		}
		return
	}
	for child, _ := range ter.Children {
		ter.Children[child].(SomaTreeUnlinker).Unlink(u)
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
					ter.Children[u.ChildId].clearParent()
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

// Interface: SomaTreeFaultReceiver
func (ter *SomaTreeElemRepository) receiveFault(r ReceiveRequest) {
	if receiveRequestCheck(r, ter) {
		switch r.ChildType {
		case "fault":
			ter.setFault(r.Fault)
			ter.Fault.setParent(ter)
			ter.updateFaultRecursive(ter.Fault)
		default:
			panic(`SomaTreeElemRepository.receiveFault`)
		}
		return
	}
	panic(`SomaTreeElemRepository.receiveFault`)
}

// Interface: SomaTreeFaultUnlinker
func (ter *SomaTreeElemRepository) unlinkFault(u UnlinkRequest) {
	if unlinkRequestCheck(u, ter) {
		switch u.ChildType {
		case "fault":
			ter.Fault = nil
			ter.updateFaultRecursive(ter.Fault)
		default:
			panic(`SomaTreeElemRepository.unlinkFault`)
		}
		return
	}
	panic(`SomaTreeElemRepository.unlinkFault`)
}

//
// Interface: SomaTreeFinder
func (ter *SomaTreeElemRepository) Find(f FindRequest, b bool) SomaTreeAttacher {
	if findRequestCheck(f, ter) {
		return ter
	}
	var wg sync.WaitGroup
	rawResult := make(chan SomaTreeAttacher, len(ter.Children))
	for child, _ := range ter.Children {
		wg.Add(1)
		go func(fr FindRequest, bl bool) {
			defer wg.Done()
			rawResult <- ter.Children[child].(SomaTreeFinder).Find(fr, bl)
		}(f, false)
	}
	wg.Wait()
	close(rawResult)

	res := make([]SomaTreeAttacher, 0)
	for sta := range rawResult {
		if sta != nil {
			res = append(res, sta)
		}
	}
	switch {
	case len(res) == 0:
		if b {
			return ter.Fault
		} else {
			return nil
		}
	case len(res) > 1:
		return ter.Fault
	}
	return res[0]
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
