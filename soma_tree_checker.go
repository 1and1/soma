package somatree

import (
	"crypto/sha512"
	"encoding/base64"
	"io"
	"sort"

	"github.com/satori/go.uuid"
)

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
	View          string
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

type SomaTreeCheckInstance struct {
	InstanceId            uuid.UUID
	CheckId               uuid.UUID
	Version               uint32
	ConstraintHash        string
	ConstraintValHash     string
	ConstraintOncall      string                         // Ids
	ConstraintService     map[string]string              // Id->value
	ConstraintSystem      map[string]string              // Id->value
	ConstraintCustom      map[string]string              // Id->value
	ConstraintNative      map[string]string              // prop->value
	ConstraintAttribute   map[string]map[string][]string // svcId->attr->[ value, value, ... ]
	InstanceServiceConfig map[string]string              // attr->value
	InstanceService       uuid.UUID
	InstanceSvcCfgHash    string
}

func (tci *SomaTreeCheckInstance) calcConstraintHash() {
	h := sha512.New()
	io.WriteString(h, tci.ConstraintOncall)

	services := []string{}
	for i, _ := range tci.ConstraintService {
		j := i
		services = append(services, j)
	}
	sort.Strings(services)
	for _, i := range services {
		io.WriteString(h, i)
	}

	systems := []string{}
	for i, _ := range tci.ConstraintSystem {
		j := i
		systems = append(systems, j)
	}
	sort.Strings(systems)
	for _, i := range systems {
		io.WriteString(h, i)
	}

	customs := []string{}
	for i, _ := range tci.ConstraintCustom {
		j := i
		customs = append(customs, j)
	}
	sort.Strings(customs)
	for _, i := range customs {
		io.WriteString(h, i)
	}

	natives := []string{}
	for i, _ := range tci.ConstraintNative {
		j := i
		natives = append(natives, j)
	}
	sort.Strings(natives)
	for _, i := range natives {
		io.WriteString(h, i)
	}

	attributes := []string{}
	for i, _ := range tci.ConstraintAttribute {
		j := i
		attributes = append(attributes, j)
	}
	sort.Strings(attributes)
	for _, i := range attributes {
		svcattr := []string{}
		for j, _ := range tci.ConstraintAttribute[i] {
			k := j
			svcattr = append(svcattr, k)
		}
		sort.Strings(svcattr)
		io.WriteString(h, i)
		for _, l := range svcattr {
			io.WriteString(h, l)
		}
	}
	tci.ConstraintHash = base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (tci *SomaTreeCheckInstance) calcConstraintValHash() {
	h := sha512.New()
	io.WriteString(h, tci.ConstraintOncall)

	services := []string{}
	for i, _ := range tci.ConstraintService {
		j := i
		services = append(services, j)
	}
	sort.Strings(services)
	for _, i := range services {
		io.WriteString(h, i)
		io.WriteString(h, tci.ConstraintService[i])
	}

	systems := []string{}
	for i, _ := range tci.ConstraintSystem {
		j := i
		systems = append(systems, j)
	}
	sort.Strings(systems)
	for _, i := range systems {
		io.WriteString(h, i)
		io.WriteString(h, tci.ConstraintSystem[i])
	}

	customs := []string{}
	for i, _ := range tci.ConstraintCustom {
		j := i
		customs = append(customs, j)
	}
	sort.Strings(customs)
	for _, i := range customs {
		io.WriteString(h, i)
		io.WriteString(h, tci.ConstraintCustom[i])
	}

	natives := []string{}
	for i, _ := range tci.ConstraintNative {
		j := i
		natives = append(natives, j)
	}
	sort.Strings(natives)
	for _, i := range natives {
		io.WriteString(h, i)
		io.WriteString(h, tci.ConstraintNative[i])
	}

	attributes := []string{}
	for i, _ := range tci.ConstraintAttribute {
		j := i
		attributes = append(attributes, j)
	}
	sort.Strings(attributes)
	for _, i := range attributes {
		svcattr := []string{}
		for j, _ := range tci.ConstraintAttribute[i] {
			k := j
			svcattr = append(svcattr, k)
		}
		sort.Strings(svcattr)
		io.WriteString(h, i)
		for _, l := range svcattr {
			io.WriteString(h, l)
			vals := make([]string, len(tci.ConstraintAttribute[i][l]))
			copy(vals, tci.ConstraintAttribute[i][l])
			sort.Strings(vals)
			for _, m := range vals {
				io.WriteString(h, m)
			}
		}
	}
	tci.ConstraintValHash = base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (tci *SomaTreeCheckInstance) calcInstanceSvcCfgHash() {
	h := sha512.New()

	attributes := []string{}
	for i, _ := range tci.InstanceServiceConfig {
		j := i
		attributes = append(attributes, j)
	}
	sort.Strings(attributes)
	for _, i := range attributes {
		io.WriteString(h, i)
		io.WriteString(h, tci.InstanceServiceConfig[i])
	}
	tci.InstanceSvcCfgHash = base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
