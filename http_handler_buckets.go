package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*
 * Read functions
 */
func ListBucket(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["bucketReadHandler"].(somaBucketReadHandler)
	handler.input <- somaBucketRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := somaproto.ProtoRequestBucket{}
	cReq.Filter = &somaproto.ProtoBucketFilter{}
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if (cReq.Filter.Name != "") || (cReq.Filter.Id != "") {
		filtered := make([]somaBucketResult, 0)
		for _, i := range result.Buckets {
			if (i.Bucket.Name == cReq.Filter.Name) || (i.Bucket.Id == cReq.Filter.Id) {
				filtered = append(filtered, i)
			}
		}
		result.Buckets = filtered
	}

skip:
	SendBucketReply(&w, &result)
}

func ShowBucket(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["bucketReadHandler"].(somaBucketReadHandler)
	handler.input <- somaBucketRequest{
		action: "show",
		reply:  returnChannel,
		Bucket: somaproto.ProtoBucket{
			Id: params.ByName("bucket"),
		},
	}
	result := <-returnChannel
	SendBucketReply(&w, &result)
}

/* Write functions
 */
func AddBucket(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestBucket{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(guidePost)
	handler.input <- treeRequest{
		RequestType: "bucket",
		Action:      "create_bucket",
		reply:       returnChannel,
		Bucket: somaBucketRequest{
			action: "add",
			Bucket: *cReq.Bucket,
		},
	}
	result := <-returnChannel
	SendBucketReply(&w, &result)
}

func AddPropertyToBucket(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestBucket{}
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	switch {
	case params.ByName("bucket") != cReq.Bucket.Id:
		DispatchBadRequest(&w,
			fmt.Errorf("Mismatched bucket ids: %s, %s",
				params.ByName("bucket"),
				cReq.Bucket.Id))
		return
	case len(*cReq.Bucket.Properties) != 1:
		DispatchBadRequest(&w,
			fmt.Errorf("Expected property count 1, actual count: %d",
				len(*cReq.Bucket.Properties)))
		return
	case params.ByName("type") != (*cReq.Bucket.Properties)[0].PropertyType:
		DispatchBadRequest(&w,
			fmt.Errorf("Mismatched property types: %s, %s",
				params.ByName("type"),
				(*cReq.Bucket.Properties)[0].PropertyType))
		return
	case (params.ByName("type") == "service") && (*cReq.Bucket.Properties)[0].Service.Name == "":
		DispatchBadRequest(&w,
			fmt.Errorf("Empty service name is invalid"))
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(guidePost)
	handler.input <- treeRequest{
		RequestType: "bucket",
		Action:      fmt.Sprintf("add_%s_property_to_bucket", params.ByName("type")),
		reply:       returnChannel,
		Bucket: somaBucketRequest{
			action: fmt.Sprintf("%s_property_new", params.ByName("type")),
			Bucket: *cReq.Bucket,
		},
	}
	result := <-returnChannel
	SendBucketReply(&w, &result)
}

/*
 * Utility
 */
func SendBucketReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.ProtoResultBucket{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.Buckets = make([]somaproto.ProtoBucket, 0)
	for _, i := range (*r).Buckets {
		result.Buckets = append(result.Buckets, i.Bucket)
		if i.ResultError != nil {
			result.Text = append(result.Text, i.ResultError.Error())
		}
	}

dispatch:
	json, err := json.Marshal(result)
	if err != nil {
		DispatchInternalError(w, err)
		return
	}
	DispatchJsonReply(w, &json)
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
