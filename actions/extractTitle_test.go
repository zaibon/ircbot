package actions

import "testing"

var instance = &TitleExtract{}

var testTable = []struct {
	input  string
	expect string
}{
	{"http://google.com", "Google"},
	{"http://linux.org", "Linux.org"},
}

func TestTitleExtract(t *testing.T) {
	for _, tt := range testTable {
		actual, err := extractTitle(tt.input)
		if err != nil {
			t.Error(err)
		}
		if actual != tt.expect {
			t.Errorf("input %s\nexpected %s\nactual %s\n", tt.input, tt.expect, actual)
		}
	}
}

var result string

func BenchmarkTitleExtract(b *testing.B) {
	var t string
	for i := 0; i < b.N; i++ {
		t, _ = extractTitle("http://google.com")
	}
	result = t
}
