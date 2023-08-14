package lib

import "strings"

var ErrorMessage = []string{
	"error",
	"fd error",
}

func IsFdError(line string) bool {
	for _, msg := range ErrorMessage {
		if strings.Contains(line, msg) {
			return true
		}
	}
	return false
}
