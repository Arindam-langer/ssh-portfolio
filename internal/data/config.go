package data

import (
	"os"
	"strings"

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
	Items []Skill `yaml:"items"`
}
type TimelineItem struct {
	Title    string   `yaml:"title"`
	Subtitle string   `yaml:"subtitle"` // e.g. Company or Degree
	Period   string   `yaml:"period"`
	Location string   `yaml:"location"`
	Bullets  []string `yaml:"bullets"`
	Tech     string   `yaml:"tech"` // Optional string of tech stack
}


type KeyValueItem struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	AppConfig = &cfg
	populateFromConfig(&cfg)
	return nil
}

func populateFromConfig(cfg *Config) {
	if cfg.Profile.AsciiLogo != "" {
		Logo = cfg.Profile.AsciiLogo
	}
	if cfg.Profile.Tagline != "" {
		Tagline = cfg.Profile.Tagline
	}
	if cfg.Profile.Name != "" {
		PersonalInfo.Name = cfg.Profile.Name
	}

	for _, sec := range cfg.Sections {
		switch sec.Title {
		case "About":
			AboutText = sec.Content
		case "Skills":
			SkillCategories = make(map[string][]Skill)
			SkillCategoryOrder = []string{}
			for _, cat := range sec.Categories {
				SkillCategoryOrder = append(SkillCategoryOrder, cat.Name)
				SkillCategories[cat.Name] = cat.Items
			}
		case "Experience":
			Experiences = []Experience{}
			for _, item := range sec.TimelineItems {
				Experiences = append(Experiences, Experience{
					Title:    item.Title,
					Company:  item.Subtitle,
					Location: item.Location,
					Period:   item.Period,
					Bullets:  item.Bullets,
					Tech:     item.Tech,
				})
			}
		case "Projects":
			Projects = sec.Projects
		case "Education":
			if len(sec.TimelineItems) > 0 {
				item := sec.TimelineItems[0]
				gpa := strings.TrimSpace(strings.TrimPrefix(item.Tech, "GPA:"))
				EducationInfo = Education{
					Degree:     item.Title,
					School:     item.Subtitle,
					Location:   item.Location,
					Period:     item.Period,
					GPA:        gpa,
					Coursework: item.Bullets,
				}
			}
		case "Contact":
			for _, kv := range sec.KeyValueItems {
				switch kv.Key {
				case "Email":
					PersonalInfo.Email = kv.Value
				case "Phone":
					PersonalInfo.Phone = kv.Value
				case "LinkedIn":
					PersonalInfo.LinkedIn = kv.Value
				case "GitHub":
					PersonalInfo.GitHub = kv.Value
				case "Location":
					PersonalInfo.Location = kv.Value
				}
			}
		}
	}
}
