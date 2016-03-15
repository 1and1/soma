package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime/debug"
	"strings"

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
	case *somaproto.ProtoRequestUnit:
		c := s.(*somaproto.ProtoRequestUnit)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestProvider:
		c := s.(*somaproto.ProtoRequestProvider)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestMetric:
		c := s.(*somaproto.ProtoRequestMetric)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestMode:
		c := s.(*somaproto.ProtoRequestMode)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestUser:
		c := s.(*somaproto.ProtoRequestUser)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestMonitoring:
		c := s.(*somaproto.ProtoRequestMonitoring)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestCapability:
		c := s.(*somaproto.ProtoRequestCapability)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestProperty:
		c := s.(*somaproto.ProtoRequestProperty)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestAttribute:
		c := s.(*somaproto.ProtoRequestAttribute)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestRepository:
		c := s.(*somaproto.ProtoRequestRepository)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestBucket:
		c := s.(*somaproto.ProtoRequestBucket)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestGroup:
		c := s.(*somaproto.ProtoRequestGroup)
		err = decoder.Decode(c)
	case *somaproto.ProtoRequestCluster:
		c := s.(*somaproto.ProtoRequestCluster)
		err = decoder.Decode(c)
	case *somaproto.CheckConfigurationRequest:
		c := s.(*somaproto.CheckConfigurationRequest)
		err = decoder.Decode(c)
	case *somaproto.HostDeploymentRequest:
		c := s.(*somaproto.HostDeploymentRequest)
		err = decoder.Decode(c)
	case *somaproto.ValidityRequest:
		c := s.(*somaproto.ValidityRequest)
		err = decoder.Decode(c)
	default:
		rt := reflect.TypeOf(s)
		//return fmt.Errorf("DecodeJsonBody: Unhandled request type: %s", rt)
		// XXX Dev Setting
		errMsg := fmt.Sprintf("DecodeJsonBody: Unhandled request type: %s", rt)
		log.Fatal(errMsg)
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
	case *somaproto.ProtoResultUnit:
		return len(r.Units)
	case *somaproto.ProtoResultProvider:
		return len(r.Providers)
	case *somaproto.ProtoResultMetric:
		return len(r.Metrics)
	case *somaproto.ProtoResultMode:
		return len(r.Modes)
	case *somaproto.ProtoResultUser:
		return len(r.Users)
	case *somaproto.ProtoResultMonitoring:
		return len(r.Systems)
	case *somaproto.ProtoResultCapability:
		return len(r.Capabilities)
	case *somaproto.ProtoResultProperty:
		return len(r.Properties)
	case *somaproto.ProtoResultAttribute:
		return len(r.Attributes)
	case *somaproto.ProtoResultRepository:
		return len(r.Repositories)
	case *somaproto.ProtoResultBucket:
		return len(r.Buckets)
	case *somaproto.ProtoResultGroup:
		return len(r.Groups)
	case *somaproto.ProtoResultCluster:
		return len(r.Clusters)
	case *somaproto.CheckConfigurationResult:
		return len(r.CheckConfigs)
	case *somaproto.DeploymentDetailsResult:
		return len(r.Deployments)
	case *somaproto.HostDeploymentResult:
		return len(r.Deployments)
	case *somaproto.ValidityResult:
		return len(r.Validity)
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

func GetPropertyTypeFromUrl(u *url.URL) (string, error) {
	// strip surrounding / and skip first path element `property|filter`
	el := strings.Split(strings.Trim(u.Path, "/"), "/")[1:]
	if el[0] == "property" {
		// looks like the path was /filter/property/...
		el = el[1:]
	}
	switch el[0] {
	case "service":
		switch el[1] {
		case "team":
			return "service", nil
		case "global":
			return "template", nil
		default:
			return "", errors.New("Unknown service property type")
		}
	default:
		return el[0], nil
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
