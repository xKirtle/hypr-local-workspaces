package main

import (
	"fmt"
	"strconv"
)

var zeroWidthDigits = []rune{
	'\u200b', // ZERO WIDTH SPACE
	'\u200c', // ZERO WIDTH NON-JOINER
	'\u200d', // ZERO WIDTH JOINER
	'\u200e', // LEFT-TO-RIGHT MARK
	'\u200f', // RIGHT-TO-LEFT MARK
	'\u2060', // WORD JOINER
	'\u2061', // FUNCTION APPLICATION
	'\u2062', // INVISIBLE TIMES
	'\u2063', // INVISIBLE SEPARATOR
	'\u2064', // INVISIBLE PLUS
}

// GetZeroWidthNameFromIndex generates a unique workspace name using zero-width characters based on the monitor ID and workspace index.
func GetZeroWidthNameFromIndex(monitorID, index int) (string, error) {
	if monitorID < 0 || monitorID > len(zeroWidthDigits)-1 {
		return "", fmt.Errorf("monitorID out of range: %d. A maximum of %d monitors is supported: ", monitorID, len(zeroWidthDigits)-1)
	}

	if index < 0 {
		return "", fmt.Errorf("index must be non-negative: %d", index)
	}

	workspaceName := strconv.Itoa(index + 1)

	// Prefix with monitor ID zero-width char to ensure uniqueness across monitors
	workspaceName += string(zeroWidthDigits[monitorID])

	// Append zero-width chars for each digit in the index
	indexAsStr := strconv.Itoa(index)
	for _, char := range indexAsStr {
		charAsNum, err := strconv.Atoi(string(char))

		// Should be impossible to hit this branch since we validated index above
		if err != nil || charAsNum < 0 || charAsNum > 9 {
			continue
		}

		workspaceName += string(zeroWidthDigits[charAsNum])
	}

	return workspaceName, nil
}

// GetZeroWidthNameToIndex extracts the workspace index from a zero-width named workspace.
func GetZeroWidthNameToIndex(name string) (int, error) {
	if name == "" {
		return -1, fmt.Errorf("empty workspace name")
	}

	// Traverse name until we hit a zero-width char
	lastDigitIndex := 0
	for lastDigitIndex < len(name) {
		char := name[lastDigitIndex]
		if char < '0' || char > '9' {
			break
		}

		lastDigitIndex++
	}

	if lastDigitIndex == 0 {
		return -1, fmt.Errorf("workspace name does not start with a digit: %q", name)
	}

	index, err := strconv.Atoi(name[:lastDigitIndex])

	// Number passed to Atoi should always be valid. Guard against overflow just in case.
	if err != nil {
		return -1, fmt.Errorf("parsing workspace name index: %w", err)
	}

	return index - 1, nil
}
