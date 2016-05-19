package util

import (
	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetClusterByUUIDOrName(c *resty.Client, cl string, b string) string {
	var id string
	bId := u.BucketByUUIDOrName(c, b)

	cId, err := uuid.FromString(cl)
	if err != nil {
		id = u.GetClusterIdByName(c, cl, bId)
	} else {
		id = cId.String()
	}
	return id
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
