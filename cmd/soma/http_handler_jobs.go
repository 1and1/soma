package main

import (
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// JobDelay function
func JobDelay(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `job`,
		Action:     `wait`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan bool)
	handler := handlerMap[`jobDelay`].(*jobDelay)
	handler.input <- waitSpec{
		JobId: params.ByName(`jobid`),
		Reply: returnChannel,
	}
	<-returnChannel
	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
}

// JobList function
func JobList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `job`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`jobs_r`].(*jobsRead)
	handler.input <- msg.Request{
		Section:    `job`,
		Action:     `list`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// JobListAll function
func JobListAll(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `runtime`,
		Action:     `job_list_all`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`jobs_r`].(*jobsRead)
	handler.input <- msg.Request{
		Section:    `runtime`,
		Action:     `job_list_all`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// JobShow function
func JobShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `job`,
		Action:     `show`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	withDetails := IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `job`,
		Action:     `details`,
	})

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`jobs_r`].(*jobsRead)
	handler.input <- msg.Request{
		Section:    `job`,
		Action:     `show`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Flag:       msg.Flags{JobDetail: withDetails},
		Job:        proto.Job{Id: params.ByName(`jobid`)},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// JobSearch function
func JobSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `job`,
		Action:     `search`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	withDetails := IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `job`,
		Action:     `details`,
	})

	cReq := proto.NewJobFilter()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`jobs_r`].(*jobsRead)
	handler.input <- msg.Request{
		Section:    `job`,
		Action:     `search/idlist`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Flag:       msg.Flags{JobDetail: withDetails},
		Search: msg.Filter{
			IsDetailed: cReq.Flags.Detailed,
			Job: proto.JobFilter{
				IdList: cReq.Filter.Job.IdList,
			},
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
