/*-
 * Copyright © 2016, Jörg Pernfuß <code.jpe@gmail.com>
 * All rights reserved.
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

// scrypth64 derives password hashes and verifies passwords against
// them according to the scrypt-h64 specification.
// It includes a number of ready to use parameter sets.
// The modular crypt format is used, and its de-facto standard as well
// as historic precedents are followed.
package scrypth64

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"strconv"
	"strings"

	"github.com/mjolnir42/hash64"

	"golang.org/x/crypto/scrypt"
)

// Mcf is a custom string type to indicate that the underlying string
// is subject to specific formatting rules
type Mcf string

func (m Mcf) String() string {
	return string(m)
}

// internal constants
const (
	fieldIdentifier = iota
	fieldParameters
	fieldSalt
	fieldDigest
	mcfString string = `$scrypt-h64$N=%d,r=%d,p=%d,l=%d,s=%d$%s$%s`
)

// Digest calculates a new scrypt digest and returns it as a scrypt-h64
// MCF formatted string. If params is nil, DefaultParams are used.
// Returns an empty string if err != nil.
func Digest(password string, params *Parameter) (Mcf, error) {
	var (
		salt, digest []byte
		err          error
	)

	// use default parameters or ensure they are valid
	if params == nil {
		params = DefaultParams()
	} else {
		if err = checkParameter(params); err != nil {
			goto fail
		}
	}

	if salt, err = newSalt(params); err != nil {
		goto fail
	}

	if digest, err = computeDigest(password, salt, params); err != nil {
		goto fail
	}

	return format(digest, salt, params), nil

fail:
	return "", err
}

// Verify takes a password and a scrypt-h64 string in MCF format and
// returns true if the password matches the provided digest.
// The result is always false if error != nil.
func Verify(password string, mcf Mcf) (bool, error) {
	var (
		p                         *Parameter
		salt, digest, pass, stunt []byte
		cmpResult                 int
		result, nop               bool
		err                       error
	)
	// unnecessary but explicit
	result = false

	// get all the required data from the MCF string
	if p, err = parameterFromMCF(mcf); err != nil {
		goto fail
	}

	if salt, err = saltFromMCF(mcf); err != nil {
		goto fail
	}

	if digest, err = digestFromMCF(mcf); err != nil {
		goto fail
	}
	// this byteslice is used if scrypt can not compute the digest of
	// the supplied password for any reason
	stunt = bytes.Repeat([]byte{0x0F}, len(digest))

	// ATTENTION!
	// from here on this function becomes timing critical, since
	// attacker supplied data is processed in the authentication
	// path

	// generate the digest of the supplied password using salt and
	// parameters extracted from the MCF string
	pass, err = computeDigest(password, salt, p)
	if err != nil {
		// there was something wrong, use the stuntslice
		cmpResult = subtle.ConstantTimeCompare(digest, stunt)
	} else {
		cmpResult = subtle.ConstantTimeCompare(digest, pass)
	}

	// avoid short-circuit operators && || and try to get the number
	// of operations for failure and success path the same
	if err == nil {
		if cmpResult == 1 {
			result = true
		} else {
			nop = true
		}
	} else {
		if cmpResult == 2 {
			nop = true
		} else {
			nop = false
		}
	}
	// always return err to unhide errors from computeDigest()
	return result, err

fail:
	// this construct silences `nop declared and not used` errors
	nop = false
	return nop, err
}

// computeDigest calculates a new scrypt digest. Returns nil if there
// was an error computing the digest.
func computeDigest(password string, salt []byte, params *Parameter) ([]byte, error) {
	var (
		digest []byte
		err    error
	)

	// timing critical: processing attacker supplied data

	if digest, err = scrypt.Key(
		[]byte(password),
		salt,
		1<<params.Iterations,
		int(params.BlockSizeFactor),
		int(params.Parallelism),
		int(params.Length),
	); err != nil {
		return nil, err
	}

	return digest, nil
}

// newSalt reads and returns a new random salt of the length requested
// in p. Returns nil if err != nil, indicating an error with the
// entropy source.
func newSalt(p *Parameter) ([]byte, error) {
	salt := make([]byte, p.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// format returns a scrypt-h64 string in MCF. Exported functions must
// verify params before calling format.
func format(digest, salt []byte, params *Parameter) Mcf {
	return Mcf(fmt.Sprintf(
		mcfString,
		params.Iterations,
		params.BlockSizeFactor,
		params.Parallelism,
		params.Length,
		params.SaltLength,
		hash64.StdEncoding.EncodeToString(salt),
		hash64.StdEncoding.EncodeToString(digest),
	))
}

// checkParameter verifies *Parameter for valid values and replaces
// unspecified zero values with values from PaperParams
func checkParameter(p *Parameter) error {
	// check and replace zero value
	if p.Iterations == 0 {
		p.Iterations = 14
	}

	// check and replace zero value
	if p.BlockSizeFactor == 0 {
		p.BlockSizeFactor = 8
	}

	// check and replace zero value
	if p.Parallelism == 0 {
		p.Parallelism = 1
	}

	// check and replace zero value
	if p.Length == 0 {
		p.Length = 32
	} else if p.Length < 16 {
		// scrypt-h64 specification defines minimum as 16 bytes
		return fmt.Errorf(
			"Minimum digest length is 16 bytes (requested: %d)",
			p.Length,
		)
	}

	// check and replace zero value
	if p.SaltLength == 0 {
		p.SaltLength = 16
	} else if p.SaltLength < 16 {
		// scrypt-h64 specification defines minimum as 16 bytes
		return fmt.Errorf(
			"Minimum salt length is 16 bytes (requested: %d)",
			p.SaltLength,
		)
	}
	return nil
}

// fieldsFromMCF slices an MCF string into its fields. It performs
// checks on field count and used identifier.
func fieldsFromMCF(mcf Mcf) ([]string, error) {
	var (
		sMcf   string
		fields []string
	)
	// remove leading and possibly trailing field separator
	sMcf = strings.Trim(string(mcf), `$`)

	// specification requires all 4 fields to be present
	fields = strings.Split(sMcf, `$`)
	if len(fields) != 4 {
		return nil, fmt.Errorf(
			"Invalid MCF string with only %d fields",
			len(fields),
		)
	}

	// check identifier
	if fields[fieldIdentifier] != `scrypt-h64` {
		return nil, fmt.Errorf(
			"Invalid MCF identifier: %s",
			fields[fieldIdentifier],
		)
	}
	return fields, nil
}

// parameterFromMCF assembles a Parameter struct from the values
// embedded inside the MCF string. Returns nil if a malformed or
// otherwise invalid string is passed in.
func parameterFromMCF(mcf Mcf) (*Parameter, error) {
	var (
		err                error
		num                uint64
		fields, parameters []string
	)
	p := &Parameter{}

	// slice MCF string
	if fields, err = fieldsFromMCF(mcf); err != nil {
		goto fail
	}

	// parse the parameters field using , as separator
	parameters = strings.Split(fields[fieldParameters], `,`)
	for _, param := range parameters {
		// individual parameters have to be of the form key=value
		opts := strings.Split(param, `=`)
		if len(opts) != 2 {
			err = fmt.Errorf("Invalid format for parameter: %s", param)
			goto fail
		}
		switch opts[0] {
		case "N":
			if num, err = strconv.ParseUint(opts[1], 10, 16); err != nil {
				goto fail
			}
			p.Iterations = uint16(num)
		case "r":
			if num, err = strconv.ParseUint(opts[1], 10, 8); err != nil {
				goto fail
			}
			p.BlockSizeFactor = uint8(num)
		case "p":
			if num, err = strconv.ParseUint(opts[1], 10, 8); err != nil {
				goto fail
			}
			p.Parallelism = uint8(num)
		case "l":
			if num, err = strconv.ParseUint(opts[1], 10, 16); err != nil {
				goto fail
			}
			p.Length = uint16(num)
		case "s":
			if num, err = strconv.ParseUint(opts[1], 10, 16); err != nil {
				goto fail
			}
			p.SaltLength = uint16(num)
		default:
			// unknown parameter
			err = fmt.Errorf("Unknown parameter: %s", opts[0])
			goto fail
		}
	}

	// check value range and fill in unspecified default values
	if err = checkParameter(p); err != nil {
		goto fail
	}

	return p, nil

fail:
	return nil, err
}

// saltFromMCF returns the salt field from an MCF string
func saltFromMCF(mcf Mcf) ([]byte, error) {
	var (
		err    error
		salt   []byte
		fields []string
	)

	// slice MCF string
	if fields, err = fieldsFromMCF(mcf); err != nil {
		return nil, err
	}

	if salt, err = hash64.StdEncoding.DecodeString(fields[fieldSalt]); err != nil {
		return nil, err
	}
	return salt, nil
}

// digestFromMCF returns the digest field of an MCF string
func digestFromMCF(mcf Mcf) ([]byte, error) {
	var (
		err    error
		digest []byte
		fields []string
	)

	// slice MCF string
	if fields, err = fieldsFromMCF(mcf); err != nil {
		return nil, err
	}

	if digest, err = hash64.StdEncoding.DecodeString(fields[fieldDigest]); err != nil {
		return nil, err
	}
	return digest, nil
}

// FromString takes a string, validates it and returns it as a custom
// typed Mcf string.
func FromString(s string) (res Mcf, err error) {
	// err+res declared as part of return types
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Not an scrypt-h64 string")
		}
	}()

	if _, err = parameterFromMCF(Mcf(s)); err != nil {
		goto fail
	}
	if _, err = saltFromMCF(Mcf(s)); err != nil {
		goto fail
	}
	if _, err = digestFromMCF(Mcf(s)); err != nil {
		goto fail
	}

	res = Mcf(s)
	return res, nil

fail:
	res = Mcf("")
	return res, err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
