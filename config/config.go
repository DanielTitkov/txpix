package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	BackgroundImages []string         `yaml:"background_images"`
	BackgroundColor  string           `yaml:"background_color"`
	ImageWidth       int              `yaml:"image_width"`
	ImageHeight      int              `yaml:"image_height"`
	FontSize         float64          `yaml:"font_size"`
	FontColor        string           `yaml:"font_color"`
	FontFile         string           `yaml:"font_file"`
	Margin           int              `yaml:"margin"`
	OutputDir        string           `yaml:"output_dir"`
	LineSpacing      int              `yaml:"line_spacing"`
	Airtable         Airtable         `yaml:"airtable"`
	Preprocess       PreprocessConfig `yaml:"preprocess"`
}

type Airtable struct {
	APIKey    string `yaml:"api_key"`
	BaseID    string `yaml:"base_id"`
	TableName string `yaml:"table_name"`
	NameID    string `yaml:"name_id"`
	BodyID    string `yaml:"body_id"`
}

type PreprocessConfig struct {
	Remove string `yaml:"remove"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
