/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"encoding/json"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/auth"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

// SendMsgResult is the output function for all requests that did not
// fail input validation and got processes by the application.
func SendMsgResult(w *http.ResponseWriter, r *msg.Result) {
	var (
		bjson  []byte
		err    error
		k      auth.Kex
		result proto.Result
	)

	// this is central error command, proceeding to log
	if r.Error != nil {
		log.Printf(LogStrErr, r.Section, r.Action, r.Code, r.Error.Error())
	}

	switch r.Section {
	case `kex`:
		k = r.Super.Kex
		if bjson, err = json.Marshal(&k); err != nil {
			log.Printf(LogStrErr, r.Section, r.Action, r.Code, err.Error())
			DispatchInternalError(w, nil)
			return
		}
		goto dispatchJSON
	case `bootstrap`, `activate`, `token`, `password`:
		// for this request type, errors are masked in responses
		switch r.Code {
		case 200:
			if r.Super.Verdict == 200 {
				log.Printf(LogStrOK, r.Section, r.Action, r.Code, 200)
				goto dispatchOCTET
			}
			log.Printf(LogStrOK, r.Section, r.Action, r.Code, 401)
			DispatchUnauthorized(w, nil)
		case 400:
			log.Printf(LogStrOK, r.Section, r.Action, r.Code, 400)
			DispatchBadRequest(w, nil)
		case 404:
			log.Printf(LogStrOK, r.Section, r.Action, r.Code, 404)
			DispatchNotFound(w, r.Error)
		case 406:
			log.Printf(LogStrOK, r.Section, r.Action, r.Code, 406)
			DispatchConflict(w, r.Error)
		default:
			log.Printf(LogStrOK, r.Section, r.Action, r.Code, 401)
			DispatchUnauthorized(w, nil)
		}
		return
	case `category`:
		result = proto.NewCategoryResult()
		*result.Categories = append(*result.Categories, r.Category...)
	case `permission`:
		result = proto.NewPermissionResult()
		*result.Permissions = append(*result.Permissions, r.Permission...)
	case `right`:
		result = proto.NewGrantResult()
		*result.Grants = append(*result.Grants, r.Grant...)
	case `section`:
		result = proto.NewSectionResult()
		*result.Sections = append(*result.Sections, r.SectionObj...)
	case `action`:
		result = proto.NewActionResult()
		*result.Actions = append(*result.Actions, r.ActionObj...)
	case `attribute`:
		result = proto.NewAttributeResult()
		*result.Attributes = append(*result.Attributes, r.Attribute...)
	case `environment`:
		result = proto.NewEnvironmentResult()
		*result.Environments = append(*result.Environments, r.Environment...)
	case `job`:
		result = proto.NewJobResult()
		*result.Jobs = append(*result.Jobs, r.Job...)
	case `tree`:
		result = proto.NewTreeResult()
		*result.Tree = r.Tree
	case `runtime`:
		switch r.Action {
		case `instance_list_all`:
			result = proto.NewInstanceResult()
			*result.Instances = append(*result.Instances, r.Instance...)
		case `job_list_all`:
			result = proto.NewJobResult()
			*result.Jobs = append(*result.Jobs, r.Job...)
		default:
			result = proto.NewSystemOperationResult()
			*result.SystemOperations = append(*result.SystemOperations, r.System...)
		}
	case `instance`:
		result = proto.NewInstanceResult()
		*result.Instances = append(*result.Instances, r.Instance...)
	case `workflow`:
		result = proto.NewWorkflowResult()
		*result.Workflows = append(*result.Workflows, r.Workflow...)
	case `state`:
		result = proto.NewStateResult()
		*result.States = append(*result.States, r.State...)
	case `entity`:
		result = proto.NewEntityResult()
		*result.Entities = append(*result.Entities, r.Entity...)
	default:
		log.Printf(LogStrErr, r.Section, r.Action, 0, `Result from unhandled subsystem`)
		DispatchInternalError(w, nil)
		return
	}

	switch r.Code {
	case 200:
		log.Printf(LogStrOK, r.Section, r.Action, r.Code, 200)
		if r.Error != nil {
			result.Error(r.Error)
		}
		result.OK()
	case 202:
		log.Printf(LogStrOK, r.Section, r.Action, r.Code, 202)
		result.JobId = r.JobId
		result.Accepted()
	case 400:
		log.Printf(LogStrOK, r.Section, r.Action, r.Code, 400)
		DispatchBadRequest(w, nil)
		return
	case 403:
		log.Printf(LogStrOK, r.Section, r.Action, r.Code, 403)
		DispatchForbidden(w, r.Error)
		return
	case 404:
		log.Printf(LogStrOK, r.Section, r.Action, r.Code, 200)
		result.NotFound()
	case 406:
		log.Printf(LogStrOK, r.Section, r.Action, r.Code, 406)
		DispatchConflict(w, r.Error)
		return
	case 500:
		log.Printf(LogStrOK, r.Section, r.Action, r.Code, 500)
		result.Error(r.Error)
	case 501:
		log.Printf(LogStrOK, r.Section, r.Action, r.Code, 501)
		result.NotImplemented()
	case 503:
		log.Printf(LogStrOK, r.Section, r.Action, r.Code, 503)
		result.Unavailable()
	default:
		log.Printf(LogStrErr, r.Section, r.Action, r.Code, `Unhandled internal result code`)
		DispatchInternalError(w, nil)
		return
	}
	goto buildJSON

dispatchOCTET:
	DispatchOctetReply(w, &r.Super.Data)
	return

buildJSON:
	if bjson, err = json.Marshal(&result); err != nil {
		log.Printf(LogStrErr, r.Section, r.Action, r.Code, err)
		DispatchInternalError(w, nil)
		return
	}

dispatchJSON:
	DispatchJsonReply(w, &bjson)
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
