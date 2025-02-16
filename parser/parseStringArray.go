/* Apache v2 license
*  Copyright (C) <2019> Intel Corporation
*
*  SPDX-License-Identifier: Apache-2.0
 */

package parser

import (
	"errors"
	"strings"

	"github.com/intel/rsp-sw-toolkit-im-suite-go-odata/parser/validatefield"
)

// OrderItem holds order key information
type OrderItem struct {
	Field string
	Order string
}

func ParseStringArray(value *string) ([]string, error) {
	result := strings.Split(*value, ",")

	// trim out space
	for idx, resultNoSpace := range result {
		result[idx] = strings.TrimSpace(resultNoSpace)
	}

	if len(result) == 0 {
		return nil, errors.New("cannot parse zero length string")
	}

	return result, nil
}

func ParseOrderArray(value *string) ([]OrderItem, error) {
	parsedArray, err := ParseStringArray(value)
	if err != nil {
		return nil, err
	}

	// Validate values for special characters
	valid := validatefield.New("~!@#$%^&*()_+-")
	for _, val := range parsedArray {
		if valid.ValidateField(val) || val == "" {
			return nil, errors.New("cannot support field " + val)
		}
	}

	result := make([]OrderItem, len(parsedArray))

	for i, v := range parsedArray {
		timmedString := strings.TrimSpace(v)
		compressedSpaces := strings.Join(strings.Fields(timmedString), " ")
		s := strings.Split(compressedSpaces, " ")

		if len(s) > 2 {
			return nil, errors.New("cannot have more than 2 items in orderby query")
		}

		if len(s) > 1 {
			if s[1] != "asc" &&
				s[1] != "desc" {
				return nil, errors.New("second value in orderby needs to be asc or desc")
			}
			result[i] = OrderItem{s[0], s[1]}
			continue
		}
		result[i] = OrderItem{compressedSpaces, "asc"}
	}
	return result, nil
}
