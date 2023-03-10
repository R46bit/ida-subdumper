package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

type Dumper struct {
	target string
	filter int
	offset int64

	subroutines []*Subroutine

	sources map[string][]*Subroutine
}

func NewDumper(target string) *Dumper {
	return &Dumper{target: target, sources: make(map[string][]*Subroutine)}
}

type Subroutine struct {
	Line int64  `json:"line"`
	Addr string `json:"addr"`
	Func string `json:"func"`

	file string
	safe string
}

func (d *Dumper) PraseSubroutines() error {
	f, err := os.Open(target + ".lst")
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	line := scanner.Text()
	parts := strings.Split(line, ";")
	if len(parts) != 2 {
		return errors.New("failed to parse line length filter and offset")
	}
	d.filter = len(line)
	d.offset, err = strconv.ParseInt(strings.TrimSpace(parts[0][5:]), 16, 64)
	if err != nil {
		return err
	}
	if d.filter == 0 || d.offset == 0 {
		return errors.New("invalid line length filter or offset")
	}
	log.Printf("parsed line length filter: %d, offset: %d", d.filter, d.offset)

	for {
		subroutine := d.NextSubroutine(scanner)
		if subroutine == nil {
			break
		}
		d.subroutines = append(d.subroutines, subroutine)
	}
	for _, subroutine := range d.subroutines {
		d.sources[subroutine.safe] = append(d.sources[subroutine.safe], subroutine)
	}
	for _, subroutines := range d.sources {
		sort.Slice(subroutines, func(i, j int) bool {
			if subroutines[i].safe == subroutines[j].safe {
				if subroutines[i].Line == subroutines[j].Line {
					return subroutines[i].Func < subroutines[j].Func
				}
				return subroutines[i].Line < subroutines[j].Line
			}
			return subroutines[i].safe < subroutines[j].safe
		})
	}
	return nil
}

func (d *Dumper) NextSubroutine(scanner *bufio.Scanner) *Subroutine {
	var subroutine *Subroutine
	var addr, line string
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) < d.filter || text[:6] != ".text:" {
			continue
		}
		if strings.TrimSpace(text[d.filter:]) == "; =============== S U B R O U T I N E =======================================" {
			if text[6:22] != addr {
				continue
			}
			subroutine = &Subroutine{}
			if strings.HasPrefix(line, "#line ") {
				parts := strings.Split(line[7:], "\" ")
				if len(parts) == 2 {
					subroutine.file = parts[0]
					subroutine.safe = path.Clean(parts[0])
					subroutine.Line, _ = strconv.ParseInt(parts[1], 10, 64)
				} else {
					subroutine.Line, _ = strconv.ParseInt(line[6:], 10, 64)
				}
				subroutine.Addr = addr
			}
			line = ""
			continue
		} else if subroutine == nil {
			temp := text[6:22]
			text = strings.TrimSpace(text[d.filter:])
			if strings.HasPrefix(text, "; #line ") {
				addr = temp
				line = text[2:]
			}
			continue
		}
		if text[6:22] == addr {
			text = text[d.filter:]
			if strings.Contains(text, " proc near") {
				if line == "" || strings.HasPrefix(line, "Attributes:") {
					subroutine.Func = text[:strings.Index(text, " proc near")]
				} else {
					subroutine.Func = line
				}
			} else if strings.HasPrefix(text, "; ") {
				line = text[2:]
			}
		} else if strings.HasSuffix(text, " endp") {
			subroutine.Addr += "-" + text[6:22]
			break
		}
	}
	return subroutine
}

func (d *Dumper) DumpFiles() error {
	f, err := os.OpenFile(d.target+".files.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	files := make([]string, 0, len(d.sources))
	for file := range d.sources {
		files = append(files, file)
	}
	sort.Strings(files)
	for _, file := range files {
		if file == "" {
			continue
		}
		f.WriteString(file + "\n")
		subroutines := d.sources[file]
		name := d.target + ".files" + file
		if _, err := os.Stat(name); os.IsNotExist(err) {
			err = os.MkdirAll(path.Dir(name), 0700)
			if err != nil {
				log.Println("failed to mkdir:", name, err)
				continue
			}
		}
		f, err := os.Create(name)
		if err != nil {
			log.Println("failed to create file:", name, err)
			continue
		}
		f.WriteString(fmt.Sprintf("// File: %s\n", file))
		for _, subroutine := range subroutines {
			f.WriteString("\n")
			f.WriteString(fmt.Sprintf("// Line %d: range %s\n", subroutine.Line, subroutine.Addr))
			f.WriteString(fmt.Sprintf("%s;\n", subroutine.Func))
		}
		f.Close()
	}
	return nil
}

func (d *Dumper) DumpSources() error {
	f, err := os.OpenFile(d.target+".sources.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	return enc.Encode(d.sources)
}

func (d *Dumper) DumpSubroutines() error {
	f, err := os.OpenFile(d.target+".subroutines.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	return enc.Encode(d.subroutines)
}
