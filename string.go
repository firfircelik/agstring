package agstring

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/mozillazg/go-unidecode"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

// Nth returns nth element of given slice or empty string if out of limits
func Nth(ls []string, n int) string {
	if len(ls) == 0 || n < 0 || n >= len(ls) {
		return ""
	}
	return ls[n]
}

// First returns the first element of given list or empty string when the list is empty.
func First(ls ...string) string { return Nth(ls, 0) }

// Last returns the last element of given list or empty string when the list is empty.
func Last(ls ...string) string { return Nth(ls, len(ls)-1) }

var nonAlphanumRegexp = regexp.MustCompile("[^[:alnum:]]")

// RemoveNonAlnum removes non-alphanumeric characters from string
func RemoveNonAlnum(s string) string {
	return nonAlphanumRegexp.ReplaceAllLiteralString(s, "")
}

// RemoveDiacritics removes diacritics from a string. If non-alphanumeric character is
// encountered diacritics are removed from it. If removing diacritics is not possible, character
// is removed.
func RemoveDiacritics(s string) string { return unidecode.Unidecode(s) }

// EmptyIf returns empty string if given string equals to one
// of the strings in empty list. Otherwise, given string is returned as it is.
func EmptyIf(s string, emptyList ...string) string {
	return ConvertIf(s, "", emptyList...)
}

// ConvertIf returns converted string if given string is one of the strings in the list
func ConvertIf(val, converted string, list ...string) string {
	for _, t := range list {
		if val == t {
			return converted
		}
	}
	return val
}

// ValueIfExists returns value from map for given key if exists, else returns the given key
func ValueIfExists(k string, m map[string]string) string {
	v, ok := m[k]
	if ok {
		return v
	}
	return k
}

// StringIterator provides a generator of names / strings
type StringIterator interface {
	Get() string
	HasNext() bool
}

// ToLower makes lowercase strings in the given slice
func ToLower(ls ...string) []string { return Map(ls, strings.ToLower) }

// Title ensures title formatting for given string
func Title(s string) string { return strings.Title(strings.ToLower(s)) }

// HasSuffix checks string has any one of given suffixes
func HasSuffix(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

// HasPrefix checks string has any one of given prefixes
func HasPrefix(s string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

// NonEmpty filters nonempty strings from given slice
func NonEmpty(ls ...string) []string {
	return nonempty(ls, UnitTransform)
}

// NonEmptyIfTrimmed filters nonempty string only if
// contains some data when whitespace is removed
func NonEmptyIfTrimmed(ls ...string) []string {
	return nonempty(ls, strings.TrimSpace)
}

// Transform modifies a string to another format
type Transform func(string) string

// UnitTransform doesn't modify its argument
func UnitTransform(s string) string { return s }

func nonempty(ls []string, t Transform) []string {
	var nonempty []string
	for _, s := range ls {
		if t(s) != "" {
			nonempty = append(nonempty, s)
		}
	}
	return nonempty
}

// IsEmpty checks if slice contains only empty strings
func IsEmpty(ls ...string) bool { return len(NonEmpty(ls...)) == 0 }

// IsEqual checks if two slices are equal after applying transforms on 2nd slice
func IsEqual(source, other []string, transforms ...Transform) bool {
	return (source == nil && other == nil) ||
		funk.Equal(source, Map(other, transforms...))
}

// MustEqual checks if two slices are equal after applying transforms on 2nd slice otherwise return error
func MustEqual(source, other []string, transforms ...Transform) error {
	if !IsEqual(source, other, transforms...) {
		return errors.Errorf("must be equal: expected %v, got %v", source, other)
	}
	return nil
}

// RemoveAllDiacritics removes diacritics from all strings in slice
func RemoveAllDiacritics(ls ...string) []string { return Map(ls, RemoveDiacritics) }

// SafeAtoi converts string, including empty string, to int
func SafeAtoi(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	n, err := strconv.Atoi(s)
	return n, errors.Wrap(err, "can't convert to int")
}

// RegexpGroups checks if regex matches to given string
// If so, returns named groups with matches in a map
func RegexpGroups(exp *regexp.Regexp, input string) (map[string]string, bool) {
	if !exp.MatchString(input) {
		return nil, false
	}
	match := exp.FindStringSubmatch(input)
	result := make(map[string]string)
	for i, name := range exp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result, true
}

// Reprintf returns a string according to fmt format escaping the parameters before.
func Reprintf(format string, params ...string) string {
	var interfaceSlice = make([]interface{}, len(params))
	for i, d := range params {
		interfaceSlice[i] = regexp.QuoteMeta(d)
	}
	return fmt.Sprintf(format, interfaceSlice...)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Map runs given modifiers for each item in slice and returns a new slice
func Map(ls []string, transforms ...Transform) []string {
	out := make([]string, len(ls))
	for i, s := range ls {
		for _, transform := range transforms {
			s = transform(s)
		}
		out[i] = s
	}
	return out
}

// Concat returns a new slice with all the elements of slices on it
func Concat(slices ...[]string) []string {
	var concatSlice []string
	for _, slice := range slices {
		concatSlice = append(concatSlice, slice...)
	}
	return concatSlice
}

// Omit excludes given string from the slice
func Omit(ls []string, omitted string) []string {
	var res []string
	for _, s := range ls {
		if s != omitted {
			res = append(res, s)
		}
	}
	return res
}
