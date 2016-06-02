package util

import (
	"fmt"

	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetBucketByUUIDOrName(c *resty.Client, b string, r string) string {
	if u.IsUUID(b) {
		return b
	}
	return u.GetBucketIdByName(c, b, r)
}

func (u SomaUtil) BucketByUUIDOrName(c *resty.Client, b string) string {
	if u.IsUUID(b) {
		return b
	}
	return u.BucketIdByName(c, b)
}

func (u SomaUtil) GetBucketIdByName(c *resty.Client, bucket string, repoId string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Bucket: &proto.BucketFilter{
				Name:         bucket,
				RepositoryId: repoId,
			},
		},
	}

	resp := u.PostRequestWithBody(c, req, "/filter/buckets/")
	repoResult := u.DecodeProtoResultBucketFromResponse(resp)

	if bucket != (*repoResult.Buckets)[0].Name {
		u.Abort("Received result set for incorrect bucket")
	}
	return (*repoResult.Buckets)[0].Id
}

func (u SomaUtil) BucketIdByName(c *resty.Client, bucket string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Bucket: &proto.BucketFilter{
				Name: bucket,
			},
		},
	}

	resp := u.PostRequestWithBody(c, req, "/filter/buckets/")
	repoResult := u.DecodeProtoResultBucketFromResponse(resp)

	if bucket != (*repoResult.Buckets)[0].Name {
		u.Abort("Received result set for incorrect bucket")
	}
	return (*repoResult.Buckets)[0].Id
}

func (u SomaUtil) GetRepositoryIdForBucket(c *resty.Client, bucket string) string {
	req := proto.NewBucketFilter()
	var b string

	if !u.IsUUID(bucket) {
		req.Filter.Bucket.Name = bucket
		resp := u.PostRequestWithBody(c, req, "/filter/buckets/")
		bucketResult := u.DecodeProtoResultBucketFromResponse(resp)
		if bucket != (*bucketResult.Buckets)[0].Name {
			u.Abort("Received result set for incorrect bucket")
		}
		b = (*bucketResult.Buckets)[0].Id
	} else {
		b = bucket
	}

	path := fmt.Sprintf("/buckets/%s", b)
	resp := u.GetRequest(c, path)
	bucketResult := u.DecodeProtoResultBucketFromResponse(resp)

	return (*bucketResult.Buckets)[0].RepositoryId
}

func (u SomaUtil) TeamIdForBucket(c *resty.Client, bucket string) string {
	req := proto.NewBucketFilter()
	var b string

	if !u.IsUUID(bucket) {
		req.Filter.Bucket.Name = bucket
		resp := u.PostRequestWithBody(c, req, "/filter/buckets/")
		bucketResult := u.DecodeProtoResultBucketFromResponse(resp)
		if bucket != (*bucketResult.Buckets)[0].Name {
			u.Abort("Received result set for incorrect bucket")
		}
		b = (*bucketResult.Buckets)[0].Id
	} else {
		b = bucket
	}

	path := fmt.Sprintf("/buckets/%s", b)
	resp := u.GetRequest(c, path)
	bucketResult := u.DecodeProtoResultBucketFromResponse(resp)

	return (*bucketResult.Buckets)[0].TeamId
}

func (u SomaUtil) DecodeProtoResultBucketFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
