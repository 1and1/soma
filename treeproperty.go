package somaproto

type TreeProperty struct {
	InstanceId       string
	SourceInstanceId string
	SourceType       string
	IsInherited      bool
	InheritedFrom    string
	Inheritance      bool
	ChildrenOnly     bool
	View             string
	PropertyType     string
	RepositoryId     string
	BucketId         string
	Custom           *TreePropertyCustom
	System           *TreePropertySystem
	Service          *TreePropertyService
	Native           *TreePropertyNative
	Oncall           *TreePropertyOncall
}

type TreePropertyCustom struct {
	CustomId     string
	RepositoryId string
	Name         string
	Value        string
}

type TreePropertySystem struct {
	Name  string
	Value string
}

type TreePropertyService struct {
	Name       string
	TeamId     string
	Attributes []TreeServiceAttribute
}

type TreePropertyNative struct {
	Name  string
	Value string
}

type TreePropertyOncall struct {
	OncallId string
	Name     string
	Number   string
}

type TreeServiceAttribute struct {
	Attribute string
	Value     string
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
