package tui

// Theme defines a color palette for the TUI using hex strings
type Theme struct {
	Name       string
	Primary    string
	Secondary  string
	Accent     string
	Muted      string
	Text       string
	TextDim    string
	Background string
	Surface    string
	Success    string
	Warning    string
	Error      string
	Border     string
	Highlight  string
}

// DefaultTheme uses the custom Arch palette specified by the user
var DefaultTheme = Theme{
	Name:       "Arch",
	Primary:    "#7EBAB5", // Accent / Main
	Secondary:  "#96CBC7", // Soft secondary accent
	Accent:     "#7EBAB5", // Main accent
	Muted:      "#454864", // Sub / Muted
	Text:       "#F6F5F5", // Crisp main text
	TextDim:    "#454864", // Sub / Muted text
	Background: "#0C0D11", // Deep dark background
	Surface:    "#171A25", // Elevated surface
	Success:    "#7EBAB5",
	Warning:    "#E2B714",
	Error:      "#CA4754",
	Border:     "#25293A", // Subtle container border
	Highlight:  "#7EBAB5",
}
