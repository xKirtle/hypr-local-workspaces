package main

import (
	"strconv"
	"testing"
)

func TestGetZeroWidthNameFromIndex(t *testing.T) {
	tests := []struct {
		monitorID      int
		index          int
		expectedOutput string
		shouldFail     bool
	}{
		{0, 0, "1\u200b\u200b", false},
		{1, 2, "3\u200c\u200d", false},
		{9, 15, "16\u2064\u200c\u2060", false},
		{10, 0, "", true},        // Unsupported monitorID
		{-1, 0, "", true},        // Invalid monitorID
		{0, -1, "0\u200b", true}, // Invalid index
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			output, err := GetZeroWidthNameFromIndex(test.monitorID, test.index)

			if test.shouldFail {
				if err == nil {
					t.Fatalf("expected failure but got success with output %q", output)
				}
				return // done; don’t look at output
			}

			if err != nil {
				t.Fatalf("expected success but got error: %v", err)
			}

			if output != test.expectedOutput {
				t.Fatalf("got %q, want %q", output, test.expectedOutput)
			}
		})
	}
}

func TestGetZeroWidthNameToIndex(t *testing.T) {
	tests := []struct {
		name           string
		expectedOutput int
		shouldFail     bool
	}{
		{"1\u200b\u200b", 0, false},
		{"3\u200c\u200d\u200b", 2, false},
		{"-1\u200b", -1, true},
		{"", -1, true},
		{"abc", -1, true},
		{"999999999999999999999999\u200b", -1, true},
		{"10\u2000", -1, true},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			output, err := GetZeroWidthNameToIndex(test.name)

			if test.shouldFail && err == nil {
				t.Fatalf("expected failure but got success with output: %q", output)
			}

			if !test.shouldFail && err != nil {
				t.Fatalf("expected success but got error: %v", err)
			}

			if !test.shouldFail && output != test.expectedOutput {
				t.Fatalf("got %q, want %q", output, test.expectedOutput)
			}
		})
	}
}
