package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config map[string]*EntityConfig

type EntityConfig struct {
	Ignore bool                    `json:"ignore,omitempty"`
	Fields map[string]*FieldConfig `json:"fields,omitempty"`
}

type FieldConfig struct {
	Association bool   `json:"association,omitempty"`
	Array       bool   `json:"array,omitempty"`
	Ignore      bool   `json:"ignore,omitempty"`
	Rename      string `json:"rename,omitempty"`

	Opt Options `json:"opt,omitempty"`
}

type Options struct {
	Error     *float64 `json:"error,omitempty"`
	Tolerance *float64 `json:"tolerance,omitempty"`
	Round     string   `json:"round,omitempty"`
}

func ReadConfigFromFile(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer file.Close()

	fileType := filepath.Ext(path)
	switch fileType {
	case ".json":
		return readJSONConfigFromReader(file)
	case ".yaml":
		return readYAMConfigFromReader(file)
	default:
		return nil, fmt.Errorf("unsupported config file type %s", fileType)
	}
}

func readJSONConfigFromReader(reader io.Reader) (Config, error) {
	var config Config
	if err := json.NewDecoder(reader).Decode(&config); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return config, nil
}

func readYAMConfigFromReader(reader io.Reader) (Config, error) {
	var config Config
	if err := yaml.NewDecoder(reader).Decode(&config); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return config, nil
}
