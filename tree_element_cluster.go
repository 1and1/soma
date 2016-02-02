package somatree

import "github.com/satori/go.uuid"

type SomaTreeElemCluster struct {
	Id              uuid.UUID
	Name            string
	Parent          SomaTreeClusterReceiver
	Children        map[string]SomaTreeClusterAttacher
	PropertyOncall  map[string]*SomaTreePropertyOncall
	PropertyService map[string]*SomaTreePropertyService
	PropertySystem  map[string]*SomaTreePropertySystem
	PropertyCustom  map[string]*SomaTreePropertyCustom
	Checks          map[string]*SomaTreeCheck
}

func NewCluster() *SomaTreeElemCluster {
	return new(SomaTreeElemCluster)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
