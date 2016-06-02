package util

import (
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetClusterByUUIDOrName(c *resty.Client, cl string, b string) string {
	if u.IsUUID(cl) {
		return cl
	}
	bId := u.BucketByUUIDOrName(c, b)
	return u.GetClusterIdByName(c, cl, bId)
}

func (u SomaUtil) GetClusterIdByName(c *resty.Client, cl string, bId string) string {
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

func (u SomaUtil) DecodeProtoResultClusterFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
