package agstring

import (
	"fmt"
	"regexp"
	"strings"
)

// TrimAndRemoveZeroWidthChars trims spaces and removes ZeroWidth chars
func TrimAndRemoveZeroWidthChars(s string) string {
	return strings.TrimSpace(ReplaceZeroWidthChars(s))
}

// TrimSpace trims spaces in the given slice
func TrimSpace(ls ...string) []string { return Map(ls, strings.TrimSpace) }

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

// TrimAllSuffixes returns a string without any of the provided trailing suffixes or spaces.
// See test for examples.
func TrimAllSuffixes(s string, suffixes ...string) string {
	if len(suffixes) == 0 || s == "" {
		return strings.TrimSpace(s)
	}

	reSufs := make([]*regexp.Regexp, 0)
	for _, suffix := range suffixes {
		if suffix == "" {
			continue
		}

		reE := fmt.Sprintf(`^(?P<rest>.*)%s\s*$`, regexp.QuoteMeta(suffix))
		reSufs = append(reSufs, regexp.MustCompile(reE))
	}

	trimAgain := true
	for trimAgain {
		trimAgain = false
		for _, reSuf := range reSufs {
			if matches, ok := RegexpGroups(reSuf, strings.TrimSpace(s)); ok {
				s = matches["rest"]
				trimAgain = true
			}
		}
	}
	return strings.TrimSpace(s)
}

// TrimAllPrefixes returns a string without any of the provided leading prefixes or spaces.
// See test for examples.
func TrimAllPrefixes(s string, prefixes ...string) string {
	if len(prefixes) == 0 || s == "" {
		return strings.TrimSpace(s)
	}

	rePres := make([]*regexp.Regexp, 0)
	for _, prefix := range prefixes {
		if prefix == "" {
			continue
		}

		reE := fmt.Sprintf(`^\s*%s(?P<rest>.*)`, regexp.QuoteMeta(prefix))
		rePres = append(rePres, regexp.MustCompile(reE))
	}

	trimAgain := true
	for trimAgain {
		trimAgain = false
		for _, rePre := range rePres {
			if matches, ok := RegexpGroups(rePre, strings.TrimSpace(s)); ok {
				s = matches["rest"]
				trimAgain = true
			}
		}
	}
	return strings.TrimSpace(s)
}

// TrimPrefixesAndSpace returns a string without any of the provided leading prefixes at word
// boundaries or spaces. See test for examples.
func TrimPrefixesAndSpace(s string, prefixes ...string) string {
	if prefixes == nil || s == "" {
		return s
	}

	rePres := make([]*regexp.Regexp, 0)
	for _, prefix := range prefixes {
		if prefix == "" {
			continue
		}

		reE := fmt.Sprintf("^\\s*%s\\b(?P<rest>.*)", regexp.QuoteMeta(prefix))
		rePres = append(rePres, regexp.MustCompile(reE))
	}

	trimAgain := true
	for trimAgain {
		trimAgain = false
		for _, rePre := range rePres {
			if matches, ok := RegexpGroups(rePre, strings.TrimSpace(s)); ok {
				s = matches["rest"]
				trimAgain = true
			}
		}
	}
	return strings.TrimSpace(s)
}
