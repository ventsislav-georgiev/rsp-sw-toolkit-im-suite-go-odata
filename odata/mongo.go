package odata

import (
	"encoding/hex"
	"strings"

	"github.com/intel/rsp-sw-toolkit-im-suite-go-odata/parser"
)

type MongoRegex struct {
	Pattern string `json:"$regex"`
	Options string `json:"$options"`
}

func ApplyFilterForMongo(node *parser.ParseNode) (map[string]interface{}, error) {
	filter := make(map[string]interface{})

	if _, ok := node.Token.Value.(string); ok {
		switch node.Token.Value {

		case "eq":
			// Escape single quotes in the case of strings
			if _, valueOk := node.Children[1].Token.Value.(string); valueOk {
				node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)
			}
			value := map[string]interface{}{"$" + node.Token.Value.(string): node.Children[1].Token.Value}
			if _, keyOk := node.Children[0].Token.Value.(string); !keyOk {
				return nil, ErrInvalidInput
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "ne":
			// Escape single quotes in the case of strings
			if _, valueOk := node.Children[1].Token.Value.(string); valueOk {
				node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)
			}
			value := map[string]interface{}{"$" + node.Token.Value.(string): node.Children[1].Token.Value}
			if _, keyOk := node.Children[0].Token.Value.(string); !keyOk {
				return nil, ErrInvalidInput
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "gt":
			var keyString string
			if keyString, ok = node.Children[0].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}

			var value map[string]interface{}
			if keyString == "_id" {
				var idString string
				if _, ok := node.Children[1].Token.Value.(string); ok {
					idString = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)
				}
				decodedString, err := hex.DecodeString(idString)
				if err != nil || len(decodedString) != 12 {
					return nil, ErrInvalidInput
				}
				value = map[string]interface{}{"$" + node.Token.Value.(string): decodedString}
			} else {
				value = map[string]interface{}{"$" + node.Token.Value.(string): node.Children[1].Token.Value}
			}
			filter[keyString] = value

		case "ge":
			value := map[string]interface{}{"$gte": node.Children[1].Token.Value}
			if _, ok := node.Children[0].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "lt":
			value := map[string]interface{}{"$" + node.Token.Value.(string): node.Children[1].Token.Value}
			if _, ok := node.Children[0].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "le":
			value := map[string]interface{}{"$lte": node.Children[1].Token.Value}
			if _, ok := node.Children[0].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "and":
			leftFilter, err := ApplyFilterForMongo(node.Children[0]) // Left children
			if err != nil {
				return nil, err
			}
			rightFilter, _ := ApplyFilterForMongo(node.Children[1]) // Right children
			if err != nil {
				return nil, err
			}
			filter["$and"] = []map[string]interface{}{leftFilter, rightFilter}

		case "or":
			leftFilter, err := ApplyFilterForMongo(node.Children[0]) // Left children
			if err != nil {
				return nil, err
			}
			rightFilter, err := ApplyFilterForMongo(node.Children[1]) // Right children
			if err != nil {
				return nil, err
			}
			filter["$or"] = []map[string]interface{}{leftFilter, rightFilter}

		//Functions
		case "startswith":
			if _, ok := node.Children[1].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)
			filter[node.Children[0].Token.Value.(string)] = MongoRegex{"^" + node.Children[1].Token.Value.(string), "g"}

		case "endswith":
			if _, ok := node.Children[1].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)
			filter[node.Children[0].Token.Value.(string)] = MongoRegex{"^.*" + node.Children[1].Token.Value.(string) + "$", "g"}

		case "contains":
			if _, ok := node.Children[1].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)
			filter[node.Children[0].Token.Value.(string)] = MongoRegex{"^.*" + node.Children[1].Token.Value.(string), "g"}

		}
	}

	return filter, nil
}
