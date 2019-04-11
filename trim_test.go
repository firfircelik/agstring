package agstring

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrimAndRemoveZeroWidthChars(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "   \u200B\u200BTest:\u200B:\u200B:Test", expected: "Test:::Test"},
		{input: "   \u200B\u200CTest:\u200C:\u200B:Test", expected: "Test:::Test"},
	}

	for _, testCase := range testCases {
		require.Equal(t, TrimAndRemoveZeroWidthChars(testCase.input), testCase.expected)
	}
}

func TestTrimSpace(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{nil, nil},
		{[]string{}, []string{}},
		{[]string{"a"}, []string{"a"}},
		{[]string{"a "}, []string{"a"}},
		{[]string{" a "}, []string{"a"}},
		{[]string{" a ", "b"}, []string{"a", "b"}},
		{[]string{" a", "b "}, []string{"a", "b"}},
	}
	for _, tt := range tests {
		require.ElementsMatch(t, tt.expected, TrimSpace(tt.input...))
	}
}

func TestTrimSuffixes(t *testing.T) {
	testCases := []struct {
		input, expected string
		suffixes        []string
	}{
		{input: "", expected: "", suffixes: nil},
		{input: "", expected: "", suffixes: []string{"a"}},
		{input: "   b   ", expected: "", suffixes: []string{"b"}},
		{input: "   c", expected: "", suffixes: []string{"c"}},
		{input: "   d\f", expected: "", suffixes: []string{"d"}},
		{input: "   d\fe\n", expected: "d\fe", suffixes: []string{"d"}},
	}

	for _, tt := range testCases {
		require.Equal(t, tt.expected, TrimSuffixes(tt.input, tt.suffixes...))
	}
}

func TestTrimPrefixesAndSpace(t *testing.T) {
	testCases := []struct {
		input    string
		prefixes []string
		expected string
	}{
		{" cif fob massive port", []string{"fob", "cif"}, "massive port"},
		{" \\cif fob massive port", []string{"fob", "\\cif"}, "massive port"},
		{" cif fobmassive port", []string{"fob", "cif"}, "fobmassive port"},
		{" cif fob massive port", []string{"big", "location"}, "cif fob massive port"},
		{" cif fob massive port", []string{""}, "cif fob massive port"},
		{" cif fob", []string{"fob", "cif"}, ""},
		{" ", []string{}, ""},
		{"", nil, ""},
	}

	for _, tt := range testCases {
		require.Equal(t, tt.expected, TrimPrefixesAndSpace(tt.input, tt.prefixes...))
	}
}

func TestTrimAllPrefixes(t *testing.T) {
	testCases := []struct {
		input    string
		prefixes []string
		expected string
	}{
		{" cif fob massive port", []string{"fob", "cif"}, "massive port"},
		{" \\ciffobmassive port", []string{"fob", "\\cif"}, "massive port"},
		{" cif fob massive port", []string{"fob"}, "cif fob massive port"},
		{" cif fob massive port", []string{"big", "location"}, "cif fob massive port"},
		{" cif fob massive port", []string{""}, "cif fob massive port"},
		{" cif fob", []string{"fob", "cif"}, ""},
		{" ", []string{}, ""},
		{"cif fob", nil, "cif fob"},
		{"", nil, ""},
	}

	for _, tt := range testCases {
		require.Equal(t, tt.expected, TrimAllPrefixes(tt.input, tt.prefixes...))
	}
}

func TestTrimAllSuffixes(t *testing.T) {
	testCases := []struct {
		input    string
		suffixes []string
		expected string
	}{
		{" cif fob massive port", []string{"massive", "port"}, "cif fob"},
		{" cif fobmassiveport", []string{"massive", "port"}, "cif fob"},
		{" cif fob massive port", []string{"massive"}, "cif fob massive port"},
		{" cif fob massive port", []string{"big", "location"}, "cif fob massive port"},
		{" cif fob massive port", []string{""}, "cif fob massive port"},
		{" cif fob", []string{"fob", "cif"}, ""},
		{" ", []string{}, ""},
		{"cif fob", nil, "cif fob"},
		{"", nil, ""},
	}

	for _, tt := range testCases {
		require.Equal(t, tt.expected, TrimAllSuffixes(tt.input, tt.suffixes...))
	}
}
