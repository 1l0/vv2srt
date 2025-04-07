package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/1l0/vv2srt/model"
	"github.com/martinlindhe/subtitles"
)

const (
	sec2nanoSec = 1000000000.0
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
	proj, err := loadProject(projectFilePath)
	if err != nil {
		return nil, err
	}
	captions := []subtitles.Caption{}
	zero := makeTime(0, 0, 0, 0)
	epoch := zero
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
			offset += item.Query.PrePhonemeLength * sec2nanoSec / speedScale
			for _, acc := range item.Query.AccentPhrases {
				for _, mo := range acc.Moras {
					if mo.Consonant != "" {
						offset += (mo.ConsonantLength * sec2nanoSec) / speedScale
					}
					if mo.Vowel != "" {
						offset += (mo.VowelLength * sec2nanoSec) / speedScale
					}
				}
				if acc.PauseMora != nil && acc.PauseMora.Vowel != "" {
					offset += (acc.PauseMora.VowelLength * item.Query.PauseLengthScale * sec2nanoSec) / speedScale
				}
			}
			offset += item.Query.PostPhonemeLength * sec2nanoSec / speedScale

			// TEMP: adjustment for the sync problem
			offset += adjustmentNanoSec

			nextEpoch := zero.Add(time.Duration(offset))
			captions = append(captions, subtitles.Caption{
				Seq:   i + 1,
				Start: epoch,
				End:   nextEpoch,
				Text:  []string{item.Text},
			})
			epoch = nextEpoch
		} else {
			return nil, fmt.Errorf("audio item not found for the key: %s", key)
		}
	}
	platform := func() string {
		if isAivis {
			return "AivisSearch"
		}
		return "VOICEVOX"
	}()
	fmt.Printf(
		"Platform: %s\nDuration: %02d:%02d:%02d.%d\n",
		platform,
		epoch.Hour(),
		epoch.Minute(),
		epoch.Second(),
		epoch.Nanosecond(),
	)
	return &subtitles.Subtitle{
		Captions: captions,
	}, nil
}

func loadProject(projectFilePath string) (*model.Project, error) {
	readfile, err := os.Open(projectFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer readfile.Close()
	dec := json.NewDecoder(readfile)
	proj := &model.Project{}
	if err := dec.Decode(proj); err != nil {
		return nil, err
	}
	return proj, nil
}

func makeTime(h int, m int, s int, ms int) time.Time {
	return time.Date(0, 1, 1, h, m, s, ms*1000*1000, time.UTC)
}
