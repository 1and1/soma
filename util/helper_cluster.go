package util

import (
	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetClusterByUUIDOrName(c string, b string) string {
	var id string
	bId := u.BucketByUUIDOrName(b)

	cId, err := uuid.FromString(c)
	if err != nil {
		id = u.GetClusterIdByName(c, bId)
	} else {
		id = cId.String()
	}
	return id
}

func (u SomaUtil) GetClusterIdByName(c string, bId string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Cluster: &proto.ClusterFilter{
				Name:     c,
				BucketId: bId,
			},
		},
	}

	resp := u.PostRequestWithBody(req, "/filter/clusters/")
	clusterResult := u.DecodeProtoResultClusterFromResponse(resp)

	return (*clusterResult.Clusters)[0].Id
}

func (u SomaUtil) DecodeProtoResultClusterFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
