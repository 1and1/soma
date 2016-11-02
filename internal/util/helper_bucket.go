package util

import (
	"fmt"

	"github.com/1and1/soma/lib/proto"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) tryGetBucketByUUIDOrName(c *resty.Client, b string, r string) string {
	if u.IsUUID(b) {
		return b
	}
	return u.getBucketIdByName(c, b, r)
}

func (u SomaUtil) bucketByUUIDOrName(c *resty.Client, b string) string {
	if u.IsUUID(b) {
		return b
	}
	return u.bucketIdByName(c, b)
}

func (u SomaUtil) getBucketIdByName(c *resty.Client, bucket string, repoId string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Bucket: &proto.BucketFilter{
				Name:         bucket,
				RepositoryId: repoId,
			},
		},
	}

	resp := u.PostRequestWithBody(c, req, "/filter/buckets/")
	repoResult := u.decodeProtoResultBucketFromResponse(resp)

	if bucket != (*repoResult.Buckets)[0].Name {
		u.abort("Received result set for incorrect bucket")
	}
	return (*repoResult.Buckets)[0].Id
}

func (u SomaUtil) bucketIdByName(c *resty.Client, bucket string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Bucket: &proto.BucketFilter{
				Name: bucket,
			},
		},
	}

	resp := u.PostRequestWithBody(c, req, "/filter/buckets/")
	repoResult := u.decodeProtoResultBucketFromResponse(resp)

	if bucket != (*repoResult.Buckets)[0].Name {
		u.abort("Received result set for incorrect bucket")
	}
	return (*repoResult.Buckets)[0].Id
}

func (u SomaUtil) GetRepositoryIdForBucket(c *resty.Client, bucket string) string {
	req := proto.NewBucketFilter()
	var b string

	if !u.IsUUID(bucket) {
		req.Filter.Bucket.Name = bucket
		resp := u.PostRequestWithBody(c, req, "/filter/buckets/")
		bucketResult := u.decodeProtoResultBucketFromResponse(resp)
		if bucket != (*bucketResult.Buckets)[0].Name {
			u.abort("Received result set for incorrect bucket")
		}
		b = (*bucketResult.Buckets)[0].Id
	} else {
		b = bucket
	}

	path := fmt.Sprintf("/buckets/%s", b)
	resp := u.GetRequest(c, path)
	bucketResult := u.decodeProtoResultBucketFromResponse(resp)

	return (*bucketResult.Buckets)[0].RepositoryId
}

func (u SomaUtil) TeamIdForBucket(c *resty.Client, bucket string) string {
	req := proto.NewBucketFilter()
	var b string

	if !u.IsUUID(bucket) {
		req.Filter.Bucket.Name = bucket
		resp := u.PostRequestWithBody(c, req, "/filter/buckets/")
		bucketResult := u.decodeProtoResultBucketFromResponse(resp)
		if bucket != (*bucketResult.Buckets)[0].Name {
			u.abort("Received result set for incorrect bucket")
		}
		b = (*bucketResult.Buckets)[0].Id
	} else {
		b = bucket
	}

	path := fmt.Sprintf("/buckets/%s", b)
	resp := u.GetRequest(c, path)
	bucketResult := u.decodeProtoResultBucketFromResponse(resp)

	return (*bucketResult.Buckets)[0].TeamId
}

func (u SomaUtil) getBucketDetails(c *resty.Client, bucketId string) *proto.Bucket {
	resp := u.GetRequest(c, fmt.Sprintf("/buckets/%s", bucketId))
	res := u.DecodeResultFromResponse(resp)
	return &(*res.Buckets)[0]
}

func (u SomaUtil) FindSourceForBucketProperty(c *resty.Client, pTyp, pName, view, bucketId string) string {
	bucket := u.getBucketDetails(c, bucketId)
	if bucket == nil {
		return ``
	}
	for _, prop := range *bucket.Properties {
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

func (u SomaUtil) decodeProtoResultBucketFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
