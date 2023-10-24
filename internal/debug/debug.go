package debug

import (
	"fmt"
	"os"
	"strconv"
)

func getDebugVal() int {
	dmStr := os.Getenv("DEBUG")
	if dmStr == "" {
		return 0
	}
	dm, err := strconv.Atoi(dmStr)
	if err != nil {
		fmt.Println("[WARNING] - Couldn't parse debug mode value from .env, defaulting to 0 (no logs)")
		return 0
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
