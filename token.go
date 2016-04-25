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

package somaauth

import "encoding/base64"

// TokenRequest is the data passed between client and server to
// authenticate the client and issue a token for it that can be used
// as HTTP Basic Auth password.
type TokenRequest struct {
	UserId    string `json:"user_id,omitempty"`
	UserName  string `json:"user_name"`
	Password  string `json:"password,omitempty"`
	Token     string `json:"token,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

// TokenExpirySeconds can be set to regulate the lifetime of newly
// issued authentication tokens. The default value is 43200, or 12
// hours.
var TokenExpirySeconds uint64 = 43200

// NewTokenRequest returns an empty TokenRequest
func NewTokenRequest() *TokenRequest {
	return &TokenRequest{}
}

// ZeroPassword ensures the Password field is set to the zero value
func (t *TokenRequest) ZeroPassword() {
	t.Password = ""
}

// ComparePassword
func (t *TokenRequest) ComparePassword(
	salt *[64]byte,
	hash string,
) (bool, error) {
	return false, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
