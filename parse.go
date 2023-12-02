package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var serid = regexp.MustCompile(`^([0-9]+)(?:_([a-z]+))?(.(?:txt|html))?$`)

type Sermon struct {
	*Info
	Body string
}

func readTextFile(dir, file string) (*Sermon, error) {
	raw, err := ioutil.ReadFile(filepath.Join(dir, file))
	if err != nil {
		return nil, err
	}
	text := string(raw)
	info, err := parseInfo(file, text)
	if err != nil {
		return nil, err
	}
	return &Sermon{Info: info, Body: text}, nil
}

func readFileList(dir string) ([]string, error) {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var list []string
	for _, fi := range fis {
		name := fi.Name()
		if serid.MatchString(name) {
			list = append(list, name)
		}
	}
	sort.Strings(list)
	return list, nil
}

type Info struct {
	ID     int
	Base   string
	Date   time.Time
	Author string
	Sunday string
	Kids   bool
	Theme  string
}

func parseInfo(file string, text string) (*Info, error) {
	base := filepath.Base(file)
	name, _, _ := strings.Cut(base, ".")
	num, vers, _ := strings.Cut(name, "_")
	id, err := strconv.Atoi(num)
	if err != nil {
		return nil, err
	}
	ser := &Info{ID: id, Base: name, Kids: vers == "kids"}

	text = strings.TrimSpace(text)
	fst, _, _ := strings.Cut(text, "\n")
	parts := strings.SplitN(fst, " - ", 4)
	// (Kinder)[Pp]redigt vom {date}
	_, date, ok := strings.Cut(strings.TrimSpace(parts[0]), " vom ")
	if !ok {
		return nil, fmt.Errorf("no date found in %s - %s", file, parts[0])
	}
	ser.Date, err = time.Parse("02.01.2006", date)
	if err != nil {
		return nil, err
	}
	ser.Author = strings.TrimSpace(parts[1])
	ser.Sunday = strings.TrimSpace(parts[2])
	if len(parts) > 3 {
		ser.Theme = strings.TrimSpace(parts[3])
	}
	return ser, nil
}
