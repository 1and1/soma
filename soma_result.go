package main

import (
	"log"
	"reflect"
)

type ErrorMarker interface {
	ErrorMark(err error, imp bool, found bool, length int) bool
}

type somaResult struct {
	RequestError   error
	NotFound       bool
	NotImplemented bool
	Servers        []somaServerResult
}

func (r *somaResult) SetRequestError(err error) bool {
	if err != nil {
		r.RequestError = err
		return true
	}
	return false
}

func (r *somaResult) SetNotFound() {
	r.NotFound = true
}

func (r *somaResult) SetNotImplemented() {
	r.NotImplemented = true
}

func (r *somaResult) Failure() bool {
	if r.NotFound || r.NotImplemented || r.RequestError != nil {
		return true
	}
	return false
}

func (r *somaResult) Append(err error, res interface{}) {
	switch res.(type) {
	case somaServerResult:
		if err != nil {
			r.Servers = append(r.Servers, somaServerResult{ResultError: err})
			break
		}
		r.Servers = append(r.Servers, res.(somaServerResult))
	default:
		log.Printf("somaResult.Append(): unhandled type %s", reflect.TypeOf(res))
	}
}

func (r *somaResult) MarkErrors(reply ErrorMarker) bool {
	return reply.ErrorMark(r.RequestError, r.NotImplemented, r.NotFound,
		len(r.Servers))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
