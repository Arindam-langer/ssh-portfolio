# SSH TUI Portfolio — Arindam Langer

An SSH-accessible terminal portfolio built with **Go**, **Bubble Tea**, **Lip Gloss**, and **Wish**, featuring the custom **Arch** theme palette.

```bash
ssh -p 2222 localhost
```

## Features

- 🎨 **Arch Theme** — Deep dark background (`#0C0D11`), elevated surface (`#171A25`), teal main accent (`#7EBAB5`), crisp white text (`#F6F5F5`), and muted accents (`#454864`)
- ✨ **Animated Splash** — Typewriter-style intro with ASCII art
- 📑 **6 Sections** — About, Skills, Experience, Projects, Education, Contact
- 📊 **Skill Bars** — Visual progress bars for technical skills
- 📂 **Expandable Projects** — Press Enter to expand/collapse project details
- 📄 **Resume Download** — SCP your resume directly from the server
- ⌨️ **Vim-style Navigation** — `h/j/k/l`, `tab`, number keys `1-6`
- ❓ **Help Overlay** — Press `?` for keyboard shortcuts
- 🐳 **Docker Ready** — Multi-stage `Dockerfile` & `docker-compose.yml` included

## Quick Start (Local)

```bash
# Build & run binary directly
go build -o ssh-portfolio .
./ssh-portfolio

# In another terminal:
ssh -p 2222 localhost
```

## Quick Start (Docker)

```bash
# Using Docker Compose
docker compose up -d

# Or building directly
docker build -t ssh-portfolio .
docker run -d -p 2222:2222 --name ssh-portfolio ssh-portfolio
```

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `tab` / `shift+tab` | Navigate sections |
| `← →` / `h l` | Navigate sections |
| `↑ ↓` / `j k` | Scroll content |
| `1`-`6` | Jump to section |
| `enter` | Expand/collapse project |
| `?` | Help overlay |
| `q` / `ctrl+c` | Quit |

## Resume Download via SCP

```bash
scp -P 2222 localhost:resume/arindam_resume.pdf ./
```


#### Current To-Do's
# Dynamic Config Refactoring Tasks

This checklist tracks the logical steps we'll take to implement the dynamic configuration architecture. We will pause after each major step (commit) so you can review the code and commit it yourself.

## [x] Commit 1: Configuration Engine
- `[x]` Define new dynamic structs in `internal/data/config.go` to represent Sections and flexible content types.
- `[x]` Write the `LoadConfig(path string)` function to read `config.yaml`.
- `[x]` Update `config.yaml` to match the new dynamic `sections` array format.
- `[x]` Update `main.go` to invoke `LoadConfig` on startup.

## [ ] Commit 2: Dynamic Tabs & Core UI Refactor
- `[ ]` Remove the hardcoded tab names and indices from `internal/tui/model.go`.
- `[ ]` Update `NewModel` to generate tabs dynamically from the loaded configuration.
- `[ ]` Refactor state tracking (like active tab index, scrolling) to be completely independent of fixed sections.

## [ ] Commit 3: Dynamic View Rendering
- `[ ]` Refactor the `View()` function in `model.go` to use a switch statement based on the active section's `Type` (e.g., text, timeline, skill_list, projects, key_value).
- `[ ]` Migrate the old hardcoded rendering logic (from `data.go` constants) to the new dynamic renderers, pulling data from the active section.

## [ ] Commit 4: Cleanup & Final Polish
- `[ ]` Remove the old `internal/data/content.go` as it is now obsolete.
- `[ ]` Ensure fallback values and error handling are robust if the user creates a malformed YAML.
- `[ ]` Test rendering across all dynamic section types.
