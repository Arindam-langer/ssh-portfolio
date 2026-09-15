package data

import (
	"os"

	"gopkg.in/yaml.v3"
)

var AppConfig *Config

type Config struct {
	Profile  Profile   `yaml:"profile"`
	Sections []Section `yaml:"sections"`
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

var DefaultAsciiLogo = `
   █████╗ ██████╗ ██╗███╗   ██╗██████╗  █████╗ ███╗   ███╗
  ██╔══██╗██╔══██╗██║████╗  ██║██╔══██╗██╔══██╗████╗ ████║
  ███████║██████╔╝██║██╔██╗ ██║██║  ██║███████║██╔████╔██║
  ██╔══██║██╔══██╗██║██║╚██╗██║██║  ██║██╔══██║██║╚██╔╝██║
  ██║  ██║██║  ██║██║██║ ╚████║██████╔╝██║  ██║██║ ╚═╝ ██║
  ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝
`

func LoadConfig(path string) error {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var cfg Config
	if err := yaml.Unmarshal(fileData, &cfg); err != nil {
		return err
	}
	if cfg.Profile.AsciiLogo == "" {
		cfg.Profile.AsciiLogo = DefaultAsciiLogo
	}
	AppConfig = &cfg
	return nil
}

