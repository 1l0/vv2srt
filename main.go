package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/1l0/voicevox2srt/model"
	"github.com/martinlindhe/subtitles"
)

const (
	sec2nanosec = 1000000000.0
)

var outputFilename string
var isAivis bool = false
var adjustmentNanoSec float64 = 0.0

func init() {
	log.SetFlags(log.Lshortfile)
	flag.StringVar(&outputFilename, "o", "subtitles.srt", "output file name")
	flag.Parse()
}

func main() {
	args := flag.Args()
	if len(args) < 1 {
		log.Fatalln(fmt.Errorf("missing project file path (.vvproj or .aisp)"))
	}
	projPath := args[0]
	exp, err := regexp.Compile(`.+\.(aisp)|(vvproj)$`)
	if err != nil {
		log.Fatalln(err)
	}
	m := exp.FindStringSubmatch(projPath)
	if len(m) < 3 {
		log.Fatalln(fmt.Errorf("unsupported project file"))
	}
	// temporal adjust due to a bug of AivisSpeech
	if m[1] == "aisp" {
		isAivis = true
		adjustmentNanoSec = 60000000.0
	}
	sub, err := parseSubtitles(projPath)
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

func parseSubtitles(projectFilePath string) (*subtitles.Subtitle, error) {
	readfile, err := os.Open(projectFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	dec := json.NewDecoder(readfile)
	proj := &model.Project{}
	if err := dec.Decode(proj); err != nil {
		return nil, err
	}
	if err := readfile.Close(); err != nil {
		return nil, err
	}
	captions := []subtitles.Caption{}
	zero := makeTime(0, 0, 0, 0)
	epoch := makeTime(0, 0, 0, 0)
	var offset float64 = 0.0

	for i, key := range proj.Talk.AudioKeys {
		if it, exist := proj.Talk.AudioItems[key]; exist {
			b, err := json.Marshal(it)
			if err != nil {
				return nil, err
			}
			item := &model.AudioItem{}
			if err := json.Unmarshal(b, item); err != nil {
				return nil, err
			}
			speedScale := item.Query.SpeedScale
			offset += item.Query.PrePhonemeLength * sec2nanosec / speedScale
			for _, acc := range item.Query.AccentPhrases {
				for _, mo := range acc.Moras {
					if mo.Consonant != "" {
						offset += (mo.ConsonantLength * sec2nanosec) / speedScale
					}
					offset += (mo.VowelLength * sec2nanosec) / speedScale
				}
				if acc.PauseMora != nil {
					offset += (acc.PauseMora.VowelLength * item.Query.PauseLengthScale * sec2nanosec) / speedScale
				}
			}
			offset += item.Query.PostPhonemeLength * sec2nanosec / speedScale
			offset += adjustmentNanoSec
			nextEpoch := zero.Add(time.Duration(offset))
			captions = append(captions, subtitles.Caption{
				Seq:   i,
				Start: epoch,
				End:   nextEpoch,
				Text:  []string{item.Text},
			})
			epoch = nextEpoch
		} else {
			return nil, fmt.Errorf("audio item not found for the key: %s", key)
		}
	}
	fmt.Printf("total time: %d:%d:%d:%d\n", epoch.Hour(), epoch.Minute(), epoch.Second(), epoch.Nanosecond())

	return &subtitles.Subtitle{
		Captions: captions,
	}, nil
}

func makeTime(h int, m int, s int, ms int) time.Time {
	return time.Date(0, 1, 1, h, m, s, ms*1000*1000, time.UTC)
}
