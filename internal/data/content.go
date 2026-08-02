package data

// Skill represents a technical skill with proficiency level
type Skill struct {
	Name  string
	Level int // 0-100
}

// Experience represents a work experience entry
type Experience struct {
	Title    string
	Company  string
	Location string
	Period   string
	Bullets  []string
	Tech     string
}

// Project represents a portfolio project
type Project struct {
	Name      string
	Tagline   string
	Tech      string
	GitHubURL string
	Bullets   []string
}

// Education represents an education entry
type Education struct {
	School     string
	Degree     string
	Period     string
	Location   string
	GPA        string
	Coursework []string
}

// Contact holds all contact information
type Contact struct {
	Name     string
	Email    string
	Phone    string
	LinkedIn string
	GitHub   string
	Location string
}

// ============================================================================
// Portfolio Data — sourced from Arindam Langer's resume
// ============================================================================

var PersonalInfo = Contact{
	Name:     "Arindam Langer",
	Email:    "arindamlanger@gmail.com",
	Phone:    "+91-7006395984",
	LinkedIn: "linkedin.com/in/arindam-langer",
	GitHub:   "github.com/Arindam-Langer",
	Location: "Jammu, India",
}

var AboutText = `Backend engineer who builds things that survive production.

I design high-performance microservices in Go and Node.js, architect 
AI/RAG pipelines that run fully offline, and ship containerized systems 
on Kubernetes with zero-downtime deploys.

Currently interning at OnClick Innovations where I've built 4 production 
AI video generation pipelines — cutting ad production from hours to minutes.
Previously architected a query engine microservice at Resolyte consumed by 
3+ downstream services.

I believe great software is invisible — it just works, scales, and stays 
out of the user's way.`

var SkillCategories = map[string][]Skill{
	"Languages": {
		{Name: "Go", Level: 92},
		{Name: "JavaScript", Level: 88},
		{Name: "Python", Level: 82},
		{Name: "C", Level: 70},
		{Name: "SQL", Level: 85},
	},
	"Frameworks": {
		{Name: "Node.js", Level: 90},
		{Name: "Express.js", Level: 88},
		{Name: "React", Level: 75},
		{Name: "Flask", Level: 78},
		{Name: "LangChainGo", Level: 72},
		{Name: "Bubble Tea", Level: 85},
		{Name: "gRPC", Level: 82},
	},
	"Databases": {
		{Name: "PostgreSQL", Level: 90},
		{Name: "pgvector", Level: 78},
		{Name: "MongoDB", Level: 75},
		{Name: "Redis", Level: 85},
		{Name: "BullMQ", Level: 80},
	},
	"AI & ML": {
		{Name: "Ollama", Level: 82},
		{Name: "OpenAI API", Level: 85},
		{Name: "Deepgram", Level: 72},
		{Name: "Google Gemini", Level: 78},
		{Name: "RAG Pipelines", Level: 88},
	},
	"DevOps": {
		{Name: "Docker", Level: 90},
		{Name: "Kubernetes", Level: 82},
		{Name: "AWS S3", Level: 75},
		{Name: "FFmpeg", Level: 70},
		{Name: "Git", Level: 92},
		{Name: "Linux", Level: 90},
	},
}

// Ordered category names for consistent rendering
var SkillCategoryOrder = []string{
	"Languages",
	"Frameworks",
	"Databases",
	"AI & ML",
	"DevOps",
}

var Experiences = []Experience{
	{
		Title:    "Software Development Intern",
		Company:  "OnClick Innovations",
		Location: "Mohali, India",
		Period:   "Aug 2025 – Present",
		Bullets: []string{
			"Built 4 production AI video generation pipelines (UGC, Veo, Template) orchestrating Google Gemini, HeyGen, ElevenLabs, and Google Veo — cutting ad production time from hours to minutes.",
			"Architected an async job queue with BullMQ and Redis to handle concurrent AI workloads, with retry logic, per-job progress tracking, and credit metering across 5+ providers.",
			"Engineered a template management system with RBAC publishing, tag/category search, B-roll upload, and a questionnaire-driven creation flow.",
			"Delivered a Live Ambient Scribe POC in 3 days: real-time audio capture via Deepgram streamed to OpenAI for live SOAP note and differential diagnosis generation.",
		},
		Tech: "Node.js, Express.js, React, Python, Flask, Redis, BullMQ, FFmpeg, Docker, OpenAI, Deepgram, AWS S3",
	},
	{
		Title:    "Backend Developer Intern",
		Company:  "Resolyte",
		Location: "Jammu, India",
		Period:   "Mar 2025 – Aug 2025",
		Bullets: []string{
			"Designed a high-performance query engine microservice from scratch in Go and PostgreSQL, consumed by 3+ downstream services.",
			"Developed dynamic SQL generation with Squirrel and exposed type-safe gRPC APIs for inter-service communication across a distributed system.",
			"Containerized with Docker and deployed on Kubernetes with CI/CD integration, achieving zero-downtime rollouts.",
		},
		Tech: "Go, PostgreSQL, gRPC, Docker, Kubernetes, Squirrel",
	},
}

var Projects = []Project{
	{
		Name:      "Authn-Service",
		Tagline:   "Authentication Microservice",
		Tech:      "Go · PostgreSQL · Redis · Docker",
		GitHubURL: "https://github.com/Arindam-Langer/authn-service",
		Bullets: []string{
			"6-endpoint auth microservice with JWT access tokens, 7-day rotating refresh tokens, bcrypt hashing, and UUID v5 phone anonymization — zero raw PII stored.",
			"Token rotation with reuse detection — a single replayed token instantly revokes every active session via Redis blocklist.",
			"11,000 req/sec on core middleware, 1,500 req/sec on JWT verification hot path, 99% responses under 15ms.",
		},
	},
	{
		Name:      "OllamaChat",
		Tagline:   "Privacy-First RAG Engine with TUI",
		Tech:      "Go · pgvector · Ollama · Docker · Bubble Tea",
		GitHubURL: "https://github.com/Arindam-Langer/OllamaChat",
		Bullets: []string{
			"Fully offline RAG pipeline: PDF ingestion, overlap-aware chunking, 768-dim vector embeddings via Ollama, cosine-similarity retrieval from Dockerized pgvector.",
			"Terminal UI in Bubble Tea (Elm architecture) with streaming chat, multi-line input, and a file browser — zero external API calls, complete data privacy.",
		},
	},
	{
		Name:      "AI Farmer Assistant",
		Tagline:   "Multimodal AI Chatbot for Agriculture",
		Tech:      "Python · Flask · LangChain · Ollama",
		GitHubURL: "https://github.com/Arindam-Langer/AI-Assistant-For-Farmers",
		Bullets: []string{
			"Multimodal AI chatbot for plant disease diagnosis from image and text inputs, fine-tuned Llama3 with RAG using LangChain and Chroma.",
			"Integrated feedback loops and real-time data retrieval to iteratively improve response accuracy and usability.",
		},
	},
}

var EducationInfo = Education{
	School:   "Central University of Jammu",
	Degree:   "BTech in Computer Science and Engineering",
	Period:   "Nov 2022 – Aug 2026",
	Location: "Jammu, India",
	GPA:      "7.6/10",
	Coursework: []string{
		"Data Structures and Algorithms (DSA)",
		"Design and Analysis of Algorithms (DAA)",
		"Operating Systems (OS)",
		"Computer Networking",
		"Network Security",
		"Database Management Systems (DBMS)",
		"Computer Organization and Architecture",
		"Theory of Computation",
		"Software Engineering",
		"Artificial Intelligence",
		"Machine Learning",
	},
}

// ASCII art logo
var Logo = `
   █████╗ ██████╗ ██╗███╗   ██╗██████╗  █████╗ ███╗   ███╗
  ██╔══██╗██╔══██╗██║████╗  ██║██╔══██╗██╔══██╗████╗ ████║
  ███████║██████╔╝██║██╔██╗ ██║██║  ██║███████║██╔████╔██║
  ██╔══██║██╔══██╗██║██║╚██╗██║██║  ██║██╔══██║██║╚██╔╝██║
  ██║  ██║██║  ██║██║██║ ╚████║██████╔╝██║  ██║██║ ╚═╝ ██║
  ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝
`

var Tagline = "Backend Engineer · Go & Node.js · AI Pipeline Architect"

var SplashFrames = []string{
	"Establishing secure connection...",
	"Loading portfolio modules...",
	"Initializing Bubble Tea runtime...",
	"Rendering terminal UI...",
	"Welcome aboard.",
}
