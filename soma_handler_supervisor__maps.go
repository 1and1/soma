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

package main

import (
	"encoding/hex"
	"sync"
	"time"


	"github.com/mjolnir42/scrypth64"
	"github.com/satori/go.uuid"
)

// supervisor internal storage format for tokens
type svToken struct {
	validFrom    time.Time
	expiresAt    time.Time
	binToken     []byte
	binExpiresAt []byte
	salt         []byte
}

// read/write locked map of tokens
type svTokenMap struct {
	// token(hex.string) -> svToken
	TMap  map[string]svToken
	mutex sync.RWMutex
}

func (t *svTokenMap) read(token string) *svToken {
	t.rlock()
	defer t.runlock()
	if tok, ok := t.TMap[token]; ok {
		return &tok
	}
	return nil
}

func (t *svTokenMap) insert(token, valid, expires, salt string) error {
	var (
		err                     error
		valTime, expTime        time.Time
		bExpTime, bSalt, bToken []byte
	)
	// convert input data into the different formats required to
	// perform later actions without conversions
	if valTime, err = time.Parse(rfc3339Milli, valid); err != nil {
		return err
	}
	if expTime, err = time.Parse(rfc3339Milli, expires); err != nil {
		return err
	}
	if bExpTime, err = expTime.MarshalBinary(); err != nil {
		return err
	}
	if bToken, err = hex.DecodeString(token); err != nil {
		return err
	}
	if bSalt, err = hex.DecodeString(salt); err != nil {
		return err
	}
	// whiteout unstable subsecond timestamp part with "random" value
	copy(bExpTime[9:], []byte{0xde, 0xad, 0xca, 0xfe})
	// acquire write lock
	t.lock()
	defer t.unlock()

	// insert token
	t.TMap[token] = svToken{
		validFrom:    valTime,
		expiresAt:    expTime,
		binToken:     bToken,
		binExpiresAt: bExpTime,
		salt:         bSalt,
	}
	return nil
}

// set writelock
func (t *svTokenMap) lock() {
	t.mutex.Lock()
}

// set readlock
func (t *svTokenMap) rlock() {
	t.mutex.RLock()
}

// release writelock
func (t *svTokenMap) unlock() {
	t.mutex.Unlock()
}

// release readlock
func (t *svTokenMap) runlock() {
	t.mutex.RUnlock()
}

// supervisor internal storage format for credentials
type svCredential struct {
	id          uuid.UUID
	validFrom   time.Time
	expiresAt   time.Time
	cryptMCF    scrypth64.Mcf
	resetActive bool
	isActive    bool
}

type svCredMap struct {
	// username -> svCredential
	CMap  map[string]svCredential
	mutex sync.RWMutex
}

func (c *svCredMap) read(user string) *svCredential {
	c.rlock()
	defer c.runlock()
	if cred, ok := c.CMap[user]; ok {
		return &cred
	}
	return nil
}

func (c *svCredMap) insert(user string, uid uuid.UUID, valid, expires time.Time, mcf scrypth64.Mcf) {
	c.lock()
	defer c.unlock()
	c.CMap[user] = svCredential{
		id:          uid,
		validFrom:   valid,
		expiresAt:   expires,
		cryptMCF:    mcf,
		resetActive: false,
		isActive:    true,
	}
}

func (c *svCredMap) restore(user string, uid uuid.UUID, valid, expires time.Time, mcf scrypth64.Mcf, reset, active bool) {
	c.lock()
	defer c.unlock()
	c.CMap[user] = svCredential{
		id:          uid,
		validFrom:   valid,
		expiresAt:   expires,
		cryptMCF:    mcf,
		resetActive: reset,
		isActive:    active,
	}
}

// set writelock
func (c *svCredMap) lock() {
	c.mutex.Lock()
}

// set readlock
func (c *svCredMap) rlock() {
	c.mutex.RLock()
}

// release writelock
func (c *svCredMap) unlock() {
	c.mutex.Unlock()
}

// release readlock
func (c *svCredMap) runlock() {
	c.mutex.RUnlock()
}

//
//
// read/write locked map of key exchanges
type svKexMap struct {
	// kexid(uuid.string) -> auth.Kex
	KMap  map[string]auth.Kex
	mutex sync.RWMutex
}

// the nonce information would normally mean returning
// a copy is problematic, but since these keys are only
// used for any client/server exchange, they are never
// put back
func (k *svKexMap) read(kexRequest string) *auth.Kex {
	k.rlock()
	defer k.runlock()
	if kex, ok := k.KMap[kexRequest]; ok {
		return &kex
	}
	return nil
}

func (k *svKexMap) insert(kex auth.Kex) {
	k.lock()
	defer k.unlock()

	k.KMap[kex.Request.String()] = kex
}

func (k *svKexMap) remove(kexRequest string) {
	k.lock()
	defer k.unlock()

	delete(k.KMap, kexRequest)
}

// set writelock
func (k *svKexMap) lock() {
	k.mutex.Lock()
}

// set readlock
func (k *svKexMap) rlock() {
	k.mutex.RLock()
}

// release writelock
func (k *svKexMap) unlock() {
	k.mutex.Unlock()
}

// release readlock
func (k *svKexMap) runlock() {
	k.mutex.RUnlock()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
