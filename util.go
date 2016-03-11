package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

func CalculateLookupId(id uint64, metric string) string {
	asset := strconv.FormatUint(id, 10)
	hash := sha256.New()
	hash.Write([]byte(asset))
	hash.Write([]byte(metric))

	return hex.EncodeToString(hash.Sum(nil))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
