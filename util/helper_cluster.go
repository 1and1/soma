package util

import (
	"bytes"
	"encoding/json"
	"fmt"

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
	var req somaproto.ProtoRequestCluster
	req.Filter.Name = c
	req.Filter.BucketId = bId

	resp := u.GetRequestWithBody(req, "/clusters/")
	clusterResult := u.DecodeProtoResultClusterFromResponse(resp)

	return clusterResult.Clusters[0].Id
}

func (u SomaUtil) DecodeProtoResultClusterFromResponse(resp *resty.Response) *somaproto.ProtoResultCluster {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body))
	var res somaproto.ProtoResultCluster
	err := decoder.Decode(&res)
	u.AbortOnError(err, "Error decoding server response body")
	if res.Code > 299 {
		s := fmt.Sprintf("Request failed: %d - %s", res.Code, res.Status)
		msgs := []string{s}
		msgs = append(msgs, res.Text...)
		u.Abort(msgs...)
	}
	return &res
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
