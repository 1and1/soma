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

func CheckErrorHandler(res interface{}, protoRes interface{}) bool {

	switch res.(type) {
	case *[]somaLevelResult:
		r := res.(*[]somaLevelResult)
		t := protoRes.(*somaproto.ProtoResultLevel)

		if len(*r) == 0 {
			t.Code = 404
			t.Status = "NOTFOUND"
			return true
		}
		// r has elements
		if (*r)[0].rErr != nil {
			t.Code = 500
			t.Status = "ERROR"
			t.Text = make([]string, 0)
			t.Text = append(t.Text, (*r)[0].rErr.Error())
			return true
		}
		t.Code = 200
		t.Status = "OK"
		return false
	default:
		log.Println(reflect.TypeOf(res))
	}
	return false
}

func DecodeJsonBody(r *http.Request, s interface{}) error {
	decoder := json.NewDecoder(r.Body)
	var err error

	switch s.(type) {
	case *somaproto.ProtoRequestLevel:
		c := s.(*somaproto.ProtoRequestLevel)
		err = decoder.Decode(c)
	default:
		rt := reflect.TypeOf(s)
		return fmt.Errorf("Unhandled request type: %s", rt)
	}
	if err != nil {
		return err
	}
	return nil
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
