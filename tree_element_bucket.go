package somatree

import "github.com/satori/go.uuid"

type SomaTreeElemBucket struct {
	Id              uuid.UUID
	Name            string
	Environment     string
	Type            string
	State           string
	Parent          SomaTreeBucketReceiver `json:"-"`
	Fault           *SomaTreeElemFault     `json:"-"`
	PropertyOncall  map[string]SomaTreeProperty
	PropertyService map[string]SomaTreeProperty
	PropertySystem  map[string]SomaTreeProperty
	PropertyCustom  map[string]SomaTreeProperty
	Checks          map[string]SomaTreeCheck
	Children        map[string]SomaTreeBucketAttacher //`json:"-"`
	Action          chan *Action                      `json:"-"`
}

//
// NEW
func NewBucket(name string, environment string, id string) *SomaTreeElemBucket {
	teb := new(SomaTreeElemBucket)
	if id == "" {
		teb.Id = uuid.NewV4()
	} else {
		teb.Id, _ = uuid.FromString(id)
	}
	teb.Name = name
	teb.Environment = environment
	teb.Type = "bucket"
	teb.State = "floating"
	teb.Parent = nil
	teb.Children = make(map[string]SomaTreeBucketAttacher)
	teb.PropertyOncall = make(map[string]SomaTreeProperty)
	teb.PropertyService = make(map[string]SomaTreeProperty)
	teb.PropertySystem = make(map[string]SomaTreeProperty)
	teb.PropertyCustom = make(map[string]SomaTreeProperty)
	teb.Checks = make(map[string]SomaTreeCheck)

	return teb
}

func (teb SomaTreeElemBucket) CloneRepository() SomaTreeRepositoryAttacher {
	cl := SomaTreeElemBucket{
		Name:        teb.Name,
		Environment: teb.Environment,
		Type:        teb.Type,
		State:       teb.State,
	}
	cl.Id, _ = uuid.FromString(teb.Id.String())
	f := make(map[string]SomaTreeBucketAttacher)
	for k, child := range teb.Children {
		f[k] = child.CloneBucket()
	}
	cl.Children = f

	pO := make(map[string]SomaTreeProperty)
	for k, prop := range teb.PropertyOncall {
		pO[k] = prop.Clone()
	}
	cl.PropertyOncall = pO

	pSv := make(map[string]SomaTreeProperty)
	for k, prop := range teb.PropertyService {
		pSv[k] = prop.Clone()
	}
	cl.PropertyService = pSv

	pSy := make(map[string]SomaTreeProperty)
	for k, prop := range teb.PropertySystem {
		pSy[k] = prop.Clone()
	}
	cl.PropertySystem = pSy

	pC := make(map[string]SomaTreeProperty)
	for k, prop := range teb.PropertyCustom {
		pC[k] = prop.Clone()
	}
	cl.PropertyCustom = pC

	cK := make(map[string]SomaTreeCheck)
	for k, chk := range teb.Checks {
		cK[k] = chk.Clone()
	}
	cl.Checks = cK

	return &cl
}

//
// Interface: SomaTreeBuilder
func (teb *SomaTreeElemBucket) GetID() string {
	return teb.Id.String()
}

func (teb *SomaTreeElemBucket) GetName() string {
	return teb.Name
}

func (teb *SomaTreeElemBucket) GetType() string {
	return teb.Type
}

func (teb *SomaTreeElemBucket) setAction(c chan *Action) {
	teb.Action = c
}

func (teb *SomaTreeElemBucket) setActionDeep(c chan *Action) {
	teb.setAction(c)
	for ch, _ := range teb.Children {
		teb.Children[ch].setActionDeep(c)
	}
}

//
// Interface: SomaTreeBucketeer
func (teb *SomaTreeElemBucket) GetBucket() SomaTreeReceiver {
	return teb
}

func (teb *SomaTreeElemBucket) GetEnvironment() string {
	return teb.Environment
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
