/* Apache v2 license
*  Copyright (C) <2019> Intel Corporation
*
*  SPDX-License-Identifier: Apache-2.0
 */

package odata

import (
	"net/url"
	"reflect"

	"github.com/intel/rsp-sw-toolkit-im-suite-go-odata/parser"
	"github.com/pkg/errors"
)

// ErrInvalidInput Client errors
var ErrInvalidInput = errors.New("odata syntax error")

type ODataQuery struct {
	Filter       *parser.ParseNode
	SelectFields *[]string
	Limit        *int
	Skip         *int
	SortFields   *map[string]int
}

func ParseODataFilter(filter string) (*parser.ParseNode, error) {
	filterQuery, err := parser.ParseFilterString(filter)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}
	return filterQuery, nil
}

func ParseODataFilterForMongo(filterQuery *parser.ParseNode) (map[string]interface{}, error) {
	filterMap, err := ApplyFilterForMongo(filterQuery)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}
	return filterMap, nil
}

func ParseODataSelect(odataSelect string) ([]string, error) {
	stringArray, err := parser.ParseStringArray(&odataSelect)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	selectSlice := reflect.ValueOf(stringArray)
	return ParseODataSelectSlice(selectSlice), nil
}

func ParseODataSelectSlice(selectSlice reflect.Value) []string {
	selectFields := make([]string, 0)
	if selectSlice.Len() > 1 && selectSlice.Index(0).Interface().(string) != "*" {
		for i := 0; i < selectSlice.Len(); i++ {
			fieldName := selectSlice.Index(i).Interface().(string)
			selectFields = append(selectFields, fieldName)
		}
	}

	return selectFields
}

func ParseODataOrderBy(odataOrderBy string) (map[string]int, error) {
	orderBySlice, err := parser.ParseOrderArray(&odataOrderBy)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}
	return ParseODataOrderBySlice(orderBySlice), nil
}

func ParseODataOrderBySlice(orderBySlice []parser.OrderItem) map[string]int {
	sortFields := make(map[string]int)
	for _, item := range orderBySlice {
		order := 1
		if item.Order == "desc" {
			order = -1
		}
		sortFields[item.Field] = order
	}
	return sortFields
}

// ParseODataURL creates a query based on odata parameters
func ParseODataURL(query url.Values) (*ODataQuery, error) {
	// Parse url values
	queryMap, err := parser.ParseURLValues(query)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	var limit *int
	if temp, ok := queryMap[parser.Top].(int); ok {
		limit = &temp
	}

	var skip *int
	if temp, ok := queryMap[parser.Skip].(int); ok {
		skip = &temp
	}

	var filter *parser.ParseNode
	if queryMap[parser.Filter] != nil {
		filter, _ = queryMap[parser.Filter].(*parser.ParseNode)
	}

	var selectFields *[]string
	if queryMap["$select"] != nil {
		selectSlice := reflect.ValueOf(queryMap["$select"])
		temp := ParseODataSelectSlice(selectSlice)
		selectFields = &temp
	}

	var sortFields *map[string]int
	if queryMap[parser.OrderBy] != nil {
		orderBySlice := queryMap[parser.OrderBy].([]parser.OrderItem)
		temp := ParseODataOrderBySlice(orderBySlice)
		sortFields = &temp
	}

	odataQuery := &ODataQuery{
		Filter:       filter,
		SelectFields: selectFields,
		Limit:        limit,
		Skip:         skip,
		SortFields:   sortFields,
	}

	return odataQuery, nil
}
