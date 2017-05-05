package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/auth"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
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
	case *proto.Request:
		c := s.(*proto.Request)
		err = decoder.Decode(c)
	case *auth.Kex:
		c := s.(*auth.Kex)
		err = decoder.Decode(c)
	default:
		rt := reflect.TypeOf(s)
		err = fmt.Errorf("DecodeJsonBody: Unhandled request type: %s", rt)
	}
	return err
}

func ResultLength(r *somaResult, t ErrorMarker) int {
	switch t.(type) {
	case *proto.Result:
		switch {
		case r.Datacenters != nil:
			return len(r.Datacenters)
		case r.Levels != nil:
			return len(r.Levels)
		case r.Predicates != nil:
			return len(r.Predicates)
		case r.Status != nil:
			return len(r.Status)
		case r.Oncall != nil:
			return len(r.Oncall)
		case r.Teams != nil:
			return len(r.Teams)
		case r.Nodes != nil:
			return len(r.Nodes)
		case r.Views != nil:
			return len(r.Views)
		case r.Servers != nil:
			return len(r.Servers)
		case r.Units != nil:
			return len(r.Units)
		case r.Providers != nil:
			return len(r.Providers)
		case r.Metrics != nil:
			return len(r.Metrics)
		case r.Modes != nil:
			return len(r.Modes)
		case r.Users != nil:
			return len(r.Users)
		case r.Capabilities != nil:
			return len(r.Capabilities)
		case r.Properties != nil:
			return len(r.Properties)
		case r.Repositories != nil:
			return len(r.Repositories)
		case r.Buckets != nil:
			return len(r.Buckets)
		case r.Groups != nil:
			return len(r.Groups)
		case r.Clusters != nil:
			return len(r.Clusters)
		case r.CheckConfigs != nil:
			return len(r.CheckConfigs)
		case r.Validity != nil:
			return len(r.Validity)
		case r.HostDeployments != nil:
			if len(r.Deployments) > len(r.HostDeployments) {
				return len(r.Deployments)
			}
			return len(r.HostDeployments)
		case r.Deployments != nil:
			return len(r.Deployments)
		}
	default:
		return 0
	}
	return 0
}

func DispatchBadRequest(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Error(*w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func DispatchUnauthorized(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusUnauthorized)
		return
	}
	http.Error(*w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func DispatchForbidden(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusForbidden)
		return
	}
	http.Error(*w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}

func DispatchNotFound(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusNotFound)
		return
	}
	http.Error(*w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func DispatchConflict(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusConflict)
		return
	}
	http.Error(*w, http.StatusText(http.StatusConflict), http.StatusConflict)
}

func DispatchInternalError(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Error(*w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func DispatchNotImplemented(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusNotImplemented)
		return
	}
	http.Error(*w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func DispatchJsonReply(w *http.ResponseWriter, b *[]byte) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(*b)
}

func DispatchOctetReply(w *http.ResponseWriter, b *[]byte) {
	(*w).Header().Set("Content-Type", `application/octet-stream`)
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

// extractAddress extracts the IP address part of the IP:port string
// set as net/http.Request.RemoteAddr. It handles IPv4 cases like
// 192.0.2.1:48467 and IPv6 cases like [2001:db8::1%lo0]:48467
func extractAddress(str string) string {
	var addr string

	switch {
	case strings.Contains(str, `]`):
		// IPv6 address [2001:db8::1%lo0]:48467
		addr = strings.Split(str, `]`)[0]
		addr = strings.Split(addr, `%`)[0]
		addr = strings.TrimLeft(addr, `[`)
	default:
		// IPv4 address 192.0.2.1:48467
		addr = strings.Split(str, `:`)[0]
	}
	return addr
}

func msgRequest(l *log.Logger, q *msg.Request) {
	l.Printf(LogStrSRq,
		q.Section,
		q.Action,
		q.AuthUser,
		q.RemoteAddr,
	)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
