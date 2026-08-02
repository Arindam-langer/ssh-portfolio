# TUI Skills — SSH Portfolio

## Tech Stack
- **Go 1.26+** with `charm.land/bubbletea/v2`, `charm.land/lipgloss/v2`, `charm.land/wish/v2`
- **Theme**: Arch Palette
  - Background: `#0C0D11`
  - Surface: `#171A25`
  - Accent/Main: `#7EBAB5`
  - Text: `#F6F5F5`
  - Sub (Muted): `#454864`
- **Containerization**: Multi-stage `Dockerfile` & `docker-compose.yml`

## Key API Notes (Bubbletea v2)
- `View()` returns `tea.View`, not `string`
- Use `tea.NewView(content)` and set `v.AltScreen = true`
- `lipgloss.Color()` is a **function** (returns `color.Color`), not a type
- `tea.Tick()` for animations
- `tea.KeyMsg.String()` for key matching
- `tea.WindowSizeMsg` for responsive layouts

## Architecture
- `internal/data/content.go` — All resume content as Go structs
- `internal/tui/model.go` — Root model with tabbed navigation, splash, scroll
- `internal/tui/theme.go` — Arch theme palette definitions
- `internal/tui/styles.go` — Lip Gloss style system derived from theme
- `main.go` — Wish SSH server on port 2222
- `Dockerfile` & `docker-compose.yml` — Containerization configs