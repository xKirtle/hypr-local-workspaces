package main

import (
	"fmt"
	"regexp"
	"strconv"
)

func zeroWidthToken(monitorID int) string {
	switch monitorID {
	case 0:
		return "\u200B\u200C"
	case 1:
		return "\u200D\u2060"
	default:
		return "\u200E\u200F"
	}
}

func GetZeroWidthName(monitorID, slot int) string {
	return strconv.Itoa(slot) + zeroWidthToken(monitorID)
}

var invisibleChars = regexp.MustCompile("[\u200B\u200C\u200D\u200E\u200F\u2060]")

func HasZeroWidthToken(slotName string) bool {
	return invisibleChars.MatchString(slotName)
}

func TrimZeroWidthToken(slotName string) string {
	return invisibleChars.ReplaceAllString(slotName, "")
}

// ParseLocalWorkspace strips any zero-widths and parses a positive int.
// Returns error if the name is empty, contains zero-width chars, or is non-numeric.
func ParseLocalWorkspace(name string) (int, error) {
	workspaceNameStr := TrimZeroWidthToken(name)
	if workspaceNameStr == "" || workspaceNameStr != name && HasZeroWidthToken(workspaceNameStr) {
		return 0, fmt.Errorf("invalid local workspace name (invisibles): %q", name)
	}

	workspaceName, err := strconv.Atoi(workspaceNameStr)
	if err != nil || workspaceName <= 0 {
		return 0, fmt.Errorf("invalid local workspace name (not a positive int): %q", name)
	}
	return workspaceName, nil
}
