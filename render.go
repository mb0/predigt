package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

func renderAll(dir string) error {
	txt := filepath.Join(dir, "text")
	glob := filepath.Join(dir, "tmpl", "*.html")
	out := filepath.Join(dir, "html")
	files, err := readFileList(txt)
	if err != nil {
		return err
	}
	tmpl, err := template.ParseGlob(glob)
	if err != nil {
		return fmt.Errorf("could not parse templates: %v", err)
	}
	var infos []*Info
	for i, file := range files {
		ser, err := readTextFile(txt, file)
		if err != nil {
			return fmt.Errorf("could not parse text file %s, %v", file, err)
		}
		tout := fmt.Sprintf("%s.html", ser.Base)
		tname := "sermon.html"
		tkids := false
		if ser.Kids {
			tname = "kids.html"
		} else if i+1 < len(files) && strings.HasPrefix(files[i+1], ser.Base) {
			tkids = true
			ser.Kids = true
		}

		err = writeTemplate(tmpl, filepath.Join(out, tout), tname, ser)
		if err != nil {
			return fmt.Errorf("could not write template %s, %v", tout, err)
		}
		if !ser.Kids || tkids {
			infos = append(infos, ser.Info)
		} else if len(infos) > 0 {
			if lst := infos[len(infos)-1]; lst.ID == ser.ID {
				lst.Kids = true
			}
		}
	}

	tout := filepath.Join(out, "index.html")
	err = writeTemplate(tmpl, tout, "index.html", infos)
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
