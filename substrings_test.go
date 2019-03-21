package agstring

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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

func TestTakeBetween(t *testing.T) {
	testCases := []struct {
		input    []string
		from     int
		to       int
		expected []string
	}{
		{
			[]string{"ąćęłńóśźż", "ĄĆĘŁŃÓŚŹŻ"},
			3,
			5,
			[]string{"łń", "ŁŃ"},
		},
		{
			[]string{"ドンキーコング", "hello"},
			100,
			9,
			[]string{"", ""},
		},
		{
			[]string{"ドンキーコング", "hello"},
			40,
			50,
			[]string{"", ""},
		},
		{
			[]string{"ドンキーコング", "hello"},
			0,
			3,
			[]string{"ドンキ", "hel"},
		},
		{
			[]string{"1234567890", "1234567890"},
			5,
			15,
			[]string{"67890", "67890"},
		},
	}

	for _, tt := range testCases {
		require.Equal(t, tt.expected, TakeBetween(tt.input, tt.from, tt.to))
	}
}

func TestTakeAround(t *testing.T) {
	testCases := []struct {
		input    []string
		from     int
		to       int
		expected []string
	}{
		{
			[]string{"ąćęłńóśźż", "ĄĆĘŁŃÓŚŹŻ"},
			3,
			5,
			[]string{"ąćęóśźż", "ĄĆĘÓŚŹŻ"},
		},
		{
			[]string{"ドンキーコング", "hello"},
			100,
			9,
			[]string{"ドンキーコング", "hello"},
		},
		{
			[]string{"ドンキーコング", "hello"},
			40,
			50,
			[]string{"ドンキーコング", "hello"},
		},
		{
			[]string{"ドンキーコング", "hello"},
			0,
			3,
			[]string{"ーコング", "lo"},
		},
		{
			[]string{"1234567890", "1234567890"},
			5,
			15,
			[]string{"12345", "12345"},
		},
		{
			[]string{"1234567890", "1234567890"},
			5,
			8,
			[]string{"1234590", "1234590"},
		},
		{
			[]string{"abc def ghi", "ab cd ef hi"},
			5,
			9,
			[]string{"abc dhi", "ab cdhi"},
		},
	}

	for _, tt := range testCases {
		require.Equal(t, tt.expected, TakeAround(tt.input, tt.from, tt.to))
	}
}
