package blog

import (
	"encoding/json"
	"io/ioutil"
)

// Config is the blog config
type Config struct {
	Logo         string `json:"logo"`
	Title        string `json:"title"`
	Path         string `json:"path"`
	TemplatePath string `json:"templatePath"`
	Listen       string `json:"listen"`
}

// DefaultConfig is the default config
func DefaultConfig() *Config {
	return &Config{
		Logo: "BL<br/>OG",
		Title: "Git based blogging",
		Path: "example.git",
		TemplatePath: "templates",
		Listen: ":80",
	}
}

// OpenConfig opens a config
func OpenConfig(file string) (*Config, error) {
	config := DefaultConfig()

	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
