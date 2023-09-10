package utils

// Utils package for /nodes/ module

import (
	"fmt"
	"reflect"
	"strings"
	"ticoma/internal/debug"
)

// Strip a string from any whitespaces, tabs, newlines and carriage returns
func StripString(str string, shouldRemoveLastChar bool) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	if shouldRemoveLastChar {
		str = strings.TrimSuffix(str, ",")
	}
	return str
}

// Extract values from a json stringified request (in-place)
func ExtractValsFromStrPkg(reqData string, ignoredStrings []string) []string {

	debug.DebugLog("[CACHE_UTILS] - Extracted values: "+reqData, debug.PLAYER)

	// Order of ignored is important
	var ignoredGenericJson = []string{"{", "}", "\"", ":"}
	ignored := append(ignoredGenericJson, ignoredStrings...)
	for i := 0; i < len(ignored); i++ {
		reqData = strings.ReplaceAll(reqData, ignored[i], "")
	}
	return strings.Split(reqData, ",")
}

// Get all: Field names (nested), Field types(optional), Json tags (optional) of a struct (can be empty)
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
