package main

import (
	"net/http"


	"github.com/julienschmidt/httprouter"
)

func JobDelay(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

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

/* Read functions
 */
func ListJobs(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	var ok, admin bool
	if ok, admin = IsAuthorized(params.ByName(`AuthenticatedUser`),
		`jobs_list`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`jobs_r`].(*jobsRead)
	handler.input <- msg.Request{
		Type:       `job`,
		Action:     `list`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		IsAdmin:    admin,
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

func ShowJob(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	var ok, admin bool
	if ok, admin = IsAuthorized(params.ByName(`AuthenticatedUser`),
		`jobs_show`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`jobs_r`].(*jobsRead)
	handler.input <- msg.Request{
		Type:       `job`,
		Action:     `show`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		IsAdmin:    admin,
		Job:        proto.Job{Id: params.ByName(`jobid`)},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

func SearchJob(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	var ok, admin bool
	if ok, admin = IsAuthorized(params.ByName(`AuthenticatedUser`),
		`jobs_search`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewJobFilter()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`jobs_r`].(*jobsRead)
	handler.input <- msg.Request{
		Type:       `job`,
		Action:     `search/idlist`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		IsAdmin:    admin,
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
