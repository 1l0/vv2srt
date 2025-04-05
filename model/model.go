package model

type Project struct {
	AppVersion string `json:"appVersion"`
	Talk       struct {
		AudioKeys  []string               `json:"audioKeys"`
		AudioItems map[string]interface{} `json:"audioItems"`
	} `json:"talk"`
}

type AudioItem struct {
	Text  string `json:"text"`
	Voice *struct {
		EngineId  string `json:"engineId"`
		SpeakerId string `json:"speakerId"`
		StyleId   int64  `json:"styleId"`
	} `json:"voice,omitempty"`
	Query struct {
		AccentPhrases      []AccentPhrase `json:"accentPhrases"`
		SpeedScale         float64        `json:"speedScale"`
		IntonationScale    float64        `json:"intonationScale"`
		TempoDynamicsScale float64        `json:"tempoDynamicsScale"`
		PitchScale         float64        `json:"pitchScale"`
		VolumeScale        float64        `json:"volumeScale"`
		PauseLengthScale   float64        `json:"pauseLengthScale"`
		PrePhonemeLength   float64        `json:"prePhonemeLength"`
		PostPhonemeLength  float64        `json:"postPhonemeLength"`
		OutputSamplingRate string         `json:"outputSamplingRate"`
		OutputStereo       bool           `json:"outputStereo"`
		Kana               string         `json:"kana"`
	} `json:"query"`
	PresetKey string `json:"presetKey,omitempty"`
}

type AccentPhrase struct {
	Moras           []Mora `json:"moras"`
	Accent          int64  `json:"accent"`
	IsInterrogative bool   `json:"isInterrogative"`
	PauseMora       *struct {
		VowelLength float64 `json:"vowelLength"`
		Vowel       string  `json:"vowel"`
	} `json:"pauseMora,omitempty"`
}

type Mora struct {
	Text            string  `json:"text"`
	Vowel           string  `json:"vowel"`
	VowelLength     float64 `json:"vowelLength"`
	Pitch           float64 `json:"pitch"`
	Consonant       string  `json:"consonant"`
	ConsonantLength float64 `json:"consonantLength"`
}
