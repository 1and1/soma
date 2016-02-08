package main

type somaTreeRequest struct {
	RequestType string
	Repository  somaRepositoryRequest
}

type somaTreeResult struct {
	ResultType  string
	ResultError error
	Repository  somaproto.ProtoRepository
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
