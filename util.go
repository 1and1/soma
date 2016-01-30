package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime/debug"

)

func PanicCatcher(w http.ResponseWriter) {
	if r := recover(); r != nil {
		log.Printf("%s\n", debug.Stack())
		msg := fmt.Sprintf("PANIC! %s", r)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

func DecodeJsonBody(r *http.Request, s interface{}) error {
	decoder := json.NewDecoder(r.Body)
	var err error

	switch s.(type) {
	case *somaproto.ProtoRequestLevel:
		c := s.(*somaproto.ProtoRequestLevel)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestPredicate:
		c := s.(*somaproto.ProtoRequestPredicate)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestStatus:
		c := s.(*somaproto.ProtoRequestStatus)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestOncall:
		c := s.(*somaproto.ProtoRequestOncall)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestTeam:
		c := s.(*somaproto.ProtoRequestTeam)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestNode:
		c := s.(*somaproto.ProtoRequestNode)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestView:
		c := s.(*somaproto.ProtoRequestView)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestServer:
		c := s.(*somaproto.ProtoRequestServer)
		err = decoder.Decode(c)
	default:
		rt := reflect.TypeOf(s)
		return fmt.Errorf("DecodeJsonBody: Unhandled request type: %s", rt)
	}
	if err != nil {
		return err
	}
	return nil
}

func ResultLength(r *somaResult, t ErrorMarker) int {
	switch t.(type) {
	case *somaproto.ProtoResultLevel:
		return len(r.Levels)
	case *somaproto.ProtoResultPredicate:
		return len(r.Predicates)
	case *somaproto.ProtoResultStatus:
		return len(r.Status)
	case *somaproto.ProtoResultOncall:
		return len(r.Oncall)
	case *somaproto.ProtoResultTeam:
		return len(r.Teams)
	case *somaproto.ProtoResultNode:
		return len(r.Nodes)
	case *somaproto.ProtoResultView:
		return len(r.Views)
	case *somaproto.ProtoResultServer:
		return len(r.Servers)
	}
	return 0
}

func DispatchBadRequest(w *http.ResponseWriter, err error) {
	http.Error(*w, err.Error(), http.StatusBadRequest)
}

func DispatchInternalError(w *http.ResponseWriter, err error) {
	http.Error(*w, err.Error(), http.StatusInternalServerError)
}

func DispatchJsonReply(w *http.ResponseWriter, b *[]byte) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(*b)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
