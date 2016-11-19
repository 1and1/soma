package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// BucketList function
func BucketList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `bucket`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["bucketReadHandler"].(*somaBucketReadHandler)
	handler.input <- somaBucketRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := proto.Request{
		Filter: &proto.Filter{
			Bucket: &proto.BucketFilter{},
		},
	}
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if (cReq.Filter.Bucket.Name != "") || (cReq.Filter.Bucket.Id != "") {
		filtered := []somaBucketResult{}
		for _, i := range result.Buckets {
			if (i.Bucket.Name == cReq.Filter.Bucket.Name) || (i.Bucket.Id == cReq.Filter.Bucket.Id) {
				filtered = append(filtered, i)
			}
		}
		result.Buckets = filtered
	}

skip:
	SendBucketReply(&w, &result)
}

// BucketShow function
func BucketShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `bucket`,
		Action:     `show`,
		BucketID:   params.ByName(`bucket`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["bucketReadHandler"].(*somaBucketReadHandler)
	handler.input <- somaBucketRequest{
		action: "show",
		reply:  returnChannel,
		Bucket: proto.Bucket{
			Id: params.ByName("bucket"),
		},
	}
	result := <-returnChannel
	SendBucketReply(&w, &result)
}

// BucketCreate function
func BucketCreate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.Request{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	nameLen := utf8.RuneCountInString(cReq.Bucket.Name)
	if nameLen < 4 || nameLen > 512 {
		DispatchBadRequest(&w, fmt.Errorf(`Illegal bucket name length (4 < x <= 512)`))
		return
	}

	if !IsAuthorized(&msg.Authorization{
		User:         params.ByName(`AuthenticatedUser`),
		RemoteAddr:   extractAddress(r.RemoteAddr),
		Section:      `bucket`,
		Action:       `create`,
		RepositoryID: cReq.Bucket.RepositoryId,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(*guidePost)
	handler.input <- treeRequest{
		RequestType: "bucket",
		Action:      "create_bucket",
		User:        params.ByName(`AuthenticatedUser`),
		reply:       returnChannel,
		Bucket: somaBucketRequest{
			action: "add",
			Bucket: *cReq.Bucket,
		},
	}
	result := <-returnChannel
	SendBucketReply(&w, &result)
}

// BucketAddProperty function
func BucketAddProperty(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `bucket`,
		Action:     `add_property`,
		BucketID:   params.ByName(`bucket`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.Request{}
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
	case params.ByName("type") != (*cReq.Bucket.Properties)[0].Type:
		DispatchBadRequest(&w,
			fmt.Errorf("Mismatched property types: %s, %s",
				params.ByName("type"),
				(*cReq.Bucket.Properties)[0].Type))
		return
	case (params.ByName("type") == "service") && (*cReq.Bucket.Properties)[0].Service.Name == "":
		DispatchBadRequest(&w,
			fmt.Errorf("Empty service name is invalid"))
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(*guidePost)
	handler.input <- treeRequest{
		RequestType: "bucket",
		Action:      fmt.Sprintf("add_%s_property_to_bucket", params.ByName("type")),
		User:        params.ByName(`AuthenticatedUser`),
		reply:       returnChannel,
		Bucket: somaBucketRequest{
			action: fmt.Sprintf("%s_property_new", params.ByName("type")),
			Bucket: *cReq.Bucket,
		},
	}
	result := <-returnChannel
	SendBucketReply(&w, &result)
}

// BucketRemoveProperty function
func BucketRemoveProperty(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `bucket`,
		Action:     `remove_property`,
		BucketID:   params.ByName(`bucket`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	bucket := &proto.Bucket{
		Id: params.ByName(`bucket`),
		Properties: &[]proto.Property{
			proto.Property{
				Type:             params.ByName(`type`),
				BucketId:         params.ByName(`bucket`),
				SourceInstanceId: params.ByName(`source`),
			},
		},
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(*guidePost)
	handler.input <- treeRequest{
		RequestType: `bucket`,
		Action: fmt.Sprintf("delete_%s_property_from_bucket",
			params.ByName("type")),
		User:  params.ByName(`AuthenticatedUser`),
		reply: returnChannel,
		Bucket: somaBucketRequest{
			action: fmt.Sprintf("%s_property_remove",
				params.ByName("type")),
			Bucket: *bucket,
		},
	}
	result := <-returnChannel
	SendRepositoryReply(&w, &result)
}

// SendBucketReply function
func SendBucketReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.Result{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	if result.Errors == nil {
		result.Errors = &[]string{}
	}
	result.Buckets = &[]proto.Bucket{}
	for _, i := range (*r).Buckets {
		*result.Buckets = append(*result.Buckets, i.Bucket)
		if i.ResultError != nil {
			*result.Errors = append(*result.Errors, i.ResultError.Error())
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
