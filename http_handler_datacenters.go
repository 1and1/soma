package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/*
 * Read functions
 */
func ListDatacenters(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`datacenters_list`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["datacenterReadHandler"].(somaDatacenterReadHandler)
	handler.input <- somaDatacenterRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendDatacenterReply(&w, &result)
}

func SyncDatacenters(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`datacenters_sync`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["datacenterReadHandler"].(somaDatacenterReadHandler)
	handler.input <- somaDatacenterRequest{
		action: `sync`,
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendDatacenterReply(&w, &result)
}

func ListDatacenterGroups(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	/*
		returnChannel := make(chan []somaDatacenterResult)

		handler := handlerMap["datacenterReadHandler"].(somaDatacenterReadHandler)
		handler.input <- somaDatacenterRequest{
			action: "grouplist",
			reply:  returnChannel,
		}

		results := <-returnChannel
		datacenters := make([]string, len(results))
		for pos, res := range results {
			datacenters[pos] = res.datacenter
		}
		json, err := json.Marshal(proto.ProtoResultDatacenterList{
			Code:        200,
			Status:      "OK",
			Datacenters: datacenters,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	*/
	result := somaResult{}
	result.SetNotImplemented()
	SendDatacenterReply(&w, &result)
}

func ShowDatacenter(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`datacenters_show`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["datacenterReadHandler"].(somaDatacenterReadHandler)
	handler.input <- somaDatacenterRequest{
		action: "show",
		Datacenter: proto.Datacenter{
			Locode: params.ByName("datacenter"),
		},
		reply: returnChannel,
	}
	result := <-returnChannel
	SendDatacenterReply(&w, &result)
}

func ShowDatacenterGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	/*
		returnChannel := make(chan []somaDatacenterResult)

		handler := handlerMap["datacenterReadHandler"].(somaDatacenterReadHandler)
		handler.input <- somaDatacenterRequest{
			action:     "groupshow",
			datacenter: params.ByName("datacentergroup"),
			reply:      returnChannel,
		}

		results := <-returnChannel
		datacenters := make([]string, len(results))
		for pos, res := range results {
			datacenters[pos] = res.datacenter
		}
		json, err := json.Marshal(proto.ProtoResultDatacenterDetail{
			Code:   200,
			Status: "OK",
			Details: proto.ProtoDatacenterDetails{
				Datacenter: params.ByName("datacentergroup"),
				UsedBy:     datacenters,
			},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	*/
	result := somaResult{}
	result.SetNotImplemented()
	SendDatacenterReply(&w, &result)
}

/*
 * Write Functions
 */
func AddDatacenter(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`datacenters_create`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.Request{}
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["datacenterWriteHandler"].(somaDatacenterWriteHandler)
	handler.input <- somaDatacenterRequest{
		action: "add",
		reply:  returnChannel,
		Datacenter: proto.Datacenter{
			Locode: cReq.Datacenter.Locode,
		},
	}
	result := <-returnChannel
	SendDatacenterReply(&w, &result)
}

func AddDatacenterToGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	result := somaResult{}
	result.SetNotImplemented()
	SendDatacenterReply(&w, &result)
	/*
		returnChannel := make(chan []somaDatacenterResult)

		// read PATCH body
		decoder := json.NewDecoder(r.Body)
		var clientRequest proto.ProtoRequestDatacenter
		err := decoder.Decode(&clientRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}

		handler := handlerMap["datacenterWriteHandler"].(somaDatacenterWriteHandler)
		handler.input <- somaDatacenterRequest{
			action:     "groupadd",
			datacenter: clientRequest.Datacenter,
			group:      params.ByName("datacentergroup"),
			reply:      returnChannel,
		}

		results := <-returnChannel
		if len(results) != 1 {
			json, _ := json.Marshal(proto.ProtoResultDatacenter{
				Code:   500,
				Status: "Internal Server Error",
				Text:   []string{"Database statement returned no/wrong number of results"},
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(json)
			return
		}

		result := results[0]
		if result.err != nil {
			json, _ := json.Marshal(proto.ProtoResultDatacenter{
				Code:   500,
				Status: "Internal Server Error",
				Text:   []string{result.err.Error()},
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(json)
			return
		}

		txt := fmt.Sprintf("Added datacenter %s to group %s",
			result.datacenter,
			params.ByName("datacentergroup"))
		json, _ := json.Marshal(proto.ProtoResultDatacenter{
			Code:   200,
			Status: "OK",
			Text:   []string{txt},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	*/
}

func DeleteDatacenter(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`datacenters_delete`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["datacenterWriteHandler"].(somaDatacenterWriteHandler)
	handler.input <- somaDatacenterRequest{
		action: "delete",
		reply:  returnChannel,
		Datacenter: proto.Datacenter{
			Locode: params.ByName("datacenter"),
		},
	}
	result := <-returnChannel
	SendDatacenterReply(&w, &result)
}

func DeleteDatacenterFromGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	result := somaResult{}
	result.SetNotImplemented()
	SendDatacenterReply(&w, &result)
	/*
		returnChannel := make(chan []somaDatacenterResult)

		// read DELETE body
		decoder := json.NewDecoder(r.Body)
		var clientRequest proto.ProtoRequestDatacenter
		err := decoder.Decode(&clientRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}

		handler := handlerMap["datacenterWriteHandler"].(somaDatacenterWriteHandler)
		handler.input <- somaDatacenterRequest{
			action:     "groupdel",
			datacenter: clientRequest.Datacenter,
			group:      params.ByName("datacentergroup"),
			reply:      returnChannel,
		}

		results := <-returnChannel
		if len(results) != 1 {
			json, _ := json.Marshal(proto.ProtoResultDatacenter{
				Code:   500,
				Status: "Internal Server Error",
				Text:   []string{"Database statement returned no/wrong number of results"},
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(json)
			return
		}

		result := results[0]
		if result.err != nil {
			json, _ := json.Marshal(proto.ProtoResultDatacenter{
				Code:   500,
				Status: "Internal Server Error",
				Text:   []string{result.err.Error()},
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(json)
			return
		}

		txt := fmt.Sprintf("Deleted datacenter %s from group %s",
			result.datacenter,
			params.ByName("datacentergroup"))
		json, _ := json.Marshal(proto.ProtoResultDatacenter{
			Code:   200,
			Status: "OK",
			Text:   []string{txt},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	*/
}

func RenameDatacenter(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`datacenters_rename`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.Request{}
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["datacenterWriteHandler"].(somaDatacenterWriteHandler)
	handler.input <- somaDatacenterRequest{
		action: "rename",
		Datacenter: proto.Datacenter{
			Locode: params.ByName("datacenter"),
		},
		rename: cReq.Datacenter.Locode,
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendDatacenterReply(&w, &result)
}

/* Utility
 */
func SendDatacenterReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewDatacenterResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Datacenters {
		*result.Datacenters = append(*result.Datacenters, i.Datacenter)
		if i.ResultError != nil {
			*result.Errors = append(*result.Errors, i.ResultError.Error())
		}
	}

dispatch:
	result.Clean()
	json, err := json.Marshal(result)
	if err != nil {
		DispatchInternalError(w, err)
		return
	}
	DispatchJsonReply(w, &json)
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
