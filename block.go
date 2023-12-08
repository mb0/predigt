package main

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode/utf8"
)

type Doc []*Block

func (d Doc) WriteTo(w io.Writer) (int64, error) {
	b := bufio.NewWriter(w)
	var nn int64
	for _, bl := range d {
		for _, ln := range bl.Lines {
			n, err := b.WriteString(ln)
			b.WriteByte('\n')
			nn += int64(n + 1)
			if err != nil {
				return nn, err
			}
		}
		b.WriteByte('\n')
		nn += 1
	}
	return nn, b.Flush()
}

type Block struct {
	Lines []string
}

func (b *Block) String() string { return "[" + strings.Join(b.Lines, ", ") + "]" }

// Transformer is fead lines one by one
type Transformer struct {
	Trans func(*bytes.Buffer, string)
	Width int
	l     bytes.Buffer
	w     bytes.Buffer
	last  *Block
	ln    int
	doc   Doc
}

func (tr *Transformer) ReadBlocks(r io.Reader) (Doc, error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		tr.line(sc.Text())
	}
	tr.flushl()
	return tr.doc, sc.Err()
}

func (tr *Transformer) line(line string) {
	if len(line) == 0 {
		tr.flushl()
		tr.last = nil
		return
	}
	line = strings.TrimSpace(line)
	var start int
	for i, r := range line {
		switch r {
		case ' ', '\t':
			if i > start {
				word := line[start:i]
				tr.word(word)
			}
			start = i + 1
		case '-':
			if tr.Width > 0 && i > start+1 && i+1 == len(line) {
				tr.w.WriteString(line[start:i])
				start = i + 1
			}
		}
	}
	if start < len(line) {
		tr.word(line[start:])
	}
	if tr.Width == 0 {
		tr.flushl()
	}
}

func (tr *Transformer) word(word string) {
	if tr.Trans != nil {
		tr.Trans(&tr.w, word)
	} else {
		tr.w.WriteString(word)
	}
	tr.flushw()
}

func TransWord(w *bytes.Buffer, word string) {
	for i, r := range word {
		switch r {
		case '\'':
			if i+1 < len(word) && word[i+1] == 's' {
				w.WriteRune('’')
				continue
			}
		case ':':
			if i+1 < len(word) && word[i+1] == '"' {
				w.WriteString(": ")
			}
		case 'ß':
			start := i - 4
			if start < 0 {
				start = 0
			}
			runes := ([]rune)(strings.ToLower(word[start : i+2]))
			if len(runes) > 3 {
				runes = runes[len(runes)-3:]
			}
			if _, ok := sz[string(runes)]; ok {
				w.WriteString("ss")
				continue
			}
		}
		w.WriteRune(r)
	}
}

func (tr *Transformer) flushl() {
	if tr.l.Len() == 0 {
		return
	}
	if tr.last == nil {
		tr.last = &Block{}
		tr.doc = append(tr.doc, tr.last)
	}
	tr.last.Lines = append(tr.last.Lines, tr.l.String())
	tr.l.Reset()
	tr.ln = 0
}

func (tr *Transformer) flushw() {
	if tr.w.Len() == 0 {
		return
	}
	w := tr.w.String()
	wn := utf8.RuneCountInString(w)
	tr.w.Reset()
	if tr.Width > 0 && tr.ln+wn+1 > tr.Width {
		tr.flushl()
		tr.ln = 0
	}
	if tr.ln > 0 {
		tr.ln++
		tr.l.WriteByte(' ')
	}
	tr.ln += wn
	tr.l.WriteString(w)
}

var ok struct{}
var sz = map[string]struct{}{
	"biß": ok, // bisschen
	//"buß": !ok Buße
	"daß": ok,
	"eß":  ok,
	//"eiß": !ok, heißt
	"faß": ok, // fasste
	"geß": ok,
	"giß": ok,
	"goß": ok,
	"guß": ok, // Guss
	"haß": ok,
	"häß": ok,
	//"ieß": !ok, hieße
	"laß": ok,
	"läß": ok,
	"loß": ok, // Schloss
	"luß": ok,
	"muß": ok,
	//"maß": !ok gleichtermaßen
	//"mäß": !ok mäßig
	"miß": ok,
	"müß": ok,
	"paß": ok, // passt Spass
	//"raß": !ok Straße
	"reß": ok, // erpresst
	"riß": ok, // frisst
	//"roß": !ok große
	//"röß": !ok Größe
	"wiß": ok,
	"wuß": ok,
	"wüß": ok,
}
