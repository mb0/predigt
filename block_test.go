package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestReadBlocks(t *testing.T) {
	tests := []struct {
		name  string
		raw   string
		width int
		want  Doc
	}{
		{"empty", "", 0, nil},
		{"empty lines", "\n\n\n", 0, nil},
		{"word", "word", 0, Doc{
			{Lines: []string{"word"}},
		}},
		{"spaces", "  foo  bar  ", 0, Doc{
			{Lines: []string{"foo bar"}},
		}},
		{"width", "  foo  bar  ", 3, Doc{
			{Lines: []string{"foo", "bar"}},
		}},
		{"blocks", "\n\nfoo bar\n\nspam egg\n\n\n", 0, Doc{
			{Lines: []string{"foo bar"}},
			{Lines: []string{"spam egg"}},
		}},
		{"join lines", "foo\nbar\n\nspam\negg\ngoo", 7, Doc{
			{Lines: []string{"foo bar"}},
			{Lines: []string{"spam", "egg goo"}},
		}},
		{"join lines", "äää\nööö\n\nüüüü\nßßß\ngoo", 7, Doc{
			{Lines: []string{"äää ööö"}},
			{Lines: []string{"üüüü", "ßßß goo"}},
		}},
		{"join word", "foo-\nbar\nspam-\negg\ngoo", 7, Doc{
			{Lines: []string{"foobar", "spamegg", "goo"}},
		}},
		{"replace sz", "Müßt Haß Daß Wißt", 0, Doc{
			{Lines: []string{"Müsst Hass Dass Wisst"}},
		}},
		{"replace apo", "wir's", 0, Doc{
			{Lines: []string{"wir’s"}},
		}},
	}
	for _, test := range tests {
		tr := &Transformer{Width: test.width, Trans: TransWord}
		doc, err := tr.ReadBlocks(strings.NewReader(test.raw))
		if err != nil {
			t.Errorf("%s error: %v", test.name, err)
			continue
		}
		if !reflect.DeepEqual(test.want, doc) {
			t.Errorf("%s want %s got %s", test.name, test.want, doc)
			continue
		}
	}
}
