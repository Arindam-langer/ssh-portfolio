package tui

import "charm.land/lipgloss/v2"

// Styles holds all the lipgloss styles derived from the active theme
type Styles struct {
	// Layout
	App         lipgloss.Style
	Header      lipgloss.Style
	TabBar      lipgloss.Style
	TabActive   lipgloss.Style
	TabInactive lipgloss.Style
	Content     lipgloss.Style
	Footer      lipgloss.Style
	StatusBar   lipgloss.Style

	// Typography
	Title     lipgloss.Style
	Subtitle  lipgloss.Style
	Heading   lipgloss.Style
	Body      lipgloss.Style
	Muted     lipgloss.Style
	Bold      lipgloss.Style
	Accent    lipgloss.Style
	Link      lipgloss.Style
	Tag       lipgloss.Style

	// Components
	Card        lipgloss.Style
	CardTitle   lipgloss.Style
	SkillBar    lipgloss.Style
	SkillFill   lipgloss.Style
	SkillEmpty  lipgloss.Style
	Bullet      lipgloss.Style
	Divider     lipgloss.Style
	Logo        lipgloss.Style
	Splash      lipgloss.Style
	HelpKey     lipgloss.Style
	HelpDesc    lipgloss.Style
}

func newStyles(t Theme) Styles {
	return Styles{
		// Layout
		App: lipgloss.NewStyle().
			Background(lipgloss.Color(t.Background)),

		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(t.Primary)).
			PaddingLeft(2).
			PaddingRight(2),

		TabBar: lipgloss.NewStyle().
			PaddingLeft(2).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(t.Border)),

		TabActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Primary)).
			Bold(true).
			Background(lipgloss.Color(t.Surface)).
			Padding(0, 2).
			BorderBottom(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color(t.Primary)),

		TabInactive: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.TextDim)).
			Padding(0, 2).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(t.Border)),

		Content: lipgloss.NewStyle().
			Padding(1, 3),

		Footer: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.TextDim)).
			PaddingLeft(2),

		StatusBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Text)).
			Background(lipgloss.Color(t.Surface)).
			Padding(0, 2),

		// Typography
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Primary)).
			Bold(true),

		Subtitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Secondary)).
			Italic(true),

		Heading: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Accent)).
			Bold(true).
			PaddingBottom(1),

		Body: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Text)),

		Muted: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.TextDim)),

		Bold: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Text)).
			Bold(true),

		Accent: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Accent)),

		Link: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Accent)).
			Underline(true),

		Tag: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Background)).
			Background(lipgloss.Color(t.Primary)).
			Padding(0, 1).
			Bold(true),

		// Components
		Card: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(t.Border)).
			Padding(1, 2).
			MarginBottom(1),

		CardTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Primary)).
			Bold(true),

		SkillBar: lipgloss.NewStyle(),

		SkillFill: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Primary)),

		SkillEmpty: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Muted)),

		Bullet: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Secondary)),

		Divider: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Border)),

		Logo: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Primary)).
			Bold(true),

		Splash: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Accent)),

		HelpKey: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Primary)).
			Bold(true),

		HelpDesc: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.TextDim)),
	}
}
