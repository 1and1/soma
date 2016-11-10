/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package adm

import (
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"unicode/utf8"

	"github.com/1and1/soma/lib/proto"
	resty "gopkg.in/resty.v0"
)

// ValidateRuneCount tests if a string's number of unicode runes is
// below an upper limit
func ValidateRuneCount(s string, l int) error {
	if cnt := utf8.RuneCountInString(s); cnt > l {
		return fmt.Errorf("Validation error, string '%s' above character limit %d.",
			s, l)
	}
	return nil
}

// ValidateRuneCountRange tests if a string's number of unicode runes is
// between an upper and lower bound
func ValidateRuneCountRange(s string, lower, higher int) error {
	if utf8.RuneCountInString(s) < lower || utf8.RuneCountInString(s) > higher {
		return fmt.Errorf("Validation error, string '%s' outside permitted length."+
			"Required: %d < len(%s) < %d.", s, lower, s, higher)
	}
	return nil
}

// IsUUID validates if a string is one very narrow formatting of a UUID,
// namely the one used by the server. Other valid formats with braces etc
// are not accepted
func IsUUID(s string) bool {
	const reUUID string = `^[[:xdigit:]]{8}-[[:xdigit:]]{4}-[1-5][[:xdigit:]]{3}-[[:xdigit:]]{4}-[[:xdigit:]]{12}$`
	const reUNIL string = `^0{8}-0{4}-0{4}-0{4}-0{12}$`
	re := regexp.MustCompile(fmt.Sprintf("%s|%s", reUUID, reUNIL))

	return re.MatchString(s)
}

// isUint64 validates if a string is a uint64 number
func isUint64(s string) (bool, uint64) {
	if u, err := strconv.ParseUint(s, 10, 64); err != nil {
		return false, 0
	} else {
		return true, u
	}
}

// ValidateOncallNumber tests if a string is a 4 digit number
func ValidateOncallNumber(n string) error {
	num, err := strconv.Atoi(n)
	if err != nil {
		return fmt.Errorf("Syntax error, argument is not a number: %s", err.Error())
	}

	if num <= 0 || num > 9999 {
		return fmt.Errorf("Phone number must be 4-digit extension")
	}
	return nil
}

// ValidateEmployeeNumber tests if a string is a non-negative number
func ValidateEmployeeNumber(s string) error {
	employeeNumber, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("Error: employeenr argument not a number")
	}
	if employeeNumber < 0 {
		return fmt.Errorf("Negative employee number is not allowed")
	}
	return nil
}

// ValidateMailAddress tests if a string is a mail address
func ValidateMailAddress(s string) error {
	if _, err := mail.ParseAddress(s); err != nil {
		return fmt.Errorf("Error: mailaddr does not parse as RFC 5322 address")
	}
	return nil
}

// ValidateBool tests that string s represents and bool and sets
// its value in b.
// Returns error if s could not be parsed as bool.
func ValidateBool(s string, b *bool) error {
	var err error
	*b, err = strconv.ParseBool(s)
	return err
}

// ValidateLBoundUint64 tests that string s is an unsigned number
// with a specific minimum value and sets it in i.
// Returns error if the number conversion failed or the minimum
// value was not reached.
func ValidateLBoundUint64(s string, i *uint64, min uint64) error {
	var err error
	if *i, err = strconv.ParseUint(s, 10, 64); err != nil {
		return err
	}
	if *i < min {
		return fmt.Errorf("Error, value %d is less than the minimun %d", i, min)
	}
	return nil
}

// ValidateUnit tests against the server if string s is a valid
// unit.
func ValidateUnit(s string) error {
	res, err := fetchObjList(`/units/`)
	if err != nil {
		return err
	}

	if res.Units != nil {
		for _, unit := range *res.Units {
			if unit.Unit == s {
				return nil
			}
		}
	}
	return fmt.Errorf("Value %s is not a valid unit", s)
}

// ValidateProvider tests aginst the server if string s is a valid
// provider.
func ValidateProvider(s string) error {
	res, err := fetchObjList(`/providers/`)
	if err != nil {
		return err
	}

	if res.Providers != nil {
		for _, prov := range *res.Providers {
			if prov.Name == s {
				return nil
			}
		}
	}
	return fmt.Errorf("Value %s is not a valid provider", s)
}

// ValidateSystemProperty tests against the server if string s is a
// valid system property.
func ValidateSystemProperty(s string) error {
	res, err := fetchObjList(`/property/system/`)
	if err != nil {
		return err
	}

	if res.Properties != nil {
		for _, prop := range *res.Properties {
			if prop.System.Name == s {
				return nil
			}
		}
	}
	return fmt.Errorf("Invalid system property requested: %s", s)
}

// ValidateEnvironment tests against the server if string s is a
// valid system property
func ValidateEnvironment(s string) error {
	res, err := fetchObjList(`/environments/`)
	if err != nil {
		return err
	}

	if res.Environments != nil {
		for _, env := range *res.Environments {
			if env.Name == s {
				return nil
			}
		}
	}
	return fmt.Errorf("Invalid environment requested: %s", s)
}

// ValidateStatus tests against the server if string s is a
// valid check instance status
func ValidateStatus(s string) error {
	res, err := fetchObjList(`/status/`)
	if err != nil {
		return err
	}

	if res.Status != nil {
		for _, st := range *res.Status {
			if st.Name == s {
				return nil
			}
		}
	}
	return fmt.Errorf("Invalid instance status requested: %s", s)
}

// ValidateInstance tests against the server if string s is a
// valid check instance id
func ValidateInstance(s string) error {
	if !IsUUID(s) {
		return fmt.Errorf("Argument is not a UUID: %s", s)
	}
	path := fmt.Sprintf("/instances/%s", s)

	if _, err := fetchObjList(path); err != nil {
		return err
	}
	// instanceId is valid if there was no 404/NotFound
	return nil
}

// ValidatePredicate tests against the server if s is a valid
// predicate
func ValidatePredicate(s string) error {
	res, err := fetchObjList(fmt.Sprintf("/predicates/%s", s))
	if err != nil {
		return err
	}

	if res.Predicates != nil || len(*res.Predicates) == 0 {
		return fmt.Errorf(`no object returned`)
	}

	if s == (*res.Predicates)[0].Symbol {
		return nil
	}
	return fmt.Errorf("Invalid predicate requested: %s", s)
}

// ValidateCategory tests against the server if s is a valid
// category
func ValidateCategory(s string) error {
	res, err := fetchObjList(fmt.Sprintf("/category/%s", s))
	if err != nil {
		return err
	}

	if res.Categories != nil || len(*res.Categories) == 0 {
		return fmt.Errorf(`no object returned`)
	}

	if s == (*res.Categories)[0].Name {
		return nil
	}
	return fmt.Errorf("Invalid category requested: %s", s)
}

// ValidateCheckConstraints tests that all specified check constraints
// resolve to a valid property or attribute.
func ValidateCheckConstraints(repoId, teamId string,
	constraints []proto.CheckConfigConstraint) (
	[]proto.CheckConfigConstraint, error) {
	valid := []proto.CheckConfigConstraint{}

	for _, prop := range constraints {
		switch prop.ConstraintType {
		case `native`:
			if _, err := fetchObjList(
				fmt.Sprintf("/property/native/%s", prop.Native.Name),
			); err != nil {
				return nil, err
			}
			valid = append(valid, prop)

		case `system`:
			if _, err := fetchObjList(
				fmt.Sprintf("/property/system/%s", prop.System.Name),
			); err != nil {
				return nil, err
			}
			valid = append(valid, prop)

		case `attribute`:
			if _, err := fetchObjList(
				fmt.Sprintf("/attributes/%s", prop.Attribute.Name),
			); err != nil {
				return nil, err
			}
			valid = append(valid, prop)

		case `oncall`:
			oncall := proto.PropertyOncall{}
			var err error
			if prop.Oncall.Name != `` {
				if oncall.Id, err = LookupOncallId(
					prop.Oncall.Name); err != nil {
					return nil, err
				}
			} else if prop.Oncall.Id != `` {
				if oncall.Id, err = LookupOncallId(
					prop.Oncall.Id); err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf(
					`Invalid oncall constraint spec`)
			}
			valid = append(valid, proto.CheckConfigConstraint{
				ConstraintType: prop.ConstraintType,
				Oncall:         &oncall,
			})

		case `service`:
			service := proto.PropertyService{}
			var err error
			if service.Name, err = LookupServicePropertyId(
				prop.Service.Name, teamId); err != nil {
				return nil, err
			}
			service.TeamId = teamId
			valid = append(valid, proto.CheckConfigConstraint{
				ConstraintType: prop.ConstraintType,
				Service:        &service,
			})

		case `custom`:
			custom := proto.PropertyCustom{}
			var err error
			if custom.Id, err = LookupCustomPropertyId(
				prop.Custom.Name, repoId); err != nil {
				return nil, err
			}
			custom.RepositoryId = repoId
			custom.Value = prop.Custom.Value
			valid = append(valid, proto.CheckConfigConstraint{
				ConstraintType: prop.ConstraintType,
				Custom:         &custom,
			})
		}
	}
	return valid, nil
}

// ValidateThresholds tests that all thresholds use the same predicate
// and all referenced levels exist. It normalizes the possible mixed
// use of long and short level names.
func ValidateThresholds(thresholds []proto.CheckConfigThreshold) (
	[]proto.CheckConfigThreshold, error) {
	valid := []proto.CheckConfigThreshold{}
	pred := ``

	for _, thr := range thresholds {
		var err error
		if pred == `` {
			pred = thr.Predicate.Symbol
		} else {
			if pred != thr.Predicate.Symbol {
				return nil, fmt.Errorf(
					"Error, threshold specification is using"+
						" multiple predicates: %s, %s",
					pred,
					thr.Predicate.Symbol,
				)
			}
		}
		t := proto.CheckConfigThreshold{
			Value:     thr.Value,
			Predicate: proto.Predicate{},
			Level:     proto.Level{},
		}
		if t.Level.Name, err = LookupLevelName(
			thr.Level.Name); err != nil {
			return nil, err
		}
		if err = ValidatePredicate(thr.Predicate.Symbol); err != nil {
			return nil, err
		}
		t.Predicate.Symbol = thr.Predicate.Symbol
		valid = append(valid, t)
	}
	return valid, nil
}

// fetchObjList is a helper for ValidateUnit and ValidateProvider
func fetchObjList(path string) (*proto.Result, error) {
	var (
		err  error
		resp *resty.Response
		res  *proto.Result
	)
	if resp, err = GetReq(path); err != nil {
		return nil, err
	}
	if res, err = decodeResponse(resp); err != nil {
		return nil, err
	}
	if err = checkApplicationError(res); err != nil {
		return nil, err
	}
	return res, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
