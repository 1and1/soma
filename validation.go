package somatree

import (
	"unicode/utf8"

	"github.com/satori/go.uuid"
)

func specRepoCheck(spec RepositorySpec) bool {
	switch {
	case spec.Id == "":
		return false
	case spec.Name == "":
		return false
	case spec.Team == "":
		return false
	}
	l := utf8.RuneCountInString(spec.Name)
	switch {
	case l < 4:
		return false
	case l > 128:
		return false
	}
	if _, err := uuid.FromString(spec.Id); err != nil {
		return false
	}
	if _, err := uuid.FromString(spec.Team); err != nil {
		return false
	}
	return true
}

func specBucketCheck(spec BucketSpec) bool {
	switch {
	case spec.Id == "":
		return false
	case spec.Name == "":
		return false
	case spec.Environment == "":
		return false
	case spec.Team == "":
		return false
	case spec.Repository == "":
		return false
	}
	l := utf8.RuneCountInString(spec.Name)
	switch {
	case l < 4:
		return false
	case l > 512:
		return false
	}
	if _, err := uuid.FromString(spec.Id); err != nil {
		return false
	}
	if _, err := uuid.FromString(spec.Team); err != nil {
		return false
	}
	if _, err := uuid.FromString(spec.Repository); err != nil {
		return false
	}
	return true
}

func specGroupCheck(spec GroupSpec) bool {
	switch {
	case spec.Id == "":
		return false
	case spec.Name == "":
		return false
	case spec.Team == "":
		return false
	}
	l := utf8.RuneCountInString(spec.Name)
	switch {
	case l < 4:
		return false
	case l > 256:
		return false
	}
	if _, err := uuid.FromString(spec.Id); err != nil {
		return false
	}
	if _, err := uuid.FromString(spec.Team); err != nil {
		return false
	}
	return true
}

func specClusterCheck(spec ClusterSpec) bool {
	switch {
	case spec.Id == "":
		return false
	case spec.Name == "":
		return false
	case spec.Team == "":
		return false
	}
	l := utf8.RuneCountInString(spec.Name)
	switch {
	case l < 4:
		return false
	case l > 256:
		return false
	}
	if _, err := uuid.FromString(spec.Id); err != nil {
		return false
	}
	if _, err := uuid.FromString(spec.Team); err != nil {
		return false
	}
	return true
}

func specNodeCheck(spec NodeSpec) bool {
	switch {
	case spec.Id == "":
		return false
	case spec.Name == "":
		return false
	case spec.Team == "":
		return false
	case spec.ServerId == "":
		return false
	}
	l := utf8.RuneCountInString(spec.Name)
	switch {
	case l < 4:
		return false
	case l > 256:
		return false
	}
	if _, err := uuid.FromString(spec.Id); err != nil {
		return false
	}
	if _, err := uuid.FromString(spec.Team); err != nil {
		return false
	}
	if _, err := uuid.FromString(spec.ServerId); err != nil {
		return false
	}
	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
