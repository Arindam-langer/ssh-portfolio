#  SSH TUI Portfolio

An interactive, SSH-accessible terminal portfolio built with **Go**, **Wish**, and **Bubble Tea**. Showcase your experience, skills, projects, and contact details directly in anyone's terminal — zero client-side dependencies required beyond a standard `ssh` client.

```bash
ssh -p 2222 localhost
```

---

## 󰒋 Features

- 󰏘 **Arch Dark Theme** — Deep background (`#0C0D11`), elevated surfaces (`#171A25`), and crisp teal accents (`#7EBAB5`).
- 󰒓 **100% Config-Driven** — Customize your entire portfolio in `config.yaml`. Add, remove, or reorder tabs without touching Go code.
-  **5 Dynamic Section Types**:
  - `text` — Free-form bio, about me, and markdown-friendly summaries.
  - `skill_list` — Categorized technical skills with optional visual progress bars (0–100%).
  - `timeline` — Work history, education, and milestones with dates, tags, and bullet points.
  - `projects` — Project showcase with expandable cards (`Enter` to toggle details) and repo links.
  - `key_value` — Clean two-column layout for contact information, socials, or specs.
- 󰈔 **Direct Resume Download (SCP)** — Visitors can download your resume directly over SCP without leaving their terminal.
-  **Vim & Number Navigation** — Move smoothly with `Tab`, `h/j/k/l`, arrow keys, or number keys `1–9`.
- 󰋖 **Help Overlay** — Press `?` at any point to view keybindings.
-  **Production & Container Ready** — Includes a multi-stage `Dockerfile` and `docker-compose.yml`.

---

##  Quick Start

### 1. Clone & Setup Configuration

```bash
git clone https://github.com/Arindam-Langer/ssh-portfolio.git
cd ssh-portfolio

# Create your personal configuration from the annotated template
cp config.example.yaml config.yaml
```

### 2. Run Locally (Go)

```bash
# Build & start the server
go build -o ssh-portfolio .
./ssh-portfolio
```

In a new terminal window:
```bash
ssh -p 2222 localhost
```

### 3. Run with Docker

```bash
# Using Docker Compose
docker compose up -d

# Or using Docker directly
docker build -t ssh-portfolio .
docker run -d -p 2222:2222 --name ssh-portfolio ssh-portfolio
```

---

##  Configuration Guide

All personal data, branding, and tabs are configured in `config.yaml`. See [config.example.yaml](config.example.yaml) for a fully annotated template.

### Compulsory vs. Optional Fields

| Field | Required? | Description |
|---|---|---|
| `profile.name` | **Required** | Your name/handle shown in SSH banners, titles, and footers. |
| `profile.tagline` | *Optional* | Subtitle displayed beneath your name. |
| `profile.ascii_logo` | *Optional* | Custom ASCII art banner (defaults to standard logo if omitted). |
| `sections` | **Required** | List of tabs to render (must include at least 1). |
| `sections[].title` | **Required** | The label shown on the navigation tab. |
| `sections[].type` | **Required** | One of `text`, `skill_list`, `timeline`, `projects`, or `key_value`. |
| `sections[].icon` | *Optional* | Nerd Font icon glyph preceding the tab title. |

### Supported Section Types

| Type | Intended Use | Key Fields |
|---|---|---|
| `text` | About Me / Philosophy | `content` (multiline text) |
| `skill_list` | Languages, Frameworks, Cloud | `categories` → `name`, `items` (`name`, `icon`, `level`) |
| `timeline` | Work Experience, Education | `timeline_items` → `title`, `subtitle`, `period`, `location`, `bullets`, `tech` |
| `projects` | GitHub / Highlighted Work | `projects` → `name`, `tagline`, `tech`, `url`, `bullets` |
| `key_value` | Contact, Socials, Hardware Specs | `key_value_items` → `key`, `value`, `icon` |

---

##  Keyboard Shortcuts

| Key | Action |
|---|---|
| `Tab` / `Shift+Tab` | Next / Previous section |
| `← →` / `h l` | Navigate sections |
| `↑ ↓` / `j k` | Scroll section content |
| `1` – `9` | Jump directly to section tab |
| `Enter` | Expand / collapse project details |
| `?` | Toggle help overlay |
| `q` / `Ctrl+C` | Disconnect / Exit |

---

## 󰈔 SCP Resume Download

Place your PDF resume inside the `resume/` directory (e.g. `resume/sample_resume.pdf`). Visitors can fetch it directly using (need to research on how to make it work with hyperlink to click and download using scp or something else):

```bash
scp -P 2222 localhost:resume/sample_resume.pdf ./
```

---

## 󰒍 Deploying to a VPS / Cloud

To run this as a public SSH service on port `22` or `2222`:

1. Forward port `2222` on your cloud firewall / security group.
2. Run via Docker Compose:
   ```bash
   docker compose up -d
   ```
3. Share your command with recruiters and colleagues:
   ```bash
   ssh yourdomain.com -p 2222
   ```
#### current To-do
- adding support for hyprlinks in values or links part of the config.