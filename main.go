package main

import (
	"fmt"
	"log"
	"os"
)

var target string

func init() {
	if len(os.Args) != 2 {
		fmt.Println("usage: subdumper <target>")
		os.Exit(1)
	}
	target = os.Args[1]
}

// in:
//   - {target}.lst
//
// out:
//   - {target}.files.txt
//   - {target}.files/
//   - {target}.sources.json
//   - {target}.sources/
//   - {target}.subroutines.json
func main() {
	log.Printf("starting subdumper for %s", target)
	dumper := NewDumper(target)
	log.Printf("clearing old files for %s", target)
	os.Remove(target + ".files.txt")
	os.RemoveAll(target + ".files")
	os.Remove(target + ".sources.json")
	os.RemoveAll(target + ".sources")
	os.Remove(target + ".subroutines.json")
	log.Printf("starting parse subroutines for %s", target)
	if err := dumper.PraseSubroutines(); err != nil {
		panic(err)
	}
	log.Printf("starting dump files for %s", target)
	if err := dumper.DumpFiles(); err != nil {
		panic(err)
	}
	log.Printf("starting dump sources for %s", target)
	if err := dumper.DumpSources(); err != nil {
		panic(err)
	}
	log.Printf("starting dump subroutines for %s", target)
	if err := dumper.DumpSubroutines(); err != nil {
		panic(err)
	}
	log.Printf("finished subdumper for %s", target)
}
