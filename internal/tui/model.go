package tui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Arindam-Langer/ssh-portfolio/internal/data"
)

// Tab indices
const (
	tabAbout = iota
	tabSkills
	tabExperience
	tabProjects
	tabEducation
	tabContact
	tabCount
)

var tabNames = []string{
	"  About",
	"  Skills",
	"  Experience",
	"  Projects",
	"  Education",
	"  Contact",
}

// Phase represents the app lifecycle
type phase int

const (
	phaseSplash phase = iota
	phaseMain
)

// tickMsg drives the splash animation
type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Model is the root Bubble Tea model
type Model struct {
	// Layout
	width  int
	height int
	term   string

	// State
	phase           phase
	activeTab       int
	styles          Styles
	showHelp        bool
	splashStep      int
	splashText      string
	splashCharIdx   int
	scrollOffset    int

	// Sub-models
	projectExpanded []bool
}

// NewModel creates a fresh portfolio model using Tokyo Night theme
func NewModel(term string, width, height int) Model {
	m := Model{
		width:           width,
		height:          height,
		term:            term,
		phase:           phaseSplash,
		activeTab:       tabAbout,
		styles:          newStyles(DefaultTheme),
		projectExpanded: make([]bool, len(data.Projects)),
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		if m.phase == phaseSplash {
			return m.updateSplash()
		}
		return m, nil

	case tea.KeyMsg:
		// Global keys
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "?":
			if m.phase == phaseMain {
				m.showHelp = !m.showHelp
			}
			return m, nil
		}

		if m.phase == phaseSplash {
			// Any key skips splash
			m.phase = phaseMain
			return m, nil
		}

		if m.showHelp {
			m.showHelp = false
			return m, nil
		}

		return m.handleMainKeys(msg)
	}
	return m, nil
}

func (m Model) handleMainKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "tab", "l", "right":
		m.activeTab = (m.activeTab + 1) % tabCount
		m.scrollOffset = 0
		return m, nil
	case "shift+tab", "h", "left":
		m.activeTab = (m.activeTab - 1 + tabCount) % tabCount
		m.scrollOffset = 0
		return m, nil
	case "j", "down":
		m.scrollOffset++
		return m, nil
	case "k", "up":
		if m.scrollOffset > 0 {
			m.scrollOffset--
		}
		return m, nil
	case "enter":
		if m.activeTab == tabProjects {
			return m.toggleProject()
		}
		return m, nil
	case "1":
		m.activeTab = tabAbout
		m.scrollOffset = 0
	case "2":
		m.activeTab = tabSkills
		m.scrollOffset = 0
	case "3":
		m.activeTab = tabExperience
		m.scrollOffset = 0
	case "4":
		m.activeTab = tabProjects
		m.scrollOffset = 0
	case "5":
		m.activeTab = tabEducation
		m.scrollOffset = 0
	case "6":
		m.activeTab = tabContact
		m.scrollOffset = 0
	}
	return m, nil
}

func (m Model) toggleProject() (Model, tea.Cmd) {
	// Cycle through projects: find next expandable based on scroll position
	idx := m.scrollOffset % len(data.Projects)
	m.projectExpanded[idx] = !m.projectExpanded[idx]
	return m, nil
}

func (m Model) updateSplash() (Model, tea.Cmd) {
	if m.splashStep >= len(data.SplashFrames) {
		m.phase = phaseMain
		return m, nil
	}

	target := data.SplashFrames[m.splashStep]
	if m.splashCharIdx < len(target) {
		m.splashCharIdx++
		m.splashText = target[:m.splashCharIdx]
		return m, tickCmd()
	}

	// Move to next frame
	m.splashStep++
	m.splashCharIdx = 0
	if m.splashStep < len(data.SplashFrames) {
		m.splashText = ""
	}
	return m, tickCmd()
}

func (m Model) View() tea.View {
	var content string
	if m.phase == phaseSplash {
		content = m.viewSplash()
	} else if m.showHelp {
		content = m.viewHelp()
	} else {
		content = m.viewMain()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// ============================================================================
// Splash Screen
// ============================================================================

func (m Model) viewSplash() string {
	s := m.styles

	logo := s.Logo.Render(data.Logo)
	tagline := s.Subtitle.Render(data.Tagline)

	// Build the progress display
	var lines []string
	for i := 0; i < m.splashStep; i++ {
		lines = append(lines, s.Muted.Render("  ✓ "+data.SplashFrames[i]))
	}
	if m.splashStep < len(data.SplashFrames) {
		cursor := "▊"
		lines = append(lines, s.Splash.Render("  → "+m.splashText+cursor))
	}

	progress := strings.Join(lines, "\n")

	block := lipgloss.JoinVertical(lipgloss.Center,
		logo,
		"",
		tagline,
		"",
		progress,
		"",
		s.Muted.Render("  press any key to skip..."),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, block)
}

// ============================================================================
// Main View
// ============================================================================

func (m Model) viewMain() string {
	header := m.renderHeader()
	tabs := m.renderTabs()
	footer := m.renderFooter()

	// Calculate content area height
	headerH := lipgloss.Height(header)
	tabsH := lipgloss.Height(tabs)
	footerH := lipgloss.Height(footer)
	contentH := m.height - headerH - tabsH - footerH - 2

	var body string
	switch m.activeTab {
	case tabAbout:
		body = m.viewAbout()
	case tabSkills:
		body = m.viewSkills()
	case tabExperience:
		body = m.viewExperience()
	case tabProjects:
		body = m.viewProjects()
	case tabEducation:
		body = m.viewEducation()
	case tabContact:
		body = m.viewContact()
	}

	// Apply scroll offset
	bodyLines := strings.Split(body, "\n")
	offset := m.scrollOffset
	if offset > len(bodyLines)-1 {
		offset = len(bodyLines) - 1
	}
	if offset < 0 {
		offset = 0
	}
	visibleLines := bodyLines[offset:]
	if len(visibleLines) > contentH {
		visibleLines = visibleLines[:contentH]
	}
	body = strings.Join(visibleLines, "\n")

	content := m.styles.Content.
		Width(m.width - 6).
		Height(contentH).
		Render(body)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		tabs,
		content,
		footer,
	)
}

func (m Model) renderHeader() string {
	s := m.styles

	name := s.Title.Render("  " + data.PersonalInfo.Name)
	role := s.Subtitle.Render(" " + data.Tagline)

	themeBadge := s.Tag.Render(DefaultTheme.Name)

	left := lipgloss.JoinHorizontal(lipgloss.Center, name, "  ", role)

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(themeBadge) - 4
	if gap < 0 {
		gap = 0
	}

	return left + strings.Repeat(" ", gap) + themeBadge + "\n"
}

func (m Model) renderTabs() string {
	s := m.styles
	var tabs []string

	for i, name := range tabNames {
		if i == m.activeTab {
			tabs = append(tabs, s.TabActive.Render(name))
		} else {
			tabs = append(tabs, s.TabInactive.Render(name))
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Bottom, tabs...)
	return s.TabBar.Width(m.width).Render(row)
}

func (m Model) renderFooter() string {
	s := m.styles
	keys := []struct{ key, desc string }{
		{"←/→", "navigate"},
		{"↑/↓", "scroll"},
		{"?", "help"},
		{"q", "quit"},
	}

	var parts []string
	for _, k := range keys {
		parts = append(parts, s.HelpKey.Render(k.key)+" "+s.HelpDesc.Render(k.desc))
	}

	return "\n" + s.Footer.Render(strings.Join(parts, "  │  "))
}

// ============================================================================
// Help Overlay
// ============================================================================

func (m Model) viewHelp() string {
	s := m.styles

	title := s.Title.Render("⌨  Keyboard Shortcuts")

	bindings := []struct{ key, desc string }{
		{"tab / shift+tab", "Navigate between sections"},
		{"← → / h l", "Navigate between sections"},
		{"↑ ↓ / j k", "Scroll content"},
		{"1-6", "Jump to section directly"},
		{"enter", "Expand/collapse project details"},
		{"?", "Toggle this help overlay"},
		{"q / ctrl+c", "Quit the portfolio"},
	}

	var lines []string
	lines = append(lines, title, "")
	for _, b := range bindings {
		key := s.HelpKey.Width(20).Render(b.key)
		desc := s.HelpDesc.Render(b.desc)
		lines = append(lines, "  "+key+"  "+desc)
	}

	lines = append(lines, "", s.Muted.Render("  Press any key to close..."))

	card := s.Card.Width(60).Render(strings.Join(lines, "\n"))
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, card)
}

// ============================================================================
// About Tab
// ============================================================================

func (m Model) viewAbout() string {
	s := m.styles
	contentWidth := m.width - 10
	if contentWidth < 40 {
		contentWidth = 40
	}

	title := s.Heading.Render("  About Me")

	// Build a nice about card
	aboutCard := s.Card.Width(contentWidth).Render(
		s.Body.Width(contentWidth - 6).Render(data.AboutText),
	)

	// Quick stats
	stats := []struct{ label, value string }{
		{"Location", data.PersonalInfo.Location},
		{"Email", data.PersonalInfo.Email},
		{"GitHub", data.PersonalInfo.GitHub},
		{"LinkedIn", data.PersonalInfo.LinkedIn},
	}

	var statLines []string
	for _, st := range stats {
		statLines = append(statLines,
			fmt.Sprintf("  %s  %s",
				s.Accent.Width(12).Render(st.label),
				s.Body.Render(st.value),
			),
		)
	}

	statsCard := s.Card.Width(contentWidth).Render(
		s.Heading.Render("  Quick Info") + "\n" +
			strings.Join(statLines, "\n"),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		aboutCard,
		"",
		statsCard,
	)
}

// ============================================================================
// Skills Tab
// ============================================================================

func (m Model) viewSkills() string {
	s := m.styles
	contentWidth := m.width - 10
	if contentWidth < 40 {
		contentWidth = 40
	}

	title := s.Heading.Render("  Technical Skills")

	var sections []string
	sections = append(sections, title)

	for _, category := range data.SkillCategoryOrder {
		skills := data.SkillCategories[category]

		catTitle := s.Accent.Bold(true).Render("  " + category)
		var skillLines []string
		skillLines = append(skillLines, catTitle)

		for _, skill := range skills {
			bar := renderSkillBar(s, skill.Name, skill.Level, contentWidth-12)
			skillLines = append(skillLines, "  "+bar)
		}

		sections = append(sections,
			s.Card.Width(contentWidth).Render(strings.Join(skillLines, "\n")),
		)
	}

	return strings.Join(sections, "\n")
}

func renderSkillBar(s Styles, name string, level int, maxWidth int) string {
	nameWidth := 16
	barWidth := maxWidth - nameWidth - 8
	if barWidth < 10 {
		barWidth = 10
	}

	filled := level * barWidth / 100
	empty := barWidth - filled

	nameStr := s.Body.Width(nameWidth).Render(name)
	fillStr := s.SkillFill.Render(strings.Repeat("█", filled))
	emptyStr := s.SkillEmpty.Render(strings.Repeat("░", empty))
	pctStr := s.Muted.Render(fmt.Sprintf(" %d%%", level))

	return nameStr + fillStr + emptyStr + pctStr
}

// ============================================================================
// Experience Tab
// ============================================================================

func (m Model) viewExperience() string {
	s := m.styles
	contentWidth := m.width - 10
	if contentWidth < 40 {
		contentWidth = 40
	}

	title := s.Heading.Render("  Work Experience")

	var cards []string
	cards = append(cards, title)

	for i, exp := range data.Experiences {
		header := fmt.Sprintf("%s  %s",
			s.CardTitle.Render(exp.Title),
			s.Muted.Render("@ "+exp.Company),
		)
		meta := fmt.Sprintf("  %s  │  %s",
			s.Accent.Render(exp.Period),
			s.Muted.Render(exp.Location),
		)

		var bullets []string
		for _, b := range exp.Bullets {
			wrapped := wrapText(b, contentWidth-10)
			bullets = append(bullets, "  "+s.Bullet.Render("▸")+" "+s.Body.Width(contentWidth-10).Render(wrapped))
		}

		tech := s.Muted.Render("  Tech: ") + s.Accent.Render(exp.Tech)

		content := strings.Join([]string{
			header,
			meta,
			"",
			strings.Join(bullets, "\n"),
			"",
			tech,
		}, "\n")

		card := s.Card.Width(contentWidth).Render(content)
		cards = append(cards, card)

		// Add timeline connector between experiences
		if i < len(data.Experiences)-1 {
			connector := s.Muted.Render("       │")
			cards = append(cards, connector)
		}
	}

	return strings.Join(cards, "\n")
}

// ============================================================================
// Projects Tab
// ============================================================================

func (m Model) viewProjects() string {
	s := m.styles
	contentWidth := m.width - 10
	if contentWidth < 40 {
		contentWidth = 40
	}

	title := s.Heading.Render("  Projects")
	hint := s.Muted.Render("  Use ↑/↓ to scroll, enter to expand/collapse")

	var cards []string
	cards = append(cards, title, hint, "")

	for i, proj := range data.Projects {
		header := fmt.Sprintf("%s  %s",
			s.CardTitle.Render("  "+proj.Name),
			s.Subtitle.Render("— "+proj.Tagline),
		)
		tech := s.Accent.Render("  "+proj.Tech)
		link := s.Link.Render("  "+proj.GitHubURL)

		lines := []string{header, tech, link}

		expanded := i < len(m.projectExpanded) && m.projectExpanded[i]
		if expanded {
			lines = append(lines, "")
			for _, b := range proj.Bullets {
				wrapped := wrapText(b, contentWidth-10)
				lines = append(lines, "  "+s.Bullet.Render("▸")+" "+s.Body.Width(contentWidth-10).Render(wrapped))
			}
		} else {
			lines = append(lines, s.Muted.Render("  ▶ press enter to expand details..."))
		}

		card := s.Card.Width(contentWidth).Render(strings.Join(lines, "\n"))
		cards = append(cards, card)
	}

	return strings.Join(cards, "\n")
}

// ============================================================================
// Education Tab
// ============================================================================

func (m Model) viewEducation() string {
	s := m.styles
	contentWidth := m.width - 10
	if contentWidth < 40 {
		contentWidth = 40
	}

	title := s.Heading.Render("  Education")

	edu := data.EducationInfo
	content := fmt.Sprintf(
		"%s\n%s\n\n  %s  │  %s\n  %s  %s",
		s.CardTitle.Render("  "+edu.School),
		s.Subtitle.Render("  "+edu.Degree),
		s.Accent.Render(edu.Period),
		s.Muted.Render(edu.Location),
		s.Body.Render("CGPA: "),
		s.Accent.Bold(true).Render(edu.GPA),
	)

	card := s.Card.Width(contentWidth).Render(content)

	// Relevant Coursework Section
	var courseLines []string
	courseLines = append(courseLines, s.Heading.Render("  Relevant Coursework"))
	for _, c := range edu.Coursework {
		courseLines = append(courseLines, "  "+s.Bullet.Render("▸")+" "+s.Body.Render(c))
	}

	courseCard := s.Card.Width(contentWidth).Render(strings.Join(courseLines, "\n"))

	// Add a "highlights" section
	facts := []string{
		"Built OllamaChat (a Bubble Tea TUI) during college",
		"Shipped production microservices before graduating",
		"This portfolio itself runs on Bubble Tea + Wish",
	}

	var factLines []string
	factLines = append(factLines, s.Heading.Render("  Highlights"))
	for _, f := range facts {
		factLines = append(factLines, "  "+s.Bullet.Render("★")+" "+s.Body.Render(f))
	}

	factsCard := s.Card.Width(contentWidth).Render(strings.Join(factLines, "\n"))

	return lipgloss.JoinVertical(lipgloss.Left, title, card, "", courseCard, "", factsCard)
}

// ============================================================================
// Contact Tab
// ============================================================================

func (m Model) viewContact() string {
	s := m.styles
	contentWidth := m.width - 10
	if contentWidth < 40 {
		contentWidth = 40
	}

	title := s.Heading.Render("  Get In Touch")

	info := data.PersonalInfo
	contactLines := []string{
		fmt.Sprintf("  %s  %s", s.Accent.Width(12).Render("Email"), s.Body.Render(info.Email)),
		fmt.Sprintf("  %s  %s", s.Accent.Width(12).Render("Phone"), s.Body.Render(info.Phone)),
		fmt.Sprintf("  %s  %s", s.Accent.Width(12).Render("LinkedIn"), s.Link.Render(info.LinkedIn)),
		fmt.Sprintf("  %s  %s", s.Accent.Width(12).Render("GitHub"), s.Link.Render(info.GitHub)),
		fmt.Sprintf("  %s  %s", s.Accent.Width(12).Render("Location"), s.Body.Render(info.Location)),
	}

	contactCard := s.Card.Width(contentWidth).Render(
		s.Heading.Render("  Contact Info") + "\n\n" +
			strings.Join(contactLines, "\n"),
	)

	// Resume download info
	resumeCard := s.Card.Width(contentWidth).Render(
		s.Heading.Render("  Resume") + "\n\n" +
			s.Body.Render("  Download my resume using SCP:") + "\n\n" +
			s.Accent.Render("  scp -P 2222 localhost:resume ./arindam_resume.pdf") + "\n\n" +
			s.Muted.Render("  Or connect via SFTP to browse available files."),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		contactCard,
		"",
		resumeCard,
	)
}

// ============================================================================
// Utilities
// ============================================================================

func wrapText(text string, maxWidth int) string {
	if maxWidth <= 0 || len(text) <= maxWidth {
		return text
	}

	var result strings.Builder
	words := strings.Fields(text)
	lineLen := 0

	for i, word := range words {
		if i > 0 && lineLen+len(word)+1 > maxWidth {
			result.WriteString("\n    ")
			lineLen = 4
		} else if i > 0 {
			result.WriteString(" ")
			lineLen++
		}
		result.WriteString(word)
		lineLen += len(word)
	}
	return result.String()
}
