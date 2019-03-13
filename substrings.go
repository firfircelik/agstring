package agstring

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

// TakeBetween gets the string from position `from` up to `to`
// from each string in the input slice
func TakeBetween(ls []string, from, to int) []string {
	return TakeFrom(TakeTo(ls, to), from)
}
