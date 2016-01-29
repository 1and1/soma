package main

import (
	"log"
	"reflect"

)

/*XXX this needs to somehow go away
 *XXX
 *XXX best idea so far is to add an interface that is implemented
 *XXX by all the somaproto.ProtoResult* types. And each of these
 *XXX types then knows how to check its []soma*Result for errors.
 *XXX
 *XXX But this means that soma has methods for types declared in
 *XXX somaproto. Is that really a win?
 *XXX For now, keep this in its own file so its ugliness does not
 *XXX become contagious.
 */

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
	case *[]somaPredicateResult:
		r := res.(*[]somaPredicateResult)
		t := protoRes.(*somaproto.ProtoResultPredicate)

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
	case *[]somaStatusResult:
		r := res.(*[]somaStatusResult)
		t := protoRes.(*somaproto.ProtoResultStatus)

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
	case *[]somaOncallResult:
		r := res.(*[]somaOncallResult)
		t := protoRes.(*somaproto.ProtoResultOncall)

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
	case *[]somaTeamResult:
		r := res.(*[]somaTeamResult)
		t := protoRes.(*somaproto.ProtoResultTeam)

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
	case *[]somaNodeResult:
		r := res.(*[]somaNodeResult)
		t := protoRes.(*somaproto.ProtoResultNode)

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
		log.Printf("CheckErrorHandler: unhandled type %s", reflect.TypeOf(res))
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
