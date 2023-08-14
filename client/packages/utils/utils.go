package utils

import (
	"fmt"
	"math/rand"
	"os/exec"
)

// Get current hash as string
func GetCommitHash() string {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		fmt.Println(err)
	}
	commitHash := string(out)
	return commitHash
}

// Gen random number in range - (inclusive, exclusive)
func RandRange(min int, max int) int {
	return rand.Intn(max-min) + min
}
