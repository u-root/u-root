package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type netdiskBootConfig struct {
	ImgURL        string `json:"ImgURL"`
	KernelPrefix  string `json:"KernelPrefix"`
	InitramPrefix string `json:"InitramPrefix"`
	Args          string `json:"Args"`
	Device        string `json:"Device"`
}

func (n netdiskBootConfig) String() string {
	var s strings.Builder
	_, err := s.WriteString("ImgURL: " + n.ImgURL + "\n")
	if err != nil {
		return ""
	}
	_, err = s.WriteString("KernelPrefix: " + n.KernelPrefix + "\n")
	if err != nil {
		return ""
	}
	_, err = s.WriteString("InitramPredix: " + n.InitramPrefix + "\n")
	if err != nil {
		return ""
	}
	_, err = s.WriteString("Args: " + n.Args + "\n")
	if err != nil {
		return ""
	}
	_, err = s.WriteString("Device: " + n.Device + "\n")
	if err != nil {
		return ""
	}
	return s.String()
}

func loadConfig(file string) (*netdiskBootConfig, error) {
	var config netdiskBootConfig
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
