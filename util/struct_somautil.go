package util

import "log"

type SomaUtil struct {
	Log *log.Logger
}

func (u *SomaUtil) SetLog(l *log.Logger) {
	u.Log = l
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
