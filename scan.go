package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Ingester struct {
	Base, Raw, Prep, Text, Conf string
}

func newIngester() *Ingester {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("no home dir: %v", err)
	}
	base := filepath.Join(home, "predigt_scan")
	return &Ingester{
		Base: base,
		Raw:  filepath.Join(base, "raw"),
		Prep: filepath.Join(base, "prep"),
		Text: filepath.Join(base, "text"),
		Conf: filepath.Join(base, "scan.conf"),
	}
}

func (ing *Ingester) scanPrep(prefix, cmd string) error {
	now := time.Now().Format("060102-150405")
	name := fmt.Sprintf("%s-%s", prefix, now)
	err := ing.scan(name, cmd)
	if err != nil {
		return err
	}
	return ing.prep(name)
}

func (ing *Ingester) prep(prefix string) error {
	files, err := files(ing.Raw, prefix)
	if err != nil {
		return err
	}
	ensure(ing.Prep)
	for _, path := range files {
		err := ing.prepFile(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ing *Ingester) prepFile(path string) error {
	fname := filepath.Base(path)
	name := splitLast(fname, ".")
	log.Printf("converting %s\n", name)
	dest := filepath.Join(ing.Prep, fmt.Sprintf("%s.png", name))
	out, err := exec.Command("convert", path,
		"-deskew", "40%", "-despeckle", "-fuzz", "4%", "+repage",
		"-quality", "40", dest,
	).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("failed to convert %s to %s:\n%s", path, dest, out)
		// send done but leave file to be picked up again
	} else {
		err = os.Remove(path)
	}
	return err
}

func (ing *Ingester) ocr(prefix string) error {
	files, err := files(ing.Prep, prefix)
	if err != nil {
		return err
	}
	var last = "-"
	var batch []string
	for _, path := range files {
		if !strings.HasPrefix(path, last) {
			if len(batch) > 0 {
				err := ing.ocrFiles(batch...)
				if err != nil {
					return err
				}
				batch = batch[:0]
			}
			dir, file := filepath.Split(path)
			name := splitFirst(file, "-")
			last = filepath.Join(dir, name)
		}
		batch = append(batch, path)
	}
	if len(batch) > 0 {
		return ing.ocrFiles(batch...)
	}
	return nil
}
func (ing *Ingester) ocrFiles(files ...string) error {
	var buf bytes.Buffer
	name := ""
	for _, path := range files {
		if name == "" {
			name = splitFirst(filepath.Base(path), "-.")
		}
		buf.WriteString(path)
		buf.WriteByte('\n')
	}
	ensure(ing.Text)
	output := filepath.Join(ing.Text, name)
	log.Printf("detecting text %s\n", name)
	cmd := exec.Command("tesseract", "stdin", output,
		"--dpi", "300",
		"--user-words", "/home/mb0/work/predigt/spell.utf-8.add",
		"-l", "deu", // set language
		"--psm", "1",
		"quiet",
	)
	cmd.Stdin = &buf
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ocr failed: %v\n%s", err, out)
	}
	return nil
}

func (ing *Ingester) scan(prefix, cmd string) error {
	conf, fallback := ing.readConf(), true
	dev := conf.Last
	if dev == "" {
		dev, fallback = discover(), false
	}
	if dev == "" {
		return fmt.Errorf("no device found")
	}
	ensure(ing.Raw)
	args := make([]string, 0, 16)
	{
		args = append(args, "-d", "")
		var src string
		if cmd == "duplex" && conf.Duplex != "" {
			src = conf.Duplex
		} else if cmd != "scan" {
			src = "ADF"
			if conf.ADF != "" {
				src = conf.ADF
			}
		} else if conf.Scan != "" {
			src = conf.Scan
		}
		if src != "" {
			args = append(args, "--source", src)
		}
		mode := "Color"
		if conf.Mode != "" {
			mode = conf.Mode
		}
		args = append(args, "--mode", mode, "--resolution", "300", "-x", "210", "-y", "297")
		args = append(args, "--format", "tiff")
		if cmd == "scan" {
			args = append(args, "--batch-count", "1")
		}
		doc := filepath.Join(ing.Raw, fmt.Sprintf("%s-%%d.tif", prefix))
		args = append(args, fmt.Sprintf("--batch=%s", doc))
	}
	log.Printf("scaning %s\n", prefix)
	for dev != "" {
		args[1] = dev
		c := exec.Command("scanimage", args...)
		// TODO read output for progress and report to std out in a way that the gui can pick it up
		output, err := c.CombinedOutput()
		if err != nil {
			if fallback {
				// TODO check for device connection issues
				dev = discover()
				continue
			}
			return fmt.Errorf("scan failed: %v\n%s", err, output)
		}
		if dev != conf.Last {
			conf.Last = dev
			ing.writeConf(conf)
		}
		// inform or start prep daemon
		break
	}
	return nil
}

type config struct {
	Last   string `json:"last"`
	Scan   string `json:"scan,omitempty"`
	ADF    string `json:"adf,omitempty"`
	Duplex string `json:"duplex,omitempty"`
	Mode   string `json:"mode,omitempty"`
}

func (ing *Ingester) readConf() (conf config) {
	f, err := os.Open(ing.Conf)
	if err != nil {
		return conf
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&conf)
	if err != nil {
		log.Printf("could not read config %v", err)
	}
	return conf
}
func (ing *Ingester) writeConf(conf config) {
	f, err := os.Create(ing.Conf)
	if err != nil {
		log.Printf("could not create config file")
		return
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(conf)
	if err != nil {
		log.Printf("could not write config file")
	}

}

func splitLast(s, cut string) string {
	for i := 0; i < len(cut); i++ {
		if idx := strings.LastIndexByte(s, cut[i]); idx >= 0 {
			return s[:idx]
		}
	}
	return s
}

func splitFirst(s, cut string) string {
	for i := 0; i < len(cut); i++ {
		if idx := strings.IndexByte(s, cut[i]); idx >= 0 {
			return s[:idx]
		}
	}
	return s
}

func ensure(dirs ...string) {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0750)
		if err != nil {
			log.Fatalf("failed to create %s: %v", dir, err)
		}
	}
}

func files(dir, pref string) ([]string, error) {
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, f := range infos {
		if !strings.HasPrefix(f.Name(), pref) {
			continue
		}
		files = append(files, filepath.Join(dir, f.Name()))
	}
	return files, nil
}

func discover() string {
	res, err := exec.Command("scanimage", "-L").CombinedOutput()
	if err != nil {
		log.Printf("could not discover any device")
	}
	pref := []byte("device `")
	if !bytes.HasPrefix(res, pref) {
		return ""
	}
	res = res[len(pref):]
	idx := bytes.IndexByte(res, '\'')
	if idx < 0 {
		return ""
	}
	return string(res[:idx])
}
