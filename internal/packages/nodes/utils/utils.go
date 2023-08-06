package utils

// Utils package for /nodes/ module

import (
	"encoding/json"
	"fmt"
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
