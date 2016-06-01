/*-
Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package adm

import (
	"encoding/json"
	"fmt"

	"gopkg.in/resty.v0"
)

func RequestToken(c *resty.Client, a *auth.Token) (*auth.Token, error) {
	var (
		kex  *auth.Kex
		err  error
		resp *resty.Response
		body []byte
	)
	jBytes := &[]byte{}
	cipher := &[]byte{}
	plain := &[]byte{}
	cred := &auth.Token{}

	if *jBytes, err = json.Marshal(a); err != nil {
		return nil, err
	}

	// establish key exchange for credential transmission
	if kex, err = KeyExchange(c); err != nil {
		return nil, err
	}

	// encrypt credentials
	if err = kex.EncryptAndEncode(jBytes, cipher); err != nil {
		return nil, err
	}

	// Send request
	if resp, err = c.R().
		SetHeader(`Content-Type`, `application/octet-stream`).
		SetBody(*cipher).
		Put(fmt.Sprintf(
			"/authenticate/token/%s", kex.Request.String())); err != nil {
		return nil, err
	} else if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("Token request failed with status code: %d", resp.StatusCode())
	}

	// decrypt reply
	body = resp.Body()
	if err = kex.DecodeAndDecrypt(&body, plain); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(*plain, *cred); err != nil {
		return nil, err
	}

	// validate token
	if resp, err = c.R().
		SetBasicAuth(a.UserName, cred.Token).
		Get(`/authenticate/validate/`); err != nil {
		return nil, err
	} else if resp.StatusCode() != 204 {
		return nil, fmt.Errorf("Token invalid (Code: %d)", resp.StatusCode())
	}

	// credentials are valid
	return cred, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
