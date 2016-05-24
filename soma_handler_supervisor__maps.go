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

//
//
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

//
//
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
// used for one client/server exchange, they are never
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

//
//
// read/write locked map of global permissions
type svPermMapGlobal struct {
	// user(uuid.string) -> permission(uuid.string) -> true
	GMap  map[string]map[string]bool
	mutex sync.RWMutex
}

func (g *svPermMapGlobal) grant(user, permission string) {
	g.lock()
	defer g.unlock()

	// zero value for maps is nil
	if m, ok := g.GMap[user]; !ok {
		g.GMap[user] = make(map[string]bool)
	} else if m == nil {
		g.GMap[user] = make(map[string]bool)
	}

	// grant permission
	g.GMap[user][permission] = true
}

func (g *svPermMapGlobal) revoke(user, permission string) {
	g.lock()
	defer g.unlock()

	// user has no permissions
	if m, ok := g.GMap[user]; !ok {
		return
	} else if m == nil {
		return
	}

	// revoke permission
	delete(g.GMap[user], permission)
}

// ATTENTION: named return parameter
func (g *svPermMapGlobal) assess(user, permission string) (verdict bool) {
	g.rlock()
	defer g.runlock()
	// map[]map[] is volatile
	defer func() {
		if r := recover(); r != nil {
			verdict = false
		}
	}()
	verdict = false

	if m, ok := g.GMap[user]; !ok {
		g.GMap[user] = make(map[string]bool)
		return
	} else if m == nil {
		g.GMap[user] = make(map[string]bool)
		return
	}

	// let zero value `false` work for us
	verdict = g.GMap[user][permission]
	return
}

func (g *svPermMapGlobal) lock() {
	g.mutex.Lock()
}

func (g *svPermMapGlobal) rlock() {
	g.mutex.RLock()
}

func (g *svPermMapGlobal) unlock() {
	g.mutex.Unlock()
}

func (g *svPermMapGlobal) runlock() {
	g.mutex.RUnlock()
}

//
//
// read/write locked map of limited permissions
type svPermMapLimited struct {
	// user(uuid.string) -> permission(uuid.string) -> repository(uuid.string)
	LMap  map[string]map[string][]string
	mutex sync.RWMutex
}

func (l *svPermMapLimited) grant(user, permission, repository string) {
	l.lock()
	defer l.unlock()

	// zero value for maps is nil
	if m, ok := l.LMap[user]; !ok {
		l.LMap[user] = make(map[string][]string)
	} else if m == nil {
		l.LMap[user] = make(map[string][]string)
	}

	// zero value for slices is nil
	s := l.LMap[user][permission]
	if s == nil {
		l.LMap[user][permission] = []string{repository}
	} else {
		s = append(s, repository)
		l.LMap[user][permission] = s
	}
}

func (l *svPermMapLimited) revoke(user, permission, repository string) {
	l.lock()
	defer l.unlock()

	// initialize to remove a future nilptr-dereference hazard, but
	// return early since the user did not have the permission in the
	// first place
	if m, ok := l.LMap[user]; !ok {
		l.LMap[user] = make(map[string][]string)
		return
	} else if m == nil {
		l.LMap[user] = make(map[string][]string)
		return
	}

	// resolve this nilptr-deref hazard as well before returning
	s := l.LMap[user][permission]
	if s == nil {
		l.LMap[user][permission] = []string{}
		return
	}

repoloop:
	for i, _ := range s {
		if s[i] == repository {
			s = append(s[:i], s[i+1:]...)
			break repoloop
		}
	}
	l.LMap[user][permission] = s
}

// ATTENTION: named return parameter
func (l *svPermMapLimited) assess(user, permission, repository string) (verdict bool) {
	l.rlock()
	defer l.runlock()
	defer func() {
		if r := recover(); r != nil {
			verdict = false
		}
	}()
	verdict = false

	if m, ok := l.LMap[user]; !ok {
		l.LMap[user] = make(map[string][]string)
		return
	} else if m == nil {
		l.LMap[user] = make(map[string][]string)
		return
	}

	s := l.LMap[user][permission]
	if s == nil {
		l.LMap[user][permission] = []string{}
		return
	}

	for _, p := range s {
		if p == repository {
			verdict = true
			return
		}
	}
	return
}

func (l *svPermMapLimited) lock() {
	l.mutex.Lock()
}

func (l *svPermMapLimited) rlock() {
	l.mutex.RLock()
}

func (l *svPermMapLimited) unlock() {
	l.mutex.Unlock()
}

func (l *svPermMapLimited) runlock() {
	l.mutex.RUnlock()
}

//
//
// read/write locked map[string]string
type svLockMap struct {
	// user uuid -> team uuid
	LockMap map[string]string
	mutex   sync.RWMutex
}

func (slm *svLockMap) insert(key, value string) {
	slm.lock()
	defer slm.unlock()
	slm.LockMap[key] = value
}

func (slm *svLockMap) remove(key string) {
	slm.lock()
	defer slm.unlock()
	delete(slm.LockMap, key)
}

func (slm *svLockMap) get(key string) (string, bool) {
	slm.rlock()
	defer slm.runlock()
	value, ok := slm.LockMap[key]
	return value, ok
}

func (slm *svLockMap) scan(value string) (string, bool) {
	slm.rlock()
	defer slm.runlock()
	for key, _ := range slm.LockMap {
		if slm.LockMap[key] == value {
			return key, true
		}
	}
	return "", false
}

// load is an unlocked insert for bulk loading at startup. The
// bulk loading mechanism must handle the locking itself
func (slm *svLockMap) load(key, value string) {
	slm.LockMap[key] = value
}

func (slm *svLockMap) lock() {
	slm.mutex.Lock()
}

func (slm *svLockMap) rlock() {
	slm.mutex.RLock()
}

func (slm *svLockMap) unlock() {
	slm.mutex.Unlock()
}

func (slm *svLockMap) runlock() {
	slm.mutex.RUnlock()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
