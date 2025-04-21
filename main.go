package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/1l0/vv2srt/construct"
)

var (
	outputFilename string
	isAivis        bool = false

	// TEMP: adjustment for the sync problem
	adjustmentNanoSec = 0.0
)

func init() {
	log.SetFlags(log.Lshortfile)
	flag.StringVar(&outputFilename, "o", "", "output file name")
	flag.Float64Var(&adjustmentNanoSec, "d", 0, "duration adjustment: nano seconds per item")
	flag.Parse()
}

func main() {
	args := flag.Args()
	if len(args) < 1 {
		log.Fatalln(fmt.Errorf("missing project file path (.vvproj or .aisp)"))
	}
	projPath := args[0]
	exp, err := regexp.Compile(`.+\.(aisp|vvproj)$`)
	if err != nil {
		log.Fatalln(err)
	}
	m := exp.FindStringSubmatch(projPath)
	if len(m) < 2 {
		log.Fatalln(fmt.Errorf("unsupported project file"))
	}
	if outputFilename == "" {
		outputFilename = m[0] + ".srt"
	}
	if m[1] == "aisp" {
		isAivis = true
	}

	//  TEMP: adjustment for the sync problem
	if adjustmentNanoSec == 0.0 {
		if isAivis {
			adjustmentNanoSec = 60000000.0
		} else {
			adjustmentNanoSec = -30000000.0
		}
	}

	sub, err := construct.Project2subtitles(projPath, adjustmentNanoSec, isAivis)
	if err != nil {
		log.Fatalln(err)
	}
	f, err := os.OpenFile(outputFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	if err := f.Truncate(0); err != nil {
		log.Fatalln(err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		log.Fatalln(err)
	}
	if _, err := f.WriteString(sub.AsSRT()); err != nil {
		log.Fatalln(err)
	}
}
