package agstring

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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

func TestMap(t *testing.T) {
	repeat := func(s string) string { return strings.Repeat(s, 3) }

	testCases := []struct {
		input      []string
		transforms []Transform
		expected   []string
	}{
		{input: []string{"hello"},
			transforms: []Transform{strings.ToTitle},
			expected:   []string{"HELLO"}},
		{input: []string{"hello"},
			transforms: []Transform{strings.ToTitle, repeat},
			expected:   []string{"HELLOHELLOHELLO"}},
	}

	for _, tt := range testCases {
		require.Equal(t, tt.expected, Map(tt.input, tt.transforms...))
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

func TestIsEqual(t *testing.T) {
	addA := func(s string) string { return s + "a" }
	testCases := []struct {
		input, probe []string
		transforms   []Transform
		expected     bool
	}{
		{input: []string{}, probe: []string{}, expected: true},
		{input: []string{}, probe: nil, expected: true},
		{input: nil, probe: nil, expected: true},
		{input: []string{"", "a"}, probe: []string{"a", ""}, expected: false},
		{input: []string{}, probe: []string{},
			transforms: []Transform{UnitTransform}, expected: true},
		{input: []string{"a"}, probe: []string{"a"},
			transforms: []Transform{UnitTransform}, expected: true},
		{input: []string{"aa"}, probe: []string{"a"},
			transforms: []Transform{addA, UnitTransform}, expected: true},
		{input: []string{"ab"}, probe: []string{"a"},
			transforms: []Transform{addA, UnitTransform}, expected: false},
		{input: []string{"ab", "aa"}, probe: []string{"a", "a"},
			transforms: []Transform{addA, UnitTransform}, expected: false},
	}

	for _, testCase := range testCases {
		require.Equalf(t, testCase.expected,
			IsEqual(testCase.input, testCase.probe, testCase.transforms...),
			"input is %v", testCase)
	}
}

func TestReprintf(t *testing.T) {
	testCases := []struct {
		format   string
		params   []string
		expected string
	}{
		{format: "%s", params: []string{`This: .+`}, expected: "This: \\.\\+"},
	}
	for _, testCase := range testCases {
		require.Equal(t, testCase.expected, Reprintf(testCase.format, testCase.params...))
	}
}

func TestOmit(t *testing.T) {
	testCases := []struct {
		input, expected []string
		omit            string
	}{
		{
			input:    nil,
			expected: nil,
			omit:     "agflow",
		},
		{
			input:    []string{},
			expected: nil,
			omit:     "agflow",
		},
		{
			input:    []string{"agriculture"},
			expected: []string{"agriculture"},
			omit:     "agflow",
		},
		{
			input:    []string{"agriculture", "agflow", "intelligence"},
			expected: []string{"agriculture", "intelligence"},
			omit:     "agflow",
		},
	}

	for _, testCase := range testCases {
		require.Equal(t, testCase.expected, Omit(testCase.input, testCase.omit))
	}
}

func TestMustEqual(t *testing.T) {
	addA := func(s string) string { return s + "a" }
	testCases := []struct {
		input, probe []string
		transforms   []Transform
		expected     error
	}{
		{input: []string{}, probe: []string{}},
		{input: []string{}, probe: nil},
		{input: nil, probe: nil},
		{input: []string{"a", "a"}, probe: []string{"a", "a"}},
		{input: []string{}, probe: []string{},
			transforms: []Transform{UnitTransform}},
		{input: []string{"a"}, probe: []string{"a"},
			transforms: []Transform{UnitTransform}},
		{input: []string{"aa"}, probe: []string{"aa"}},
		{input: []string{"aa"}, probe: []string{"a"},
			transforms: []Transform{addA, UnitTransform}},
	}

	for _, testCase := range testCases {
		require.Equalf(t, testCase.expected,
			MustEqual(testCase.input, testCase.probe, testCase.transforms...),
			"input is %v", testCase)
	}
}
