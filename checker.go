package somatree

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
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

type CheckGetter interface {
	GetCheckId() string
	GetSourceCheckId() string
	GetCheckConfigId() string
	GetSourceType() string
	GetIsInherited() bool
	GetInheritedFrom() string
	GetInheritance() bool
	GetChildrenOnly() bool
	GetView() string
	GetCapabilityId() string
	GetInterval() uint64
	GetItemId(objType string, objId uuid.UUID) uuid.UUID
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
	Items         []CheckItem
}

type CheckItem struct {
	ObjectId   uuid.UUID
	ObjectType string
	ItemId     uuid.UUID
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
	InstanceConfigId      uuid.UUID
	Version               uint64
	ConstraintHash        string
	ConstraintValHash     string
	ConstraintOncall      string                         // Ids
	ConstraintService     map[string]string              // svcName->value
	ConstraintSystem      map[string]string              // Id->value
	ConstraintCustom      map[string]string              // Id->value
	ConstraintNative      map[string]string              // prop->value
	ConstraintAttribute   map[string]map[string][]string // svcId->attr->[ value, value, ... ]
	InstanceServiceConfig map[string]string              // attr->value
	InstanceService       string
	InstanceSvcCfgHash    string
}

func (c *Check) GetItemId(objType string, objId uuid.UUID) uuid.UUID {
	if !uuid.Equal(c.Id, uuid.Nil) {
		return c.Id
	}
	for _, item := range c.Items {
		if objType == item.ObjectType && uuid.Equal(item.ObjectId, objId) {
			return item.ItemId
		}
	}
	return uuid.Nil
}

func (c *Check) GetCheckId() string {
	return c.Id.String()
}

func (c *Check) GetSourceCheckId() string {
	return c.SourceId.String()
}

func (c *Check) GetCheckConfigId() string {
	return c.ConfigId.String()
}

func (c *Check) GetSourceType() string {
	return c.SourceType
}

func (c *Check) GetIsInherited() bool {
	return c.Inherited
}

func (c *Check) GetInheritedFrom() string {
	return c.InheritedFrom.String()
}

func (c *Check) GetInheritance() bool {
	return c.Inheritance
}

func (c *Check) GetChildrenOnly() bool {
	return c.ChildrenOnly
}

func (c *Check) GetView() string {
	return c.View
}

func (c *Check) GetCapabilityId() string {
	return c.CapabilityId.String()
}

func (c *Check) GetInterval() uint64 {
	return c.Interval
}

func (c *Check) MakeAction() Action {
	return Action{
		Check: somaproto.TreeCheck{
			CheckId:       c.GetCheckId(),
			SourceCheckId: c.GetSourceCheckId(),
			CheckConfigId: c.GetCheckConfigId(),
			SourceType:    c.GetSourceType(),
			IsInherited:   c.GetIsInherited(),
			InheritedFrom: c.GetInheritedFrom(),
			Inheritance:   c.GetInheritance(),
			ChildrenOnly:  c.GetChildrenOnly(),
			CapabilityId:  c.GetCapabilityId(),
		},
	}
}

func (tci *CheckInstance) Clone() CheckInstance {
	cl := CheckInstance{
		Version:            tci.Version,
		ConstraintHash:     tci.ConstraintHash,
		ConstraintValHash:  tci.ConstraintValHash,
		ConstraintOncall:   tci.ConstraintOncall,
		InstanceSvcCfgHash: tci.InstanceSvcCfgHash,
		InstanceService:    tci.InstanceService,
	}
	cl.InstanceConfigId, _ = uuid.FromString(tci.InstanceConfigId.String())
	cl.InstanceId, _ = uuid.FromString(tci.InstanceId.String())
	cl.CheckId, _ = uuid.FromString(tci.CheckId.String())
	cl.ConfigId, _ = uuid.FromString(tci.ConfigId.String())
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

func (c *Check) clone() Check {
	cl := Check{
		SourceType:   c.SourceType,
		Inherited:    c.Inherited,
		Inheritance:  c.Inheritance,
		ChildrenOnly: c.ChildrenOnly,
		View:         c.View,
		Interval:     c.Interval,
	}
	cl.Id, _ = uuid.FromString(c.Id.String())
	cl.SourceId, _ = uuid.FromString(c.SourceId.String())
	cl.InheritedFrom, _ = uuid.FromString(c.InheritedFrom.String())
	cl.CapabilityId, _ = uuid.FromString(c.CapabilityId.String())
	cl.ConfigId, _ = uuid.FromString(c.ConfigId.String())
	cl.Thresholds = make([]CheckThreshold, len(c.Thresholds))
	for i, thr := range c.Thresholds {
		n := CheckThreshold{
			Predicate: thr.Predicate,
			Level:     thr.Level,
			Value:     thr.Value,
		}
		cl.Thresholds[i] = n
	}
	cl.Constraints = make([]CheckConstraint, len(c.Constraints))
	for i, ctr := range c.Constraints {
		n := CheckConstraint{
			Type:  ctr.Type,
			Key:   ctr.Key,
			Value: ctr.Value,
		}
		cl.Constraints[i] = n
	}
	cl.Items = make([]CheckItem, len(c.Items))
	for i, item := range c.Items {
		n := CheckItem{
			ObjectType: item.ObjectType,
		}
		n.ItemId, _ = uuid.FromString(item.ItemId.String())
		n.ObjectId, _ = uuid.FromString(item.ObjectId.String())
		cl.Items[i] = n
	}
	return cl
}

func (ci CheckInstance) MakeAction() Action {
	serviceCfg, err := json.Marshal(ci.InstanceServiceConfig)
	if err != nil {
		serviceCfg = []byte{}
	}

	return Action{
		CheckInstance: somaproto.TreeCheckInstance{
			InstanceId:            ci.InstanceId.String(),
			CheckId:               ci.CheckId.String(),
			ConfigId:              ci.ConfigId.String(),
			InstanceConfigId:      ci.InstanceConfigId.String(),
			Version:               ci.Version,
			ConstraintHash:        ci.ConstraintHash,
			ConstraintValHash:     ci.ConstraintValHash,
			InstanceSvcCfgHash:    ci.InstanceSvcCfgHash,
			InstanceService:       ci.InstanceService,
			InstanceServiceConfig: string(serviceCfg),
		},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
