package main

import (
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestParseGotoArgs_Success_WithTrailing(t *testing.T) {
    v, trailing, err := parseGotoArgs([]string{"3", "--no-compact"})
    assert.NoError(t, err)
    assert.Equal(t, 3, v)
    assert.Equal(t, []string{"--no-compact"}, trailing)
}

func TestParseGotoArgs_Errors(t *testing.T) {
    // Missing positional
    _, _, err := parseGotoArgs([]string{})
    assert.Error(t, err)

    // Non-integer
    _, _, err = parseGotoArgs([]string{"x"})
    assert.Error(t, err)

    // Out of range low
    _, _, err = parseGotoArgs([]string{"0"})
    assert.Error(t, err)

    // Out of range high
    _, _, err = parseGotoArgs([]string{"10"})
    assert.Error(t, err)
}

func TestParseMoveArgs_Success(t *testing.T) {
    // With flag and trailing
    v, all, trailing, err := parseMoveArgs([]string{"--all", "2", "--no-compact"})
    assert.NoError(t, err)
    assert.Equal(t, 2, v)
    assert.True(t, all)
    assert.Equal(t, []string{"--no-compact"}, trailing)

    // Minimal
    v, all, trailing, err = parseMoveArgs([]string{"2"})
    assert.NoError(t, err)
    assert.Equal(t, 2, v)
    assert.False(t, all)
    assert.Empty(t, trailing)
}

func TestParseMoveArgs_Errors(t *testing.T) {
    // No positional
    _, _, _, err := parseMoveArgs([]string{"--all"})
    assert.Error(t, err)

    // Non-integer positional
    _, _, _, err = parseMoveArgs([]string{"x"})
    assert.Error(t, err)

    // Out of range
    _, _, _, err = parseMoveArgs([]string{"0"})
    assert.Error(t, err)

    // Unknown flag causes parse error
    _, _, _, err = parseMoveArgs([]string{"--wat"})
    assert.Error(t, err)
}

func TestParseCycleArgs_Success_WithTrailing(t *testing.T) {
    dir, trailing, err := parseCycleArgs([]string{"next", "--no-compact"})
    assert.NoError(t, err)
    assert.Equal(t, "next", dir)
    assert.Equal(t, []string{"--no-compact"}, trailing)
}

func TestParseCycleArgs_Errors(t *testing.T) {
    // Missing positional
    _, _, err := parseCycleArgs([]string{})
    assert.Error(t, err)

    // Invalid value
    _, _, err = parseCycleArgs([]string{"left"})
    assert.Error(t, err)
}

func TestParseTrailingGlobalFlags(t *testing.T) {
    // Default compact true when no flags
    g, err := parseTrailingGlobalFlags([]string{})
    assert.NoError(t, err)
    assert.True(t, g.Compact)

    // --no-compact flips to false
    g, err = parseTrailingGlobalFlags([]string{"--no-compact"})
    assert.NoError(t, err)
    assert.False(t, g.Compact)

    // Unknown flag
    _, err = parseTrailingGlobalFlags([]string{"--wat"})
    assert.Error(t, err)

    // Unexpected extra args
    _, err = parseTrailingGlobalFlags([]string{"--no-compact", "extra"})
    assert.Error(t, err)
    if err != nil {
        assert.True(t, strings.Contains(err.Error(), "unexpected arguments"))
    }
}

func TestParseMoveArgs_FlagAfterPos_TreatedAsTrailing(t *testing.T) {
    v, all, trailing, err := parseMoveArgs([]string{"2", "--all", "--no-compact"})
    assert.NoError(t, err)
    assert.Equal(t, 2, v)
    // '--all' after positional should not be treated as subcommand flag
    assert.False(t, all)
    assert.Equal(t, []string{"--all", "--no-compact"}, trailing)
}
