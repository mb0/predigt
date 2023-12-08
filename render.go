package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

func renderAll(dir string) error {
	txt := filepath.Join(dir, "text")
	tdir := filepath.Join(dir, "tmpl")
	out := filepath.Join(dir, "html")
	tmpl, err := template.ParseGlob(filepath.Join(tdir, "*.html"))
	if err != nil {
		return fmt.Errorf("could not parse templates: %v", err)
	}
	rc, err := readRenderConf(filepath.Join(tdir, "render.json"))
	if err != nil {
		return fmt.Errorf("could not parse render config: %v", err)
	}
	evs, nfos, err := rc.readEvents(txt)
	if err != nil {
		return err
	}
	for _, ev := range evs {
		for i, nfo := range ev.Vers {
			ser, err := readTextFile(txt, nfo)
			if err != nil {
				return fmt.Errorf("could not parse text file %d.txt, %v", nfo.ID, err)
			}
			if ev.Date.IsZero() {
				ev.Date = ser.Date
			}
			tout := fmt.Sprintf("%d.html", ser.ID)
			if !ser.Kids && i+1 < len(ev.Vers) && ev.Vers[i+1].Kids {
				ser.Kids = true
			}
			ctx := struct {
				*Sermon
				Ev *event
			}{Sermon: ser, Ev: ev}
			err = writeTemplate(tmpl, filepath.Join(out, tout), "sermon.html", ctx)
			if err != nil {
				return fmt.Errorf("could not write template %s, %v", tout, err)
			}
		}
	}
	ctx := struct {
		Evs     []*event
		Feature []*Feature `json:"feature"`
		Count   int
	}{Evs: evs, Feature: rc.Feature, Count: len(nfos)}
	for _, feat := range ctx.Feature {
		for _, id := range feat.IDs {
			for _, nfo := range nfos {
				if nfo.ID == id {
					feat.Infos = append(feat.Infos, nfo)
					break
				}
			}
		}
	}
	tout := filepath.Join(out, "index.html")
	err = writeTemplate(tmpl, tout, "index.html", ctx)
	if err != nil {
		return fmt.Errorf("could not write template %s, %v", tout, err)
	}
	tout = filepath.Join(out, "archiv.html")
	err = writeTemplate(tmpl, tout, "archiv.html", ctx)
	if err != nil {
		return fmt.Errorf("could not write template %s, %v", tout, err)
	}
	return nil
}

func writeTemplate(tmpl *template.Template, out, name string, ctx interface{}) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.ExecuteTemplate(f, name, ctx)
}

type Feature struct {
	Text  string
	IDs   []int
	Infos []*Info
}

type renderConfig struct {
	Ignore  []int      `json:"ignore"`
	Feature []*Feature `json:"feature"`
}

func readRenderConf(path string) (rc renderConfig, _ error) {
	f, err := os.Open(path)
	if err != nil {
		return rc, nil
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&rc)
	return rc, err
}

func (c renderConfig) Ignored(id int) bool {
	num := id / 10
	for _, ign := range c.Ignore {
		if ign == id || ign == num {
			return true
		}
	}
	return false
}

type event struct {
	Num    int
	Date   time.Time
	Sunday string
	Vers   []*Info
}

func (rc *renderConfig) readEvents(dir string) (evs []*event, res []*Info, _ error) {
	files, err := readFileList(dir)
	if err != nil {
		return nil, nil, err
	}
	var lst *event
	for _, file := range files {
		n, err := parseFilename(file)
		if err != nil {
			return nil, nil, err
		}
		if rc.Ignored(n.ID) {
			continue
		}
		res = append(res, n)
		num := n.ID / 10
		if lst != nil && lst.Num == num {
			lst.Vers = append(lst.Vers, n)
		} else {
			lst = &event{Num: num, Vers: []*Info{n}}
			evs = append(evs, lst)
		}
	}
	return evs, res, nil
}
