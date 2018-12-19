package agstring

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/mozillazg/go-unidecode"
	"github.com/pkg/errors"
	funk "github.com/thoas/go-funk"
)

var stripper = regexp.MustCompile("  +")

// ReplaceMultispace replaces multiple spaces with one space and
// also trims space from both ends
func ReplaceMultispace(s string) string {
	return strings.TrimSpace(stripper.ReplaceAllString(s, " "))
}

// First returns the first element of given list or empty string when the list is empty.
func First(strings []string) string {
	if len(strings) != 0 {
		return strings[0]
	}
	return ""
}

// Lastt returns the last  element of given list or empty string when the list is empty.
func Last(strings []string) string {
	if len(strings) != 0 {
		return strings[len(strings)-1]
	}
	return ""
}

// TrimSuffixes returns s without any of the provided trailing suffixes strings.
func TrimSuffixes(s string, suffixes ...string) string {
	s = strings.TrimSpace(s)
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return strings.TrimSpace(strings.TrimSuffix(s, suffix))
		}
	}
	return s
}

// TrimPrefixesAndSpace returns a string without any of the provided leading prefixes at word
// boundaries or spaces. See test for examples.
func TrimPrefixesAndSpace(s string, prefixes []string) string {
	if prefixes == nil || s == "" {
		return s
	}
	trimAgain := true
	rePre := make([]*regexp.Regexp, len(prefixes))
	for trimAgain {
		trimAgain = false
		for i, prefix := range prefixes {
			if prefix == "" {
				continue
			}
			if rePre[i] == nil {
				reE := fmt.Sprintf("^\\s*%s\\b(?P<rest>.*)", prefix)
				rePre[i] = regexp.MustCompile(reE)
			}
			if matches, ok := RegexpGroups(rePre[i], strings.TrimSpace(s)); ok {
				s = matches["rest"]
				trimAgain = true
			}
		}
	}
	return strings.TrimSpace(s)
}

var nonAlphanumRegexp = regexp.MustCompile("[^[:alnum:]]")

// RemoveNonAlnum removes non-alphanumeric characters from string
func RemoveNonAlnum(s string) string {
	return nonAlphanumRegexp.ReplaceAllLiteralString(s, "")
}

// ContainsAll checks if given slice contains all searched strings
func ContainsAll(holder []string, searched ...string) bool {
	for _, s := range searched {
		if !funk.ContainsString(holder, s) {
			return false
		}
	}
	return true
}

// StringContainsAll checks if given string contains all searched strings
func StringContainsAll(holder string, searched ...string) bool {
	for _, s := range searched {
		if !strings.Contains(holder, s) {
			return false
		}
	}
	return true
}

// ContainsAny checks if source slice contains any of given strings
func ContainsAny(src []string, qs ...string) bool {
	for _, q := range qs {
		if funk.ContainsString(src, q) {
			return true
		}
	}
	return false
}

// StringContainsAny is similar to ContainsAny but source is a string
func StringContainsAny(s string, ls ...string) bool {
	for _, e := range ls {
		if strings.Contains(s, e) {
			return true
		}
	}
	return false
}

// MatchesPrefixes checks if given string has a prefix from given prefix list
func MatchesPrefixes(s string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

// RemoveDiacritics removes diacritics from a string. If non-alphanumeric character is
// encountered diacritics are removed from it. If removing diacritics is not possible, character
// is removed.
func RemoveDiacritics(s string) string { return unidecode.Unidecode(s) }

// Normalize tries to remove the diacritics, removes remaining non-alphanumeric characters and
// then changes case to lower
func Normalize(s string) string {
	s = RemoveDiacritics(s)
	s = RemoveNonAlnum(s)
	return strings.ToLower(s)
}

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

// ReplaceWholeWord replaces old into new only if old occurs as a whole word.
func ReplaceWholeWord(s, old, replacement string) string {
	s = " " + s + " "
	old = " " + old + " "
	replacement = " " + replacement + " "
	s = strings.Replace(s, old, replacement, -1)
	return s[1 : len(s)-1]
}

// StringIterator provides a generator of names / strings
type StringIterator interface {
	Get() string
	HasNext() bool
}

// TrimSpace trims spaces in the given slice
func TrimSpace(ls []string) []string {
	for i := range ls {
		ls[i] = strings.TrimSpace(ls[i])
	}
	return ls
}

// ToLower makes lowercase strings in the given slice
func ToLower(ls []string) []string {
	for i := range ls {
		ls[i] = strings.ToLower(ls[i])
	}
	return ls
}

// Title ensures title formatting for given string
func Title(s string) string { return strings.Title(strings.ToLower(s)) }

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
func NonEmpty(ls []string) []string {
	var nonempty []string
	for _, s := range ls {
		if s != "" {
			nonempty = append(nonempty, s)
		}
	}
	return nonempty
}

// IsEmpty checks if slice contains only empty strings
func IsEmpty(ls []string) bool {
	for _, s := range ls {
		if s != "" {
			return false
		}
	}
	return true
}

// RemoveAllDiacritics removes diacritics from all strings in slice
func RemoveAllDiacritics(ls []string) []string {
	res := make([]string, len(ls))
	for i, s := range ls {
		res[i] = RemoveDiacritics(s)
	}
	return res
}

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

// TakeTo truncates each string in the input slice up to `n` characters.
func TakeTo(ls []string, n int) []string {
	out := make([]string, 0, len(ls))
	for _, s := range ls {
		rs := []rune(s)
		o := string(rs[:min(len(rs), n)])
		out = append(out, o)
	}
	return out
}

// TakeFrom removes the first `n` characters from each string in the input slice
func TakeFrom(ls []string, n int) []string {
	out := make([]string, 0, len(ls))
	for _, s := range ls {
		rs := []rune(s)
		o := string(rs[min(len(rs), n):])
		out = append(out, o)
	}
	return out
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ReplaceDayOrdinal replaces day ordinals (`st`, `nd`, `rd`, `th`)
// Default replaces with empty string.
func ReplaceDayOrdinal(s string, replacements ...string) string {
	var rep string
	if len(replacements) > 0 {
		rep = replacements[0]
	}
	ordinal := strings.NewReplacer("st", rep, "nd", rep, "th", rep, "rd", rep)
	return ordinal.Replace(s)
}

// ReplaceNewline replaces the newline character `\n`
// Default replaces with empty string.
func ReplaceNewline(s string, replacements ...string) string {
	var rep string
	if len(replacements) > 0 {
		rep = replacements[0]
	}
	return strings.Replace(s, "\n", rep, -1)
}

func Map(ls []string, funcs ...func(string) string) []string {
	out := make([]string, len(ls))
	for i, s := range ls {
		tmp := s
		for _, f := range funcs {
			tmp = f(tmp)
		}
		out[i] = tmp
	}
	return out
}
