package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetRepositoryByUUIDOrName(s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetRepositoryIdByName(s)
	}
	return id.String()
}

func (u SomaUtil) GetRepositoryIdByName(repo string) string {
	var req somaproto.ProtoRequestRepository
	req.Filter = &somaproto.ProtoRepositoryFilter{}
	req.Filter.Name = repo

	resp := u.GetRequestWithBody(req, "/repository/")
	repoResult := u.DecodeProtoResultRepositoryFromResponse(resp)

	if repo != repoResult.Repositories[0].Name {
		u.Abort("Received result set for incorrect rep[ository")
	}
	return repoResult.Repositories[0].Id
}

func (u SomaUtil) DecodeProtoResultRepositoryFromResponse(resp *resty.Response) *somaproto.ProtoResultRepository {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	var res somaproto.ProtoResultRepository
	err := decoder.Decode(&res)
	u.AbortOnError(err, "Error decoding server response body")
	if res.Code > 299 {
		s := fmt.Sprintf("Request failed: %d - %s", res.Code, res.Status)
		msgs := []string{s}
		msgs = append(msgs, res.Text...)
		u.Abort(msgs...)
	}
	return &res
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
