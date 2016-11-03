package util

import (
	"fmt"

	"github.com/1and1/soma/lib/proto"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetClusterByUUIDOrName(c *resty.Client, cl string, b string) string {
	if u.isUUID(cl) {
		return cl
	}
	bId := u.bucketByUUIDOrName(c, b)
	return u.getClusterIdByName(c, cl, bId)
}

func (u SomaUtil) getClusterIdByName(c *resty.Client, cl string, bId string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Cluster: &proto.ClusterFilter{
				Name:     cl,
				BucketId: bId,
			},
		},
	}

	resp := u.PostRequestWithBody(c, req, "/filter/clusters/")
	clusterResult := u.DecodeProtoResultClusterFromResponse(resp)

	return (*clusterResult.Clusters)[0].Id
}

func (u SomaUtil) getClusterDetails(c *resty.Client, clusterId string) *proto.Cluster {
	resp := u.GetRequest(c, fmt.Sprintf("/clusters/%s", clusterId))
	res := u.DecodeResultFromResponse(resp)
	return &(*res.Clusters)[0]
}

func (u SomaUtil) findSourceForClusterProperty(c *resty.Client, pTyp, pName, view, clusterId string) string {
	cluster := u.getClusterDetails(c, clusterId)
	if cluster == nil {
		return ``
	}
	for _, prop := range *cluster.Properties {
		// wrong type
		if prop.Type != pTyp {
			continue
		}
		// wrong view
		if prop.View != view {
			continue
		}
		// inherited property
		if prop.InstanceId != prop.SourceInstanceId {
			continue
		}
		switch pTyp {
		case `system`:
			if prop.System.Name == pName {
				return prop.SourceInstanceId
			}
		case `oncall`:
			if prop.Oncall.Name == pName {
				return prop.SourceInstanceId
			}
		case `custom`:
			if prop.Custom.Name == pName {
				return prop.SourceInstanceId
			}
		case `service`:
			if prop.Service.Name == pName {
				return prop.SourceInstanceId
			}
		}
	}
	return ``
}

func (u SomaUtil) DecodeProtoResultClusterFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
