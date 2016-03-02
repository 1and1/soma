package somatree

import (
	"crypto/sha512"
	"encoding/base64"
	"io"
	"sort"

	"github.com/satori/go.uuid"
)

type Checker interface {
	SetCheck(c Check)

	inheritCheck(c Check)
	inheritCheckDeep(c Check)
	storeCheck(c Check)
	syncCheck(childId string)
	checkCheck(checkId string) bool
}

type Check struct {
	Id            uuid.UUID
	SourceId      uuid.UUID
	SourceType    string
	Inherited     bool
	InheritedFrom uuid.UUID
	CapabilityId  uuid.UUID
	ConfigId      uuid.UUID
	Inheritance   bool
	ChildrenOnly  bool
	View          string
	Interval      uint64
	Thresholds    []CheckThreshold
	Constraints   []CheckConstraint
}

type CheckThreshold struct {
	Predicate string
	Level     uint8
	Value     int64
}

type CheckConstraint struct {
	Type  string
	Key   string
	Value string
}

type CheckInstance struct {
	InstanceId            uuid.UUID
	CheckId               uuid.UUID
	ConfigId              uuid.UUID
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

func (tc Check) Clone() Check {
	cl := Check{
		SourceType:   tc.SourceType,
		Inherited:    tc.Inherited,
		Inheritance:  tc.Inheritance,
		ChildrenOnly: tc.ChildrenOnly,
		View:         tc.View,
		Interval:     tc.Interval,
	}
	cl.Id, _ = uuid.FromString(tc.Id.String())
	cl.SourceId, _ = uuid.FromString(tc.SourceId.String())
	cl.InheritedFrom, _ = uuid.FromString(tc.InheritedFrom.String())
	cl.CapabilityId, _ = uuid.FromString(tc.CapabilityId.String())
	cl.ConfigId, _ = uuid.FromString(tc.ConfigId.String())
	cl.Thresholds = make([]CheckThreshold, 0)
	for _, thr := range tc.Thresholds {
		t := CheckThreshold{
			Predicate: thr.Predicate,
			Level:     thr.Level,
			Value:     thr.Value,
		}
		cl.Thresholds = append(cl.Thresholds, t)
	}
	cl.Constraints = make([]CheckConstraint, 0)
	for _, cstr := range tc.Constraints {
		c := CheckConstraint{
			Type:  cstr.Type,
			Key:   cstr.Key,
			Value: cstr.Value,
		}
		cl.Constraints = append(cl.Constraints, c)
	}

	return cl
}

func (tci *CheckInstance) Clone() CheckInstance {
	cl := CheckInstance{
		Version:            tci.Version,
		ConstraintHash:     tci.ConstraintHash,
		ConstraintValHash:  tci.ConstraintValHash,
		ConstraintOncall:   tci.ConstraintOncall,
		InstanceSvcCfgHash: tci.InstanceSvcCfgHash,
	}
	cl.InstanceId, _ = uuid.FromString(tci.InstanceId.String())
	cl.CheckId, _ = uuid.FromString(tci.CheckId.String())
	cl.ConfigId, _ = uuid.FromString(tci.ConfigId.String())
	cl.InstanceService, _ = uuid.FromString(tci.InstanceService.String())
	for k, v := range tci.ConstraintService {
		t := v
		cl.ConstraintService[k] = t
	}
	for k, v := range tci.ConstraintSystem {
		t := v
		cl.ConstraintSystem[k] = t
	}
	for k, v := range tci.ConstraintCustom {
		t := v
		cl.ConstraintCustom[k] = t
	}
	for k, v := range tci.ConstraintNative {
		t := v
		cl.ConstraintNative[k] = t
	}
	for k, v := range tci.InstanceServiceConfig {
		t := v
		cl.InstanceServiceConfig[k] = t
	}
	cl.ConstraintAttribute = make(map[string]map[string][]string, 0)
	for k, _ := range tci.ConstraintAttribute {
		for k2, aVal := range tci.ConstraintAttribute[k] {
			for _, val := range aVal {
				t := val
				cl.ConstraintAttribute[k][k2] = append(cl.ConstraintAttribute[k][k2], t)
			}
		}
	}

	return cl
}

func (tci *CheckInstance) calcConstraintHash() {
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

func (tci *CheckInstance) calcConstraintValHash() {
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

func (tci *CheckInstance) calcInstanceSvcCfgHash() {
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
