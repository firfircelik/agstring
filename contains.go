package agstring

import (
	"reflect"
	"strings"

	"github.com/thoas/go-funk"
)

// ContainsSubString checks if the searched string is a sub string of
// the any string in the source slice
func ContainsSubString(src []string, searched string) bool {
	return IndexContainingSubString(src, searched) != -1
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
	return StringIndexContainingSubString(s, ls...) != -1
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

// IndexContainingSubString returns the first index in given slice
// which contains the searched string
func IndexContainingSubString(holder []string, searched string) int {
	for i := range holder {
		if strings.Contains(holder[i], searched) {
			return i
		}
	}
	return -1
}

// StringIndexContainingSubString returns the first index in given slice
// which is a sub string of the source string
func StringIndexContainingSubString(s string, ls ...string) int {
	for i := range ls {
		if strings.Contains(s, ls[i]) {
			return i
		}
	}
	return -1
}

func stringApplyTransforms(s string, transforms ...Transform) string {
	for _, t := range transforms {
		s = t(s)
	}
	return s
}

func applyTransforms(slice []string, transforms ...Transform) []string {
	transformed := make([]string, len(slice))
	for i := range slice {
		transformed[i] = stringApplyTransforms(slice[i], transforms...)
	}
	return transformed
}

// SliceContains checks if `slice` contains `s`
func SliceContains(slice [][]string, s []string, transform ...Transform) bool {
	for _, elem := range slice {
		if reflect.DeepEqual(s, applyTransforms(elem, transform...)) {
			return true
		}
	}
	return false
}
