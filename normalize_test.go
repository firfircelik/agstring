package agstring

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"   sdasdaSDAda", "sdasdasdada"},
		{"  qwesDsfsdf      ", "qwesdsfsdf"},
		{"sdAAAaA    sdas \n", "sdaaaaa    sdas"},
		{" \t\n Hello World \n\t\r\n", "hello world"},
		{"\u200BHowever thÎs prÔject is fÛn", "however thîs prôject is fûn"},
		{"However thÎs prÔject is fÛn", "however thîs prôject is fûn"},
		{"", ""},
	}

	for _, tt := range tests {
		require.Equal(t, tt.expected, Normalize(tt.input))
	}
}

func TestNormalizeDiacritics(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ąćęłńóśźż ĄĆĘŁŃÓŚŹŻ", "acelnoszz acelnoszz"},                         // Polish
		{"čďěňřšťž ČĎĚŇŘŠŤŽ", "cdenrstz cdenrstz"},                             // Czech
		{"éàèùâêîôûëïöüÿç ÉÀÈÙÂÊÎÔÛËÏÖÜŸÇ", "eaeuaeioueiouyc eaeuaeioueiouyc"}, // French
		{"ğöüçş ĞÖÜÇŞ", "goucs goucs"},
		{"However thÎs prÔject is fÛn", "however this project is fun"},
		{"", ""},
		{"7/8/9月", "7/8/9yue"},
	}

	for _, tt := range tests {
		require.Equal(t, tt.expected, NormalizeDiacritics(tt.input))
	}
}

func TestNormalizeDiacriticsAndNonAlnum(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello-World", "helloworld"},
		{"  dont';'; mind    ", "dontmind"},
		{"the pretty hou[]$%$żżse un89", "theprettyhouzzseun89"},
		{"However thÎs prÔject is fÛn", "howeverthisprojectisfun"},
		{"", ""},
	}

	for _, tt := range tests {
		require.Equal(t, tt.expected, NormalizeDiacriticsAndNonAlnum(tt.input))
	}
}
