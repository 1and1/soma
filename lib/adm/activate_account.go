package adm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/1and1/soma/lib/auth"
	"gopkg.in/resty.v0"
)

func ActivateAccount(c *resty.Client, a *auth.Token) (*auth.Token, error) {
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
			"/authenticate/activate/%s", kex.Request.String())); err != nil {
		return nil, err
	} else if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("Activation failed: %s[%d], %s", http.StatusText(resp.StatusCode()), resp.StatusCode(), resp.String())
	}

	// decrypt reply
	body = resp.Body()
	if err = kex.DecodeAndDecrypt(&body, plain); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(*plain, cred); err != nil {
		return nil, err
	}

	return cred, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
