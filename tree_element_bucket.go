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
	f := make(map[string]SomaTreeBucketAttacher)
	for k, child := range teb.Children {
		f[k] = child.CloneBucket()
	}
	teb.Children = f
	return &teb
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

//
// Interface: SomaTreeBucketeer
func (teb *SomaTreeElemBucket) GetBucket() SomaTreeReceiver {
	return teb
}

func (teb *SomaTreeElemBucket) GetEnvironment() string {
	return teb.Environment
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
