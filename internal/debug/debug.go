package debug

import (
	"fmt"
	"os"
	"strconv"
)

func getDebugVal() int {
	dm, err := strconv.Atoi(os.Getenv("DEBUG"))
	if dm != 0 && err != nil {
		panic("[SETUP] - Couldn't parse debug mode value from .env")
	}
	return dm
}

type Source int8

const (
	NO_LOGS Source = iota
	ALL_LOGS
	PLAYER
	NETWORK
)

//	DEBUG MODES
//
// e.g: .env conf: DEBUG=0
//
// 0 = no logs at all
//
// 1 = enable all logs
//
// 2 = only logs from player module
//
// 3 = only logs from network module
func DebugLog(msg string, source Source) {

	dm := getDebugVal()
	switch dm {
	case 0:
		return
	case 1:
		fmt.Println("[DEBUG]", msg)
	case 2:
		if source != PLAYER {
			return
		}
		fmt.Println("[DEBUG]", msg)
	case 3:
		if source != NETWORK {
			return
		}
		fmt.Println("[DEBUG]", msg)
	default:
		return
	}
}
