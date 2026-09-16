package tui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Arindam-Langer/ssh-portfolio/internal/data"
	"github.com/charmbracelet/bubbles/viewport"
)

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
	phase         phase
	activeTab     int
	styles        Styles
	showHelp      bool
	splashStep    int
	splashText    string
	splashCharIdx int
	scrollOffset  int
	viewport      viewport.Model
	// Sub-models / expansion tracking (secIdx_projIdx -> bool)
	projectExpanded map[string]bool
}

// NewModel creates a fresh portfolio model
func NewModel(term string, width, height int) Model {
	m := Model{
		width:           width,
		height:          height,
		term:            term,
		phase:           phaseSplash,
		activeTab:       0,
		styles:          newStyles(DefaultTheme),
		projectExpanded: make(map[string]bool),
		viewport:        viewport.New(width, height),
	}
	m.updateViewport()
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
		m.updateViewport()
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
	secCount := 0
	if data.AppConfig != nil {
		secCount = len(data.AppConfig.Sections)
	}

	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "tab", "l", "right":
		if secCount > 0 {
			m.activeTab = (m.activeTab + 1) % secCount
			m.scrollOffset = 0
			m.viewport.GotoTop()
			m.updateViewport()
		}
		return m, nil
	case "shift+tab", "h", "left":
		if secCount > 0 {
			m.activeTab = (m.activeTab - 1 + secCount) % secCount
			m.scrollOffset = 0
			m.viewport.GotoTop()
			m.updateViewport()
		}
		return m, nil
	case "j", "down":
		m.viewport.ScrollDown(1)
		// m.scrollOffset++ // causing infinite scroll need a fix.
		return m, nil
	case "k", "up":
		m.viewport.ScrollUp(1)
		// if m.scrollOffset > 0 {
		// 	m.scrollOffset--
		// }
		return m, nil
	case "enter":
		if secCount > 0 && m.activeTab < secCount {
			if data.AppConfig.Sections[m.activeTab].Type == "projects" {
				m, cmd := m.toggleProject()
				m.updateViewport()
				return m, cmd
			}
		}
		return m, nil
	case "1", "2", "3", "4", "5", "6", "7", "8", "9":
		idx := int(msg.String()[0] - '1')
		if idx < secCount {
			m.activeTab = idx
			m.scrollOffset = 0
			m.viewport.GotoTop()
			m.updateViewport()
		}
		return m, nil
	}
	return m, nil
}

func (m Model) toggleProject() (Model, tea.Cmd) {
	if data.AppConfig == nil || m.activeTab >= len(data.AppConfig.Sections) {
		return m, nil
	}
	sec := data.AppConfig.Sections[m.activeTab]
	if len(sec.Projects) == 0 {
		return m, nil
	}
	idx := m.scrollOffset % len(sec.Projects)
	key := fmt.Sprintf("%d_%d", m.activeTab, idx)
	m.projectExpanded[key] = !m.projectExpanded[key]
	return m, nil
}

func (m Model) updateSplash() (Model, tea.Cmd) {
	frames := data.SplashFrames
	if m.splashStep >= len(frames) {
		m.phase = phaseMain
		return m, nil
	}

	target := frames[m.splashStep]
	if m.splashCharIdx < len(target) {
		m.splashCharIdx++
		m.splashText = target[:m.splashCharIdx]
		return m, tickCmd()
	}

	// Move to next frame
	m.splashStep++
	m.splashCharIdx = 0
	if m.splashStep < len(frames) {
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

// Splash Screen

func (m Model) viewSplash() string {
	s := m.styles

	logoStr := data.DefaultAsciiLogo
	taglineStr := ""
	if data.AppConfig != nil {
		if data.AppConfig.Profile.AsciiLogo != "" {
			logoStr = data.AppConfig.Profile.AsciiLogo
		}
		taglineStr = data.AppConfig.Profile.Tagline
	}

	logo := s.Logo.Render(logoStr)
	tagline := s.Subtitle.Render(taglineStr)

	frames := data.SplashFrames

	// Build the progress display
	var lines []string
	for i := 0; i < m.splashStep && i < len(frames); i++ {
		lines = append(lines, s.Muted.Render("  ✓ "+frames[i]))
	}
	if m.splashStep < len(frames) {
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

func (m Model) renderActiveSection() string {
	var body string

	if data.AppConfig != nil && len(data.AppConfig.Sections) > 0 {
		activeIdx := m.activeTab
		if activeIdx >= len(data.AppConfig.Sections) {
			activeIdx = 0
		}
		sec := data.AppConfig.Sections[activeIdx]

		switch sec.Type {
		case "text":
			body = m.renderTextSection(sec)
		case "skill_list":
			body = m.renderSkillListSection(sec)
		case "timeline":
			body = m.renderTimelineSection(sec)
		case "projects":
			body = m.renderProjectsSection(sec)
		case "key_value":
			body = m.renderKeyValueSection(sec)
		default:
			body = m.renderFallbackSection(sec)
		}
	} else {
		body = m.styles.Muted.Render("No sections configured.")
	}

	return body
}

// updateViewport calculates the current terminal content area and updates
// the viewport model. This must run from Update, not View, because
// viewport.SetContent mutates the viewport state.
func (m *Model) updateViewport() {
	header := m.renderHeader()
	tabs := m.renderTabs()
	footer := m.renderFooter()

	headerH := lipgloss.Height(header)
	tabsH := lipgloss.Height(tabs)
	footerH := lipgloss.Height(footer)

	contentH := m.height - headerH - tabsH - footerH - 2
	if contentH < 5 {
		contentH = 5
	}

	contentWidth := m.width - 6
	if contentWidth < 20 {
		contentWidth = 20
	}

	// The viewport is rendered inside Content, which may have padding/borders.
	// Its usable size must exclude that style's frame, otherwise the viewport
	// thinks it can display more lines than the outer Content box can actually
	// show. That makes the final lines appear unreachable.
	viewportWidth := contentWidth - m.styles.Content.GetHorizontalFrameSize()
	viewportHeight := contentH - m.styles.Content.GetVerticalFrameSize()

	if viewportWidth < 1 {
		viewportWidth = 1
	}
	if viewportHeight < 1 {
		viewportHeight = 1
	}

	m.viewport.Width = viewportWidth
	m.viewport.Height = viewportHeight
	m.viewport.SetContent(m.renderActiveSection())
}

func (m Model) viewMain() string {
	header := m.renderHeader()
	tabs := m.renderTabs()
	footer := m.renderFooter()

	headerH := lipgloss.Height(header)
	tabsH := lipgloss.Height(tabs)
	footerH := lipgloss.Height(footer)

	contentH := m.height - headerH - tabsH - footerH - 2
	if contentH < 5 {
		contentH = 5
	}

	contentWidth := m.width - 6
	if contentWidth < 20 {
		contentWidth = 20
	}

	// Old manual scroll logic kept for reference:
	// bodyLines := strings.Split(body, "\n")
	// offset := m.scrollOffset
	// if offset > len(bodyLines)-1 {
	// 	offset = len(bodyLines) - 1
	// }
	// if offset < 0 {
	// 	offset = 0
	// }
	// visibleLines := bodyLines[offset:]
	// if len(visibleLines) > contentH {
	// 	visibleLines = visibleLines[:contentH]
	// }
	// body = strings.Join(visibleLines, "\n")

	// The viewport is updated in Update(). View() only reads its current state.
	body := m.viewport.View()

	content := m.styles.Content.
		Width(contentWidth).
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

	nameStr := ""
	taglineStr := ""
	if data.AppConfig != nil {
		nameStr = data.AppConfig.Profile.Name
		taglineStr = data.AppConfig.Profile.Tagline
	}

	name := s.Title.Render("  " + nameStr)
	role := s.Subtitle.Render(" " + taglineStr)

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

	if data.AppConfig != nil {
		for i, sec := range data.AppConfig.Sections {
			title := sec.Title
			if sec.Icon != "" {
				title = sec.Icon + " " + sec.Title
			}
			title = "  " + title
			if i == m.activeTab {
				tabs = append(tabs, s.TabActive.Render(title))
			} else {
				tabs = append(tabs, s.TabInactive.Render(title))
			}
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

// Help Overlay

func (m Model) viewHelp() string {
	s := m.styles

	title := s.Title.Render("⌨  Keyboard Shortcuts")

	secCount := 0
	if data.AppConfig != nil {
		secCount = len(data.AppConfig.Sections)
	}

	numKeyDesc := fmt.Sprintf("1-%d", secCount)
	if secCount == 0 {
		numKeyDesc = "1-N"
	}

	bindings := []struct{ key, desc string }{
		{"tab / shift+tab", "Navigate between sections"},
		{"← → / h l", "Navigate between sections"},
		{"↑ ↓ / j k", "Scroll content"},
		{numKeyDesc, "Jump to section directly"},
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

// Dynamic Section Renderers

func (m Model) renderTextSection(sec data.Section) string {
	s := m.styles
	contentWidth := m.width - 10
	if contentWidth < 40 {
		contentWidth = 40
	}

	titleText := sec.Title
	if sec.Icon != "" {
		titleText = sec.Icon + " " + titleText
	}
	title := s.Heading.Render("  " + titleText)

	card := s.Card.Width(contentWidth).Render(
		s.Body.Width(contentWidth - 6).Render(sec.Content),
	)

	return lipgloss.JoinVertical(lipgloss.Left, title, card)
}

func (m Model) renderSkillListSection(sec data.Section) string {
	s := m.styles
	contentWidth := m.width - 10
	if contentWidth < 40 {
		contentWidth = 40
	}

	titleText := sec.Title
	if sec.Icon != "" {
		titleText = sec.Icon + " " + titleText
	}
	title := s.Heading.Render("  " + titleText)

	var sections []string
	sections = append(sections, title)

	for _, category := range sec.Categories {
		catName := category.Name
		if category.Icon != "" {
			catName = category.Icon + " " + catName
		}
		catTitle := s.Accent.Bold(true).Render("  " + catName)
		var skillLines []string
		skillLines = append(skillLines, catTitle)

		for _, skill := range category.Items {
			name := skill.Name
			if skill.Icon != "" {
				name = skill.Icon + " " + name
			}
			if skill.Level > 0 {
				bar := renderSkillBar(s, name, skill.Level, contentWidth-12)
				skillLines = append(skillLines, "  "+bar)
			} else {
				line := fmt.Sprintf("  %s  %s", s.Bullet.Render("▸"), s.Body.Render(name))
				skillLines = append(skillLines, line)
			}
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

func (m Model) renderTimelineSection(sec data.Section) string {
	s := m.styles
	contentWidth := m.width - 10
	if contentWidth < 40 {
		contentWidth = 40
	}

	titleText := sec.Title
	if sec.Icon != "" {
		titleText = sec.Icon + " " + titleText
	}
	title := s.Heading.Render("  " + titleText)

	var cards []string
	cards = append(cards, title)

	for i, item := range sec.TimelineItems {
		itemTitle := item.Title
		if item.Icon != "" {
			itemTitle = item.Icon + " " + itemTitle
		}

		var header string
		if item.Subtitle != "" {
			header = fmt.Sprintf("%s  %s",
				s.CardTitle.Render("  "+itemTitle),
				s.Muted.Render("@ "+item.Subtitle),
			)
		} else {
			header = s.CardTitle.Render("  " + itemTitle)
		}

		var metaParts []string
		if item.Period != "" {
			metaParts = append(metaParts, s.Accent.Render(item.Period))
		}
		if item.Location != "" {
			metaParts = append(metaParts, s.Muted.Render(item.Location))
		}
		meta := "  " + strings.Join(metaParts, "  │  ")

		var lines []string
		lines = append(lines, header)
		if len(metaParts) > 0 {
			lines = append(lines, meta)
		}
		if len(item.Bullets) > 0 {
			lines = append(lines, "")
			for _, b := range item.Bullets {
				wrapped := wrapText(b, contentWidth-10)
				lines = append(lines, "  "+s.Bullet.Render("▸")+" "+s.Body.Width(contentWidth-10).Render(wrapped))
			}
		}
		if item.Tech != "" {
			lines = append(lines, "", s.Muted.Render("  Info/Tech: ")+s.Accent.Render(item.Tech))
		}

		card := s.Card.Width(contentWidth).Render(strings.Join(lines, "\n"))
		cards = append(cards, card)

		if i < len(sec.TimelineItems)-1 {
			connector := s.Muted.Render("       │")
			cards = append(cards, connector)
		}
	}

	return strings.Join(cards, "\n")
}

func (m Model) renderProjectsSection(sec data.Section) string {
	s := m.styles
	contentWidth := m.width - 10
	if contentWidth < 40 {
		contentWidth = 40
	}

	titleText := sec.Title
	if sec.Icon != "" {
		titleText = sec.Icon + " " + titleText
	}
	title := s.Heading.Render("  " + titleText)
	hint := s.Muted.Render("  Use ↑/↓ to scroll, enter to expand/collapse")

	var cards []string
	cards = append(cards, title, hint, "")

	for i, proj := range sec.Projects {
		projName := proj.Name
		if proj.Icon != "" {
			projName = proj.Icon + " " + projName
		}
		header := fmt.Sprintf("%s  %s",
			s.CardTitle.Render("  "+projName),
			s.Subtitle.Render("— "+proj.Tagline),
		)
		tech := s.Accent.Render("  " + proj.Tech)
		link := s.Link.Render("  " + proj.GitHubURL)

		lines := []string{header, tech, link}

		key := fmt.Sprintf("%d_%d", m.activeTab, i)
		expanded := m.projectExpanded[key]
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

func (m Model) renderKeyValueSection(sec data.Section) string {
	s := m.styles
	contentWidth := m.width - 10
	if contentWidth < 40 {
		contentWidth = 40
	}

	titleText := sec.Title
	if sec.Icon != "" {
		titleText = sec.Icon + " " + titleText
	}
	title := s.Heading.Render("  " + titleText)

	var kvLines []string
	for _, kv := range sec.KeyValueItems {
		keyText := kv.Key
		if kv.Icon != "" {
			keyText = kv.Icon + " " + keyText
		}
		var valStr string
		if strings.HasPrefix(kv.Value, "http") || strings.Contains(kv.Value, ".com") || strings.Contains(kv.Value, "github") {
			valStr = s.Link.Render(kv.Value)
		} else {
			valStr = s.Body.Render(kv.Value)
		}
		kvLines = append(kvLines,
			fmt.Sprintf("  %s  %s",
				s.Accent.Width(14).Render(keyText),
				valStr,
			),
		)
	}

	card := s.Card.Width(contentWidth).Render(
		strings.Join(kvLines, "\n"),
	)

	return lipgloss.JoinVertical(lipgloss.Left, title, card)
}

func (m Model) renderFallbackSection(sec data.Section) string {
	s := m.styles
	contentWidth := m.width - 10
	if contentWidth < 40 {
		contentWidth = 40
	}

	titleText := sec.Title
	if sec.Icon != "" {
		titleText = sec.Icon + " " + titleText
	}
	title := s.Heading.Render("  " + titleText)
	body := s.Body.Width(contentWidth - 6).Render(sec.Content)
	if body == "" {
		body = s.Muted.Render("No content available for this section.")
	}

	card := s.Card.Width(contentWidth).Render(body)
	return lipgloss.JoinVertical(lipgloss.Left, title, card)
}

// Utilities

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
