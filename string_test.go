package agstring

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestReplaceMultispace(t *testing.T) {
	tests := []string{"   a   b   c d    ef  ", "     123   "}
	expects := []string{"a b c d ef", "123"}

	for i, tt := range tests {
		res := ReplaceMultispace(tt)
		if res != expects[i] {
			t.Errorf("%d '%s': Expected %s, got %s", i, tt, expects[i], res)
		}
	}
}

func TestFirst(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{[]string{}, ""},
		{[]string{"one"}, "one"},
		{[]string{"one", "two"}, "one"},
	}

	for _, tt := range tests {
		if res := First(tt.input); res != tt.expected {
			t.Errorf("For input %v, expected %s, got %s", tt.input, tt.expected, res)
		}
	}
}

func TestLast(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{[]string{}, ""},
		{[]string{"one"}, "one"},
		{[]string{"one", "two"}, "two"},
	}

	for _, tt := range tests {
		if res := Last(tt.input); res != tt.expected {
			t.Errorf("For input %v, expected %s, got %s", tt.input, tt.expected, res)
		}
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
		if res := RemoveDiacritics(tt.input); res != tt.expected {
			t.Errorf("For input %v, expected %s, got %s", tt.input, tt.expected, res)
		}
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
	for i, tt := range tests {
		if res := EmptyIf(tt.s, tt.emptyList...); res != tt.expected {
			t.Errorf("%d: For input %q with empty list %v, expected %q, got %q",
				i, tt.s, tt.emptyList, tt.expected, res)
		}
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
	for i, tt := range tests {
		if res := TrimSpace(tt.input); !reflect.DeepEqual(res, tt.expected) {
			t.Errorf("%d: For input %v, expected %v, not %v", i, tt.input, tt.expected, res)
		}
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
	for i, tt := range tests {
		if res := Title(tt.input); res != tt.expected {
			t.Errorf("%d: For input %v, expected %v, not %v", i, tt.input, tt.expected, res)
		}
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

	for _, testCase := range testCases {
		require.Equalf(t,
			testCase.expected,
			ContainsAll(testCase.holder, testCase.searched...),
			"%v is searched in %v and fails")
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

	for _, testCase := range testCases {
		require.Equalf(t,
			testCase.expected,
			StringContainsAll(testCase.holder, testCase.searched...),
			"%v is searched in %q and fails")
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

	for _, testCase := range testCases {
		require.Equalf(t,
			testCase.expected,
			ContainsAny(testCase.holder, testCase.searched...),
			"%v is searched in %v and fails", testCase.holder, testCase.searched)
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

	for _, testCase := range testCases {
		require.Equalf(t,
			testCase.expected,
			StringContainsAny(testCase.holder, testCase.searched...),
			"%v is searched in %q and fails")
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
	for _, testCase := range testCases {
		match, ok := RegexpGroups(testCase.re, testCase.input)
		require.Equal(t, testCase.matches, ok, "for input %s", testCase.input)
		if testCase.matches {
			require.Equal(t, testCase.expected, match)
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

	for _, testCase := range testCases {
		result := TrimSuffixes(testCase.input, testCase.suffixes...)
		require.Equal(t, testCase.expected, result)
	}
}

func TestTrimPrefixesAndSpace(t *testing.T) {
	testCases := []struct {
		input    string
		prefixes []string
		expected string
	}{
		{" cif fob massive port", []string{"fob", "cif"}, "massive port"},
		{" cif fobmassive port", []string{"fob", "cif"}, "fobmassive port"},
		{" cif fob massive port", []string{"big", "location"}, "cif fob massive port"},
		{" cif fob massive port", []string{""}, "cif fob massive port"},
		{" cif fob", []string{"fob", "cif"}, ""},
		{" ", []string{}, ""},
		{"", nil, ""},
	}

	for _, testCase := range testCases {
		require.Equal(t,
			testCase.expected,
			TrimPrefixesAndSpace(testCase.input, testCase.prefixes),
		)
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

	for _, testCase := range testCases {
		require.Equal(t,
			testCase.expected,
			TakeTo(testCase.input, testCase.truncateTill),
		)
	}
}

func TestTakeFrom(t *testing.T) {
	testCases := []struct {
		input        []string
		truncateTill int
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

	for _, testCase := range testCases {
		require.Equal(t,
			testCase.expected,
			TakeFrom(testCase.input, testCase.truncateTill),
		)
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

	for _, testCase := range testCases {
		require.Equal(t,
			testCase.expected,
			ReplaceDayOrdinal(testCase.input, testCase.replaceWith),
		)
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

	for _, testCase := range testCases {
		require.Equal(t,
			testCase.expected,
			ReplaceNewline(input, testCase.replaceWith),
		)
	}

	require.Equal(t, "HiThere", ReplaceNewline(input))
}

func TestMap(t *testing.T) {
	f1 := func(s string) string { return strings.ToTitle(s) }
	f2 := func(s string) string { return strings.Repeat(s, 3) }

	testCases := []struct {
		holder   []string
		funcs    []func(string) string
		expected []string
	}{
		{[]string{"hello"}, []func(string) string{f1}, []string{"HELLO"}},
		{[]string{"hello"}, []func(string) string{f1, f2}, []string{"HELLOHELLOHELLO"}},
	}

	for i, testCase := range testCases {
		require.Equalf(t,
			testCase.expected,
			Map(testCase.holder, testCase.funcs...),
			"%v index is fails", i)
	}

}
