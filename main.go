package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-audio/wav"
	"github.com/martinlindhe/subtitles"
)

var outputName string

func init() {
	log.SetFlags(log.Lshortfile)
	flag.StringVar(&outputName, "o", "subtitle.srt", "output file name")
	flag.Parse()
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	if len(flag.Args()) > 0 {
		dir = flag.Args()[0]
	}

	files, err := collectFiles(dir)
	if err != nil {
		log.Fatalln(err)
	}
	sortFiles(files)

	captions := []subtitles.Caption{}
	offset := makeTime(0, 0, 0, 0)

	for _, file := range files {
		f, err := os.Open(file.AudioFilePath)
		if err != nil {
			log.Fatalln(err)
		}
		dur, err := wav.NewDecoder(f).Duration()
		if err != nil {
			log.Fatalln(err)
		}
		if err := f.Close(); err != nil {
			log.Fatalln(err)
		}

		f2, err := os.Open(file.TextFilePath)
		if err != nil {
			log.Fatalln(err)
		}
		b, err := io.ReadAll(f2)
		if err != nil {
			log.Fatalln(err)
		}
		if err := f2.Close(); err != nil {
			log.Fatalln(err)
		}

		extended := offset.Add(dur)

		captions = append(captions, subtitles.Caption{
			Seq:   file.Seq,
			Start: offset,
			End:   extended,
			Text:  []string{string(b)},
		})
		offset = extended
	}

	if len(captions) < 1 {
		log.Fatalln(fmt.Errorf("captions not found"))
	}

	s := subtitles.Subtitle{
		Captions: captions,
	}

	target := filepath.Join(dir, outputName)
	f, err := os.OpenFile(target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	if _, err := f.WriteString(s.AsSRT()); err != nil {
		log.Fatalln(err)
	}
}

func makeTime(h int, m int, s int, ms int) time.Time {
	return time.Date(0, 1, 1, h, m, s, ms*1000*1000, time.UTC)
}
