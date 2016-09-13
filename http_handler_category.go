package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListCategory(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`category_list`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Type:       `supervisor`,
		Action:     `category`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Super: &msg.Supervisor{
			Action: `list`,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

func ShowCategory(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`category_show`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Type:       `supervisor`,
		Action:     `category`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Super: &msg.Supervisor{
			Action: `show`,
		},
		Category: proto.Category{
			Name: params.ByName(`category`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

/* Write functions
 */

func AddCategory(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`category_create`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.Request{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Type:       `supervisor`,
		Action:     `category`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Super: &msg.Supervisor{
			Action: `add`,
		},
		Category: proto.Category{
			Name: cReq.Category.Name,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`category_delete`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Type:       `supervisor`,
		Action:     `category`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Super: &msg.Supervisor{
			Action: `delete`,
		},
		Category: proto.Category{
			Name: params.ByName(`category`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
