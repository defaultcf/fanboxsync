package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type config struct {
	Default  *profileConfig
	Profiles []*profileConfig
}

type profileConfig struct {
	CreatorId string `yaml:"creator_id"`
	SessionId string `yaml:"session_id"`
}

func newConfig() (*config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configFilePath := filepath.Join(home, ".config", "fanboxsync", "config.yaml")
	f, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return loadConfig(f)
}

func loadConfig(r io.Reader) (*config, error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var profiles config
	err = yaml.Unmarshal(bytes, &profiles)
	if err != nil {
		return nil, err
	}

	defaultProfile := profiles.Default
	if defaultProfile == nil {
		defaultProfile = &profileConfig{}
	}
	// TODO: default 以外のキーを profiles にまとめているが、これを default と同じ階層にフラットに置けるようにしたい

	return &config{
		Default:  defaultProfile,
		Profiles: []*profileConfig{},
	}, nil
}
