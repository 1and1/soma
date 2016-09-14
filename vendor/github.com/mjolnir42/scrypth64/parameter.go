/*-
 * Copyright © 2016, Jörg Pernfuß <code.jpe@gmail.com>
 * All rights reserved.
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package scrypth64

// Parameter contains all required parameters for scrypth64
type Parameter struct {
	Iterations      uint16 // log2 of the iteration count
	BlockSizeFactor uint8  // scrypt r parameter
	Parallelism     uint8  // scrypt p parameter
	Length          uint16 // output length in bytes, dkLen parameter
	SaltLength      uint16 // salt length in bytes
}

// DefaultParams returns a prefilled *Parameter with default values.
// The defaults represent a 64MiByte memory consumption by scrypt per
// calculation and are intended for 2016 high performance datacenter
// equipment, not mobile gadgets.
// Uses a 512bit salt and 384bit digest.
func DefaultParams() *Parameter {
	return &Parameter{
		Iterations:      15,
		BlockSizeFactor: 16,
		Parallelism:     2,
		Length:          48,
		SaltLength:      64,
	}
}

// PaperParams returns a prefilled *Parameter with the recommended
// values from the scrypt paper. These represent 16MiByte of used
// memory per scrypt calculation.
// Additional values are chosen with popular values.
func PaperParams() *Parameter {
	return &Parameter{
		Iterations:      14,
		BlockSizeFactor: 8,
		Parallelism:     1,
		Length:          32,
		SaltLength:      16,
	}
}

// HighParams returns a prefilled *Parameter with the recommended high
// security parameters from the scrypt paper. These represent 1GiByte
// of used memory per scrypt calculation.
func HighParams() *Parameter {
	return &Parameter{
		Iterations:      20,
		BlockSizeFactor: 8,
		Parallelism:     1,
		Length:          32,
		SaltLength:      32,
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
