package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var serid = regexp.MustCompile(`^([0-9]+)(.(?:txt|html))?$`)

type Sermon struct {
	*Info
	Doc Doc
}

func readTextFile(dir string, nfo *Info) (*Sermon, error) {
	f, err := os.Open(filepath.Join(dir, fmt.Sprintf("%d.txt", nfo.ID)))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tr := &Transformer{}
	doc, err := tr.ReadBlocks(f)
	if err != nil {
		return nil, err
	}
	err = parseDocInfo(nfo, doc)
	if err != nil {
		return nil, err
	}
	return &Sermon{Info: nfo, Doc: doc}, nil
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
	Date   time.Time
	Author string
	Sunday string
	Kids   bool
	Theme  string
}

func (n *Info) String() string {
	name, xtra := "Predigt", ""
	num := n.ID / 10
	lst := n.ID % 10
	if lst == 7 {
		name = "Taufpredigt"
	} else if lst%2 == 1 {
		name = "Kinderpredigt"
	}
	if lst > 1 && lst < 7 {
		idx := (lst - 2) / 2
		xtra = " " + "abc"[idx:idx+1]
	}
	return fmt.Sprintf("%s %d%s", name, num, xtra)
}
func parseFilename(file string) (*Info, error) {
	base := filepath.Base(file)
	num, _, _ := strings.Cut(base, ".")
	id, err := strconv.Atoi(num)
	if err != nil {
		return nil, err
	}
	return &Info{ID: id, Kids: id%2 == 1}, nil
}

func parseDocInfo(nfo *Info, doc Doc) (err error) {
	if len(doc) == 0 || len(doc[0].Lines) == 0 {
		return fmt.Errorf("no lines in doc %d", nfo.ID)
	}
	fst := doc[0].Lines[0]
	parts := strings.SplitN(fst, " - ", 4)
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	// (Kinder)[Pp]redigt vom {date}
	_, date, ok := strings.Cut(parts[0], " vom ")
	if !ok {
		return fmt.Errorf("no date found in doc %d - %q", nfo.ID, parts[0])
	}
	nfo.Date, err = time.Parse("02.01.2006", date)
	if err != nil {
		return err
	}
	kids := strings.HasPrefix(parts[0], "Kinder")
	switch len(parts) {
	case 2:
		if strings.HasPrefix(parts[1], "Pastor") {
			nfo.Author = parts[1]
		} else {
			nfo.Theme = parts[1]
		}
	case 3:
		nfo.Author = parts[1]
		if kids {
			nfo.Sunday = parts[2]
		} else {
			nfo.Theme = parts[2]
		}
	case 4:
		nfo.Author = parts[1]
		nfo.Sunday = parts[2]
		nfo.Theme = parts[3]
	}
	return nil
}

var matchTheme = regexp.MustCompile(`((?:\d+.)?\s?\w+.?(:?\s?\d+,))\s?(\d+(?:-\d|ff))`)
