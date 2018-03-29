package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type GMina struct {
	APIKey   string
	CacheDir string
}

func New(apiKey string, cacheDir string) (*GMina, error) {
	if !isFileExist(cacheDir) {
		err := os.MkdirAll(cacheDir, 0755)
		if err != nil {
			return nil, err
		}
	}
	return &GMina{
		APIKey:   apiKey,
		CacheDir: cacheDir,
	}, nil
}

func (g *GMina) TTS(text string) (string, error) {
	input := map[string]interface{}{
		"input": map[string]interface{}{
			"text": text,
		},
		"voice": map[string]interface{}{
			"languageCode": "en-US",
			"name":         "en-US-Wavenet-D",
		},
		"audioConfig": map[string]interface{}{
			"audioEncoding": "MP3",
			"pitch":         0.00,
			"speakingRate":  1.00,
		},
	}

	jsonBytes, _ := json.Marshal(input)

	req, err := http.NewRequest("POST", "https://texttospeech.googleapis.com/v1beta1/text:synthesize?key="+g.APIKey, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}

	resp, err := g.request(req)
	if err != nil {
		return "", err
	}

	jsonBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var output struct {
		AudioContent string `json:"audioContent"`
	}

	err = json.Unmarshal(jsonBytes, &output)
	return output.AudioContent, err
}
