package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

)

func SendMsgResult(w *http.ResponseWriter, r *msg.Result) {
	var (
		bjson  []byte
		err    error
		k      auth.Kex
		result proto.Result
	)

	// this is central error command, proceeding to log
	if r.Error != nil {
		log.Printf(LogStrErr, r.Type, r.Action, r.Code, r.Error.Error())
	}

	switch r.Type {
	case `supervisor`:
		switch r.Action {
		case `kex_reply`:
			k = r.Super.Kex
			if bjson, err = json.Marshal(&k); err != nil {
				log.Printf(LogStrErr, r.Type, r.Action, r.Code, err.Error())
				DispatchInternalError(w, nil)
				return
			}
			goto dispatchJSON
		case `bootstrap_root`:
			fallthrough
		case `activate_user`:
			fallthrough
		case `issue_token`:
			// for this request type, errors are masked in responses
			switch r.Code {
			case 200:
				if r.Super.Verdict == 200 {
					log.Printf(LogStrOK, r.Type, r.Action, r.Code, 200)
					goto dispatchOCTET
				}
				log.Printf(LogStrOK, r.Type, r.Action, r.Code, 401)
				DispatchUnauthorized(w, nil)
			case 400:
				log.Printf(LogStrOK, r.Type, r.Action, r.Code, 400)
				DispatchBadRequest(w, nil)
			case 404:
				log.Printf(LogStrOK, r.Type, r.Action, r.Code, 404)
				DispatchNotFound(w, r.Error)
			case 406:
				log.Printf(LogStrOK, r.Type, r.Action, r.Code, 406)
				DispatchConflict(w, r.Error)
			default:
				log.Printf(LogStrOK, r.Type, r.Action, r.Code, 401)
				DispatchUnauthorized(w, nil)
			}
			return
		case `category`:
			result = proto.NewCategoryResult()
			*result.Categories = append(*result.Categories, r.Category...)
			goto UnmaskedReply
		case `permission`:
			result = proto.NewPermissionResult()
			*result.Permissions = append(*result.Permissions, r.Permission...)
			goto UnmaskedReply
		case `right`:
			result = proto.NewGrantResult()
			*result.Grants = append(*result.Grants, r.Grant...)
			goto UnmaskedReply
		default:
			log.Printf(LogStrErr, r.Type, r.Action, 0, `Unhandled supervisor action`)
			// supervisor as auth-lord has special default to avoid
			// accidental leakage
			DispatchUnauthorized(w, nil)
			return
		} // end supervisor
	default:
		log.Printf(LogStrErr, r.Type, ``, 0, `Result from unhandled subsystem`)
		DispatchInternalError(w, nil)
		return
	}

UnmaskedReply:
	switch r.Code {
	case 200:
		log.Printf(LogStrOK, r.Type, fmt.Sprintf("%s/%s", r.Action, r.Super.Action), r.Code, 200)
		if r.Error != nil {
			result.Error(r.Error)
		}
		result.OK()
	case 202:
		log.Printf(LogStrOK, r.Type, fmt.Sprintf("%s/%s", r.Action, r.Super.Action), r.Code, 202)
		result.JobId = r.JobId
		result.Accepted()
	case 400:
		log.Printf(LogStrOK, r.Type, fmt.Sprintf("%s/%s", r.Action, r.Super.Action), r.Code, 400)
		DispatchBadRequest(w, nil)
		return
	case 403:
		log.Printf(LogStrOK, r.Type, fmt.Sprintf("%s/%s", r.Action, r.Super.Action), r.Code, 403)
		DispatchForbidden(w, r.Error)
		return
	case 404:
		log.Printf(LogStrOK, r.Type, fmt.Sprintf("%s/%s", r.Action, r.Super.Action), r.Code, 200)
		result.NotFound()
	case 406:
		log.Printf(LogStrOK, r.Type, fmt.Sprintf("%s/%s", r.Action, r.Super.Action), r.Code, 406)
		DispatchConflict(w, r.Error)
		return
	case 500:
		log.Printf(LogStrOK, r.Type, fmt.Sprintf("%s/%s", r.Action, r.Super.Action), r.Code, 500)
		result.Error(r.Error)
	default:
		log.Printf(LogStrErr, r.Type, r.Action, r.Code, `Unhandled internal result code`)
		DispatchInternalError(w, nil)
		return
	}
	goto buildJSON

dispatchOCTET:
	DispatchOctetReply(w, &r.Super.Data)
	return

buildJSON:
	if bjson, err = json.Marshal(&result); err != nil {
		log.Printf(LogStrErr, r.Type, r.Action, r.Code, err)
		DispatchInternalError(w, nil)
		return
	}

dispatchJSON:
	DispatchJsonReply(w, &bjson)
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
