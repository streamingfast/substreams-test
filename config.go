package subcmp

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	GraphURL         string                 `json:"graph_url"`
	TestOutputPath   string                 `json:"test_output_path"`
	SubstreamsModule string                 `json:"substreams_module"`
	Tests            map[string]*TestConfig `json:"tests"`
}

type TestConfig struct {
	Paths []*Path             `json:"paths"`
	Vars  []map[string]string `json:"vars"`
}

type Path struct {
	Graph      string `json:"graph"`
	Substreams string `json:"substreams"`
}

func ReadConfigFromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer file.Close()

	return ReadConfigFromReader(file)
}

func ReadConfigFromReader(reader io.Reader) (*Config, error) {
	var config *Config
	if err := json.NewDecoder(reader).Decode(&config); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return config, nil
}
