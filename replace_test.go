package agstring

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReplaceAll(t *testing.T) {
	testCasesWithASlice := []struct {
		input        string
		toBeReplaced string
		replacements []string
		expected     string
	}{
		{"January20190301June", "j", []string{"J"},
			"january20190301june"},
		{"My saster het a heth pet which as", "", []string{"et", "as"},
			"My ster h a hh p which "},
	}

	for _, tt := range testCasesWithASlice {
		require.Equal(t, tt.expected, ReplaceAll(tt.input, tt.toBeReplaced, tt.replacements...))
	}
}

func TestZeroWidthSpace(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "Test:\u200B:\u200B:Test", expected: "Test:::Test"},
		{input: "Test:\u200C:\u200B:Test", expected: "Test:::Test"},
		{input: "Test:\u200D:\u200B:Test", expected: "Test:::Test"},
		{input: "Test:\uFEFF:\u200B:Test", expected: "Test:::Test"},
	}

	for _, testCase := range testCases {
		require.Equal(t, ReplaceZeroWidthChars(testCase.input), testCase.expected)
	}
}

func TestReplaceMultispace(t *testing.T) {
	tests := []string{"   a   b   c d    ef  ", "     123   ", "1\n\n123"}
	expects := []string{"a b c d ef", "123", "1 123"}

	for i, tt := range tests {
		require.Equal(t, ReplaceMultispace(tt), expects[i])
	}
}

func TestReplaceWholeWords(t *testing.T) {
	tests := []struct {
		s           string
		old         string
		replacement string
		expected    string
	}{
		{"port saint jose", "nt", "nntt", "port saint jose"},
		{"po saint jose", "po", "port", "port saint jose"},
		{"port st jose", "st", "saint", "port saint jose"},
		{"port saint hose", "hose", "jose", "port saint jose"},
	}
	for i, tt := range tests {
		if res := ReplaceWholeWord(tt.s, tt.old, tt.replacement); res != tt.expected {
			t.Errorf("%d: For input %q replacing %s to %s, expected %q, got %q",
				i, tt.s, tt.old, tt.replacement, tt.expected, res)
		}
	}
}

func TestDayOrdinalReplacer(t *testing.T) {
	testCases := []struct {
		input       string
		replaceWith string
		expected    string
	}{
		{
			"2nd Jan",
			"",
			"2 Jan",
		}, {
			"3rd Jan",
			"XX",
			"3XX Jan",
		}, {
			"4th Jan",
			"ZZZ",
			"4ZZZ Jan",
		},
	}

	for _, tt := range testCases {
		require.Equal(t, tt.expected, ReplaceDayOrdinal(tt.input, tt.replaceWith))
	}

	require.Equal(t, "1 Jan", ReplaceDayOrdinal("1st Jan"))
}

func TestReplaceNewline(t *testing.T) {
	input := "Hi\nThere"
	testCases := []struct {
		replaceWith string
		expected    string
	}{
		{
			"",
			"HiThere",
		}, {
			" ",
			"Hi There",
		}, {
			"XXX",
			"HiXXXThere",
		},
	}

	for _, tt := range testCases {
		require.Equal(t, tt.expected, ReplaceNewline(input, tt.replaceWith))
	}

	require.Equal(t, "HiThere", ReplaceNewline(input))
}
