package data

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var AppConfig *Config

type Config struct {
	Profile      Profile   `yaml:"profile"`
	SplashFrames []string  `yaml:"splash_frames,omitempty"`
	Sections     []Section `yaml:"sections"`
}

type Profile struct {
	Name      string `yaml:"name"`
	Tagline   string `yaml:"tagline"`
	AsciiLogo string `yaml:"ascii_logo"`
}

type Section struct {
	Title string `yaml:"title"`
	Icon  string `yaml:"icon,omitempty"`
	Type  string `yaml:"type"` // text, skill_list, timeline, projects, key_value

	// For type: text
	Content string `yaml:"content,omitempty"`

	// For type: skill_list
	Categories []SkillCategory `yaml:"categories,omitempty"`

	// For type: timeline
	TimelineItems []TimelineItem `yaml:"timeline_items,omitempty"`

	// For type: projects
	Projects []Project `yaml:"projects,omitempty"`

	// For type: key_value
	KeyValueItems []KeyValueItem `yaml:"key_value_items,omitempty"`
}

type SkillCategory struct {
	Name  string  `yaml:"name"`
	Icon  string  `yaml:"icon,omitempty"`
	Items []Skill `yaml:"items"`
}

type Skill struct {
	Name  string `yaml:"name"`
	Icon  string `yaml:"icon,omitempty"`
	Level int    `yaml:"level"` // 0-100
}

type TimelineItem struct {
	Title    string   `yaml:"title"`
	Subtitle string   `yaml:"subtitle"` // e.g. Company or Degree
	Icon     string   `yaml:"icon,omitempty"`
	Period   string   `yaml:"period"`
	Location string   `yaml:"location"`
	Bullets  []string `yaml:"bullets"`
	Tech     string   `yaml:"tech"` // Optional string of tech stack
}

type Project struct {
	Name      string   `yaml:"name"`
	Tagline   string   `yaml:"tagline"`
	Icon      string   `yaml:"icon,omitempty"`
	Tech      string   `yaml:"tech"`
	GitHubURL string   `yaml:"url"`
	Bullets   []string `yaml:"bullets"`
}

type KeyValueItem struct {
	Key   string `yaml:"key"`
	Icon  string `yaml:"icon,omitempty"`
	Value string `yaml:"value"`
}

var SplashFrames = []string{
	"Establishing secure connection...",
	"Loading portfolio modules...",
	"Initializing Bubble Tea runtime...",
	"Rendering terminal UI...",
	"Welcome aboard.",
}

func (c *Config) splashFramesOrDefault() {
	if len(c.SplashFrames) == 0 {
		c.SplashFrames = append([]string(nil), SplashFrames...)
	}
}

var DefaultAsciiLogo = `
▄▖▄▖▖▖  ▄▖    ▗ ▐▘  ▜ ▘  
▚ ▚ ▙▌  ▙▌▛▌▛▘▜▘▜▘▛▌▐ ▌▛▌
▄▌▄▌▌▌  ▌ ▙▌▌ ▐▖▐ ▙▌▐▖▌▙▌
                         
`

func (c *Config) Validate() error {
	if strings.TrimSpace(c.Profile.Name) == "" {
		return fmt.Errorf("'profile.name' is required in config.yaml")
	}
	if len(c.Sections) == 0 {
		return fmt.Errorf("at least one section is required under 'sections' in config.yaml")
	}
	for i, sec := range c.Sections {
		if strings.TrimSpace(sec.Title) == "" {
			return fmt.Errorf("section #%d is missing a 'title' in config.yaml", i+1)
		}
		if strings.TrimSpace(sec.Type) == "" {
			return fmt.Errorf("section '%s' is missing a 'type' in config.yaml", sec.Title)
		}
	}
	return nil
}

func LoadConfig(path string) error {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(fileData, &cfg); err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
	}
	if err := cfg.Validate(); err != nil {
		return err
	}
	if cfg.Profile.AsciiLogo == "" {
		cfg.Profile.AsciiLogo = DefaultAsciiLogo
	}
	cfg.splashFramesOrDefault()
	AppConfig = &cfg
	return nil
}
