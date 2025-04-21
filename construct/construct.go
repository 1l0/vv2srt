package construct

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/1l0/vv2srt/model"
	"github.com/martinlindhe/subtitles"
)

const (
	sec2nanoSec = 1000000000.0
)

func Project2subtitles(projectFilePath string, adjustmentNanoSec float64, isAivis bool) (*subtitles.Subtitle, error) {
	proj, err := loadProject(projectFilePath)
	if err != nil {
		return nil, err
	}
	return generateSubtitles(proj, adjustmentNanoSec, isAivis)
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

func generateSubtitles(proj *model.Project, adjustmentNanoSec float64, isAivis bool) (*subtitles.Subtitle, error) {
	captions := []subtitles.Caption{}
	zero := makeTime(0, 0, 0, 0)
	epoch := zero
	var offset float64 = 0.0

	for i, key := range proj.Talk.AudioKeys {
		it, exist := proj.Talk.AudioItems[key]
		if !exist {
			return nil, fmt.Errorf("audio item not found for the key: %s", key)
		}
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

func makeTime(h int, m int, s int, ms int) time.Time {
	return time.Date(0, 1, 1, h, m, s, ms*1000*1000, time.UTC)
}
