package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetBucketByUUIDOrName(b string, r string) string {
	id, err := uuid.FromString(b)
	if err != nil {
		// aborts on failure
		return u.GetBucketIdByName(b, r)
	}
	return id.String()
}

func (u SomaUtil) BucketByUUIDOrName(b string) string {
	id, err := uuid.FromString(b)
	if err != nil {
		// aborts on failure
		return u.BucketIdByName(b)
	}
	return id.String()
}

func (u SomaUtil) GetBucketIdByName(bucket string, repoId string) string {
	var req somaproto.ProtoRequestBucket
	req.Filter.Name = bucket
	req.Filter.RepositoryId = repoId

	resp := u.GetRequestWithBody(req, "/buckets/")
	repoResult := u.DecodeProtoResultBucketFromResponse(resp)

	if bucket != repoResult.Buckets[0].Name {
		u.Abort("Received result set for incorrect bucket")
	}
	return repoResult.Buckets[0].Id
}

func (u SomaUtil) BucketIdByName(bucket string) string {
	var req somaproto.ProtoRequestBucket
	req.Filter.Name = bucket

	resp := u.GetRequestWithBody(req, "/buckets/")
	repoResult := u.DecodeProtoResultBucketFromResponse(resp)

	if bucket != repoResult.Buckets[0].Name {
		u.Abort("Received result set for incorrect bucket")
	}
	return repoResult.Buckets[0].Id
}

func (u SomaUtil) GetRepositoryIdForBucket(bucket string) string {
	var req somaproto.ProtoRequestBucket
	receivedUuidArgument := false

	id, err := uuid.FromString(bucket)
	if err != nil {
		req.Filter.Name = bucket
	} else {
		receivedUuidArgument = true
		req.Filter.Id = id.String()
	}

	resp := u.GetRequestWithBody(req, "/buckets/")
	bucketResult := u.DecodeProtoResultBucketFromResponse(resp)

	if receivedUuidArgument {
		if bucket != bucketResult.Buckets[0].Id {
			u.Abort("Received result set for incorrect bucket")
		}
	} else {
		if bucket != bucketResult.Buckets[0].Name {
			u.Abort("Received result set for incorrect bucket")
		}
	}
	return bucketResult.Buckets[0].Repository
}

func (u SomaUtil) DecodeProtoResultBucketFromResponse(resp *resty.Response) *somaproto.ProtoResultBucket {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	var res somaproto.ProtoResultBucket
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
