package util

import (
	"log"
	"net/url"
)

type SomaUtil struct {
	Log           *log.Logger
	ApiUrl        *url.URL
	PropertyTypes []string
	Views         []string
}

func (u *SomaUtil) SetLog(l *log.Logger) {
	u.Log = l
}

func (u *SomaUtil) SetUrl(str string) {
	url, err := url.Parse(str)
	if err != nil {
		u.Log.Printf("Error parsing API address from config file")
		u.Log.Fatal(err)
	}
	u.ApiUrl = url
}

func (u *SomaUtil) SetPropertyTypes(sl []string) {
	u.PropertyTypes = sl
}

func (u *SomaUtil) SetViews(sl []string) {
	u.Views = sl
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
