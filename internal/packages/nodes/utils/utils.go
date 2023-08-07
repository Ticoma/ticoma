package utils

// Utils package for /nodes/ module

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"ticoma/packages/nodes/interfaces"
)

// Strip a string from any whitespaces, tabs, newline chars, etc...
func StripString(str string, removeLastCharToo bool) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	if removeLastCharToo {
		str = strings.TrimSuffix(str, ",")
	}
	return str
}

// Conv ADPT to ADP (test util)
func StripPkgFromTimestamp(adpt *interfaces.ActionDataPackageTimestamped) *interfaces.ActionDataPackage {
	adp := &interfaces.ActionDataPackage{
		PlayerId:     adpt.PlayerId,
		PubKey:       adpt.PubKey,
		Position:     &interfaces.Position{X: adpt.Position.X, Y: adpt.Position.Y},
		DestPosition: &interfaces.DestPosition{X: adpt.DestPosition.X, Y: adpt.DestPosition.Y},
	}
	return adp
}

// Extract all the values of fields from a published Pkg
// Pkg types must be OK
// This needs some improvement
func ExtractValsFromStrPkg(pkg string) []string {
	fmt.Println("EXTRACT GOT ", pkg) // DEBUG

	// Order of ignored is important
	var ignored = []string{"{", "}", "\"", ":", "playerId", "pubKey", "posX", "posY", "pos", "destPosX", "destPosY", "destPos"}
	for i := 0; i < len(ignored); i++ {
		pkg = strings.ReplaceAll(pkg, ignored[i], "")
	}
	return strings.Split(pkg, ",")
}

// ADP / ADPT interface object -> String conversion
func StringifyPkg(pkg interface{}, trimLastChar bool) string {
	switch v := pkg.(type) {
	case interfaces.ActionDataPackage, interfaces.ActionDataPackageTimestamped:
		marshaled, err := json.Marshal(v)
		if err != nil {
			fmt.Println("[UTILS] Couldn't strngify ADP object. err: ", err)
		}
		return StripString(string(marshaled), trimLastChar)
	default:
		fmt.Println("[UTILS] Invalid function parameter.")
		return ""
	}
}

// Get all: Field names (nested), Json tags (optional) of a struct (can be empty)
//
// Returns (field names, field json tags, field types)
func GetInterfaceFieldData(obj interface{}, wantJsonTags bool, wantTypes bool) ([]string, []string, []interface{}) {
	var fieldNames []string
	var jsonTags []string
	var fieldTypes []interface{}
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			ft := t.Field(i).Type            // field type
			fn := t.Field(i).Name            // field name
			jt := t.Field(i).Tag.Get("json") // json tag

			// fmt.Println("FIELD ", fn, "JSON ", t.Field(i).Tag.Get("json"), "TYPE ", ft)
			fieldNames = append(fieldNames, fn)
			if wantJsonTags {
				jsonTags = append(jsonTags, jt)
			}

			// This part below could be wrapped to a smaller, recursive func to support more depth
			if ft.Kind() == reflect.Ptr {
				v := ft.Elem()
				if v.Kind() == reflect.Struct {
					for i := 0; i < v.NumField(); i++ {
						// fmt.Println("NESTED FIELD ", v.Field(i).Name, "JSON ", v.Field(i).Tag.Get("json"), "TYPE ", v.Field(i).Type)
						fieldNames = append(fieldNames, v.Field(i).Name)
						if wantJsonTags {
							jsonTags = append(jsonTags, v.Field(i).Tag.Get("json"))
						}
						if wantTypes {
							fieldTypes = append(fieldTypes, v.Field(i).Type)
						}
					}
				}
			} else if wantTypes {
				fieldTypes = append(fieldTypes, ft)
			}
		}
	} else {
		fmt.Println("[UTILS] Can't get fields - provided object is not a struct.")
	}
	return fieldNames, jsonTags, fieldTypes
}
