/*-
Copyright (c) 2016, Jörg Pernfuß <code.jpe@gmail.com>
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

package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"net"
	"net/http"
	"time"

	"github.com/mjolnir42/scrypth64"
)

// Token is the data passed between client and server to
// authenticate the client and issue a token for it that can be used
// as HTTP Basic Auth password.
type Token struct {
	UserName  string `json:"username"`
	Password  string `json:"password,omitempty"`
	Token     string `json:"token,omitempty"`
	ValidFrom string `json:"validFrom,omitempty"`
	ExpiresAt string `json:"expiresAt,omitempty"`
	Salt      string `json:"-"`
	SourceIP  net.IP `json:"-"`
}

// NewToken returns an empty Token
func NewToken() *Token {
	return &Token{}
}

// Generate verifies a the embedded credentials in Token and
// issues a new token to be returned to the user. Calling
// GenerateToken consumes the embedded password regardless of outcome.
// Returns ErrAuth if the password is incorrect.
func (t *Token) Generate(mcf scrypth64.Mcf, key, seed []byte) error {
	var (
		err                 error
		ok                  bool
		stuntSeed, stuntKey []byte
	)
	// do not generate Tokens with an incorrectly configured seed or
	// key
	if seed == nil || key == nil || len(seed) == 0 || len(key) == 0 {
		err = ErrInput
		goto fail
	}

	// net.IP is a bit of a pain since an unset net.IP does not equal
	// nil, and .String() prints `<nil>`
	if t.SourceIP.Equal(net.IP{}) {
		err = ErrInput
		goto fail
	}

	// Username and password have to be set
	if t.UserName == "" || t.Password == "" {
		err = ErrInput
		goto fail
	}

	// stunt data
	stuntSeed = bytes.Repeat([]byte{0x0F}, len(seed))
	stuntKey = bytes.Repeat([]byte{0xAB}, len(key))

	// start of timing critical path
	if ok, err = t.comparePassword(mcf); !ok || err != nil {
		// spend some quality time computing garbage
		t.mixToken(stuntKey, stuntSeed)
		goto fail
	}

	if err = t.mixToken(key, seed); err != nil {
		goto fail
	}

	t.zeroPassword()
	return nil

fail:
	t.zeroPassword()
	if err != nil {
		return err
	}
	return ErrAuth
}

// zeroPassword ensures the Password field is set to the zero value
func (t *Token) zeroPassword() {
	t.Password = ""
}

// comparePassword is used to verify the user supplied password
// against a scrypth64.Mcf
func (t *Token) comparePassword(mcf scrypth64.Mcf) (bool, error) {
	return scrypth64.Verify(t.Password, mcf)
}

// mixToken generates a new password token
func (t *Token) mixToken(key, seed []byte) error {
	var (
		salt, bin, btime []byte
		err              error
		valid, expires   time.Time
	)
	// expiry time
	valid = time.Now().UTC()
	expires = valid.Add(
		time.Duration(TokenExpirySeconds) * time.Second,
	).UTC()
	if btime, err = expires.MarshalBinary(); err != nil {
		goto fail
	}
	t.ValidFrom = valid.Format(rfc3339Milli)
	t.ExpiresAt = expires.Format(rfc3339Milli)

	// add random salt
	salt = make([]byte, 16)
	if _, err = rand.Read(salt); err != nil {
		goto fail
	}
	t.Salt = hex.EncodeToString(salt)

	// compute the token
	bin = computeToken(
		[]byte(t.UserName),
		key,
		seed,
		btime,
		salt,
		[]byte(t.SourceIP.String()),
	)
	t.Token = hex.EncodeToString(bin)

	return nil

fail:
	return err
}

// SetIPAddress records the client's IP address
func (t *Token) SetIPAddress(r *http.Request) {
	t.SetIPAddressString(r.RemoteAddr)
}

// SetIPAddressString records the client's IP address
func (t *Token) SetIPAddressString(addr string) {
	t.SourceIP = net.ParseIP(extractAddress(addr))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
