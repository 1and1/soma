package somatree

import "github.com/satori/go.uuid"

type SomaTreeChecker interface {
	SetCheck(c SomaTreeCheck)

	inheritCheck(c SomaTreeCheck)
	inheritCheckDeep(c SomaTreeCheck)
	storeCheck(c SomaTreeCheck)
	syncCheck(childId string)
	checkCheck(checkId string) bool
}

type SomaTreeCheck struct {
	Id            uuid.UUID
	Inherited     bool
	InheritedFrom uuid.UUID
	Inheritance   bool
	ChildrenOnly  bool
	CapabilityId  uuid.UUID
	Interval      uint64
	Thresholds    []SomaTreeCheckThreshold
	Constraints   []SomaTreeCheckConstraint
}

type SomaTreeCheckThreshold struct {
	Predicate string
	Level     uint8
	Value     int64
}

type SomaTreeCheckConstraint struct {
	Type  string
	Key   string
	Value string
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
