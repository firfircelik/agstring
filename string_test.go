package agstring

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReplaceMultispace(t *testing.T) {
	tests := []string{"   a   b   c d    ef  ", "     123   ", "1\n\n123"}
	expects := []string{"a b c d ef", "123", "1 123"}

	for i, tt := range tests {
		require.Equal(t, ReplaceMultispace(tt), expects[i])
	}
}

func TestFirst(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{nil, ""},
		{[]string{}, ""},
		{[]string{"one"}, "one"},
		{[]string{"one", "two"}, "one"},
	}

	for _, tt := range tests {
		require.Equal(t, tt.expected, First(tt.input...))
	}
}

func TestLast(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{nil, ""},
		{[]string{}, ""},
		{[]string{"one"}, "one"},
		{[]string{"one", "two"}, "two"},
	}

	for _, tt := range tests {
		require.Equal(t, tt.expected, Last(tt.input...))
	}
}

func TestNonEmpty(t *testing.T) {
	tests := []struct {
		input, expected, ifTrimmed []string
	}{
		{nil, nil, nil},
		{[]string{}, nil, nil},
		{[]string{"one", "", "   "}, []string{"one", "   "}, []string{"one"}},
		{[]string{"one", "", "two", "  \n"}, []string{"one", "two", "  \n"}, []string{"one", "two"}},
	}

	for _, tt := range tests {
		require.ElementsMatch(t, tt.expected, NonEmpty(tt.input...))
		require.ElementsMatch(t, tt.ifTrimmed, NonEmptyIfTrimmed(tt.input...))
	}
}

func TestRemoveDiacritics(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ąćęłńóśźż ĄĆĘŁŃÓŚŹŻ", "acelnoszz ACELNOSZZ"},                         // Polish
		{"čďěňřšťž ČĎĚŇŘŠŤŽ", "cdenrstz CDENRSTZ"},                             // Czech
		{"éàèùâêîôûëïöüÿç ÉÀÈÙÂÊÎÔÛËÏÖÜŸÇ", "eaeuaeioueiouyc EAEUAEIOUEIOUYC"}, // French
		{"ğöüçş ĞÖÜÇŞ", "goucs GOUCS"},                                         // Turkish
	}

	for _, tt := range tests {
		require.Equal(t, tt.expected, RemoveDiacritics(tt.input))
	}
}

func TestEmptyIf(t *testing.T) {
	tests := []struct {
		s         string
		emptyList []string
		expected  string
	}{
		{"", nil, ""},
		{"agflow", nil, "agflow"},
		{"agflow", []string{}, "agflow"},
		{"-", []string{"-"}, ""},
		{"?", []string{"-", "tbd"}, "?"},
		{"agflow", []string{""}, "agflow"},
	}
	for _, tt := range tests {
		require.Equal(t, tt.expected, EmptyIf(tt.s, tt.emptyList...))
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

func TestTitle(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"abc abc", "Abc Abc"},
		{"ABC ABC", "Abc Abc"},
		{"Abc Abc", "Abc Abc"},
		{"ABc aBC", "Abc Abc"},
	}
	for _, tt := range tests {
		require.Equal(t, tt.expected, Title(tt.input))
	}
}

func TestSafeAtoi(t *testing.T) {
	// Given
	s1 := "100"
	// When
	i1, err1 := SafeAtoi(s1)
	// Then
	require.Equal(t, 100, i1)
	require.Nil(t, err1)

	// Given
	s2 := ""
	// When
	i2, err2 := SafeAtoi(s2)
	// Then
	require.Equal(t, 0, i2)
	require.Nil(t, err2)

	// Given
	s3 := "hello"
	// When
	_, err3 := SafeAtoi(s3)
	// Then
	require.NotNil(t, err3)
}

func TestContainsAll(t *testing.T) {
	testCases := []struct {
		holder, searched []string
		expected         bool
	}{
		{nil, nil, true},
		{[]string{}, []string{}, true},
		{nil, []string{"a"}, false},
		{[]string{"a"}, []string{"a"}, true},
		{[]string{"a", "b"}, []string{"a"}, true},
		{[]string{"a", "b"}, []string{"a", "c"}, false},
	}

	for _, tt := range testCases {
		require.Equal(t, tt.expected, ContainsAll(tt.holder, tt.searched...))
	}
}

func TestStringContainsAll(t *testing.T) {
	testCases := []struct {
		holder   string
		searched []string
		expected bool
	}{
		{"", nil, true},
		{"", []string{}, true},
		{"", []string{"a"}, false},
		{"a", []string{"a"}, true},
		{"ab", []string{"a"}, true},
		{"ab", []string{"a", "c"}, false},
	}

	for _, tt := range testCases {
		require.Equal(t, tt.expected, StringContainsAll(tt.holder, tt.searched...))
	}
}
func TestContainsAny(t *testing.T) {
	testCases := []struct {
		holder, searched []string
		expected         bool
	}{
		{nil, []string{"a"}, false},
		{[]string{"a"}, []string{"a"}, true},
		{[]string{"a", "b"}, []string{"a"}, true},
		{[]string{"a", "b"}, []string{"a", "c"}, true},
		{[]string{"a", "b"}, []string{"d", "e"}, false},
	}

	for _, tt := range testCases {
		require.Equal(t, tt.expected, ContainsAny(tt.holder, tt.searched...))
	}
}

func TestStringContainsAny(t *testing.T) {
	testCases := []struct {
		holder   string
		searched []string
		expected bool
	}{
		{"", []string{"a"}, false},
		{"a", []string{"a"}, true},
		{"ab", []string{"a"}, true},
		{"ab", []string{"a", "c"}, true},
		{"ab", []string{"c", "d"}, false},
	}

	for _, tt := range testCases {
		require.Equal(t, tt.expected, StringContainsAny(tt.holder, tt.searched...))
	}
}

func TestRegexpGroups(t *testing.T) {
	testCases := []struct {
		re       *regexp.Regexp
		input    string
		matches  bool
		expected map[string]string
	}{
		{
			re:       regexp.MustCompile(`^(?P<month>[a-zA-Z]+?)\s*(?P<year>\d\d(?:\d\d)?)?$`),
			input:    "Mrch19",
			matches:  true,
			expected: map[string]string{"month": "Mrch", "year": "19"},
		},
		{
			re:      regexp.MustCompile(`^(?P<month>[a-zA-Z]+?)\s*(?P<year>\d\d(?:\d\d)?)?$`),
			input:   "Mrch199",
			matches: false,
		},
		{
			re:       regexp.MustCompile(`^([a-zA-Z]+?)\s*(?P<year>\d\d(?:\d\d)?)?$`),
			input:    "Mrch19",
			matches:  true,
			expected: map[string]string{"year": "19"},
		},
	}
	for _, tt := range testCases {
		match, ok := RegexpGroups(tt.re, tt.input)
		require.Equal(t, tt.matches, ok, "for input %s", tt.input)
		if tt.matches {
			require.Equal(t, tt.expected, match)
		}
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

func TestTakeTo(t *testing.T) {
	testCases := []struct {
		input        []string
		truncateTill int
		expected     []string
	}{
		{
			[]string{"ąćęłńóśźż", "ĄĆĘŁŃÓŚŹŻ"},
			5,
			[]string{"ąćęłń", "ĄĆĘŁŃ"},
		},
		{
			[]string{"ドンキーコング", "hello"},
			100,
			[]string{"ドンキーコング", "hello"},
		},
		{
			[]string{"ドンキーコング", "hello"},
			4,
			[]string{"ドンキー", "hell"},
		},
		{
			[]string{"ドンキーコング", "hello"},
			7,
			[]string{"ドンキーコング", "hello"},
		},
		{
			[]string{"1234567890", "1234567890"},
			7,
			[]string{"1234567", "1234567"},
		},
	}

	for _, tt := range testCases {
		require.Equal(t, tt.expected, TakeTo(tt.input, tt.truncateTill))
	}
}

func TestTakeFrom(t *testing.T) {
	testCases := []struct {
		input        []string
		truncateFrom int
		expected     []string
	}{
		{
			[]string{"ąćęłńóśźż", "ĄĆĘŁŃÓŚŹŻ"},
			5,
			[]string{"óśźż", "ÓŚŹŻ"},
		},
		{
			[]string{"ドンキーコング", "hello"},
			100,
			[]string{"", ""},
		},
		{
			[]string{"ドンキーコング", "hello"},
			4,
			[]string{"コング", "o"},
		},
		{
			[]string{"ドンキーコング", "hello"},
			7,
			[]string{"", ""},
		},
		{
			[]string{"1234567890", "1234567890"},
			7,
			[]string{"890", "890"},
		},
	}

	for _, tt := range testCases {
		require.Equal(t, tt.expected, TakeFrom(tt.input, tt.truncateFrom))
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

func TestMap(t *testing.T) {
	f1 := func(s string) string { return strings.ToTitle(s) }
	f2 := func(s string) string { return strings.Repeat(s, 3) }

	testCases := []struct {
		input    []string
		funcs    []func(string) string
		expected []string
	}{
		{[]string{"hello"}, []func(string) string{f1}, []string{"HELLO"}},
		{[]string{"hello"}, []func(string) string{f1, f2}, []string{"HELLOHELLOHELLO"}},
	}

	for _, tt := range testCases {
		require.Equal(t, tt.expected, Map(tt.input, tt.funcs...))
	}
}

func TestHasSuffix(t *testing.T) {
	testCases := []struct {
		input    string
		suffixes []string
		expected bool
	}{
		{" cif fob massive port", []string{"port", "market"}, true},
		{" cif fob massive port", []string{""}, true},
		{"cif fob massive port", []string{"rt"}, true},
		{"cif fob massive port", []string{"cif", "massive"}, false},
	}
	for _, tt := range testCases {
		require.Equal(t, tt.expected, HasSuffix(tt.input, tt.suffixes...))
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

func TestConcat(t *testing.T) {
	testCasesWithASlice := []struct {
		input    []string
		expected []string
	}{
		{[]string{"slice1", "slice2"}, []string{"slice1", "slice2"}},
		{[]string{"Rome", "June"}, []string{"Rome", "June"}},
	}

	for _, tt := range testCasesWithASlice {
		require.Equal(t, tt.expected, Concat(tt.input))
	}

	testCasesForSlices := []struct {
		input    [][]string
		expected []string
	}{
		{[][]string{{"slice1", "slice2"}, {"slice3", "slice4"}},
			[]string{"slice1", "slice2", "slice3", "slice4"}},
		{[][]string{{"slice1", "slice2"}, {"slice3", "slice4"}, {"slice5", "slice6"}},
			[]string{"slice1", "slice2", "slice3", "slice4", "slice5", "slice6"}},
	}

	for _, tt := range testCasesForSlices {
		require.Equal(t, tt.expected, Concat(tt.input...))
	}

	sliceToUpdate := []string{"slice1", "slice2"}

	testCasesWithASliceUpdate := []struct {
		input    [][]string
		expected []string
	}{
		{[][]string{sliceToUpdate, {"slice3", "slice4"}},
			[]string{"slice1", "slice2", "slice3", "slice4"}},
	}

	for _, tt := range testCasesWithASliceUpdate {
		result := Concat(tt.input...)
		sliceToUpdate[0] = "abc"
		require.Equal(t, tt.expected, result)
		sliceToUpdate = []string{"slice6", "slice7"}
		require.Equal(t, tt.expected, result)
	}

	require.Equal(t, 0, len(Concat()))
}
