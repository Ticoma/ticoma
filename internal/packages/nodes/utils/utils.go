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

// Extract all the values of fields from a published Pkg
// Pkg types must be OK
// This needs some improvement
func ExtractValsFromStrPkg(pkg string) []string {
	fmt.Println("EXTRACT GOT ", pkg)
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
// Returns (names, jsons)
func GetInterfaceFieldNames(obj interface{}, wantJsonTags bool) ([]string, []string) {
	var fieldNames []string
	var jsonTags []string
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			ft := t.Field(i).Type            // field type
			fn := t.Field(i).Name            // field name
			jt := t.Field(i).Tag.Get("json") // json tag

			fmt.Println("FIELD ", fn, "JSON ", t.Field(i).Tag.Get("json"))

			fieldNames = append(fieldNames, fn)
			if wantJsonTags {
				jsonTags = append(jsonTags, jt)
			}

			// This part below could be wrapped to a smaller, recursive func to support more depth
			if ft.Kind() == reflect.Ptr {
				v := ft.Elem()
				if v.Kind() == reflect.Struct {
					for i := 0; i < v.NumField(); i++ {
						fmt.Println("NESTED FIELD ", v.Field(i).Name, "JSON ", v.Field(i).Tag.Get("json"))
						fieldNames = append(fieldNames, v.Field(i).Name)
						if wantJsonTags {
							jsonTags = append(jsonTags, v.Field(i).Tag.Get("json"))
						}
					}
				}
			}

		}
	} else {
		fmt.Println("not a stuct")
	}

	return fieldNames, jsonTags
}

func ConstructPkg(model interface{}, vals interface{}, timestamped bool) {

}
