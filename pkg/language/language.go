package language

import (
	"encoding/json"
	"io/ioutil"
)

// ...
const (
	DefaultLanguageCode  = "en"
	JapaneseLanguageCode = "ja"
)

// LoadLanguages ...
func LoadLanguages(path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &Languages)
}

// Languages ...
var Languages []Language

// Language ...
type Language struct {
	Label      string       `json:"label"`
	Code       string       `json:"code"`
	SpeechCode []SpeechCode `json:"speech_code"`
}

// SpeechCode struct
type SpeechCode struct {
	Label string `json:"label"`
	Code  string `json:"code"`
}

// IsMultipleSpeechCode ...
func IsMultipleSpeechCode(code string) bool {
	for _, lang := range Languages {
		if lang.Code == code {
			if len(lang.SpeechCode) > 1 {
				return true
			}
		}
	}
	return false
}
