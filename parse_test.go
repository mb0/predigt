package main

import (
	"log"
	"testing"
	"time"
)

func TestParseInfo(t *testing.T) {
	ser532 := &Info{
		ID:     5320,
		Date:   dateISO("1985-12-01"),
		Author: "Pastor Schnabel",
		Sunday: "1. Advent",
		Kids:   false,
		Theme:  "Röm. 13, 8-12",
	}
	ser532kids := &Info{
		ID:     5321,
		Date:   dateISO("1985-12-01"),
		Author: "Pastor Schnabel",
		Sunday: "1. Advent",
		Kids:   true,
		Theme:  "",
	}
	tests := []struct {
		file string
		text string
		want *Info
	}{
		{"5320.txt", "Predigt vom 01.12.1985 - Pastor Schnabel - 1. Advent - Röm. 13, 8-12\n", ser532},
		{"5320.txt", "\n\nPredigt vom 01.12.1985 - Pastor Schnabel - 1. Advent - Röm. 13, 8-12\n", ser532},
		{"5321.txt", "Kinderpredigt vom 01.12.1985 - Pastor Schnabel - 1. Advent", ser532kids},
	}
	for _, test := range tests {
		nfo, err := parseFilename(test.file)
		if err != nil {
			t.Errorf("%s error: %v", test.file, err)
		}
		err = parseDocInfo(nfo, Doc{{Lines: []string{test.text}}})
		if err != nil {
			t.Errorf("%s error: %v", test.file, err)
		}
		if *nfo != *test.want {
			t.Errorf("%s want %v got %v", test.file, *test.want, *nfo)
		}
	}
}

func dateISO(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		log.Printf("failed to parse date %s: %v", s, err)
	}
	return t
}
