package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	choices  []string
	cursor   int
	selected int
	Step     int

	// Configs
	ProjectType string
	ProjectName string
	Runtime     string // node, bun
	Framework   string // fastify, nestjs
	ORM         string // prisma, drizzle, none
	Database    string // postgres, sqlite, mongo
	UseDocker   bool
	UseTurbo    bool

	textInput textinput.Model
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "meu-projeto-incrivel"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30

	return Model{
		choices:   []string{"Backend API (Clean Arch)", "Universal App (Web/Mobile/Desktop)", "Monorepo Genérico"},
		selected:  -1,
		Step:      0,
		textInput: ti,
		Runtime:   "node",
		Framework: "fastify",
		ORM:       "prisma",
		Database:  "postgres",
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnterAltScreen)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	isBackend := strings.Contains(m.ProjectType, "Backend")

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.Step == 0 { // Type
				if m.cursor > 0 {
					m.cursor--
				}
			} else if m.Step == 2 { // Runtime
				if m.cursor > 0 {
					m.cursor--
				}
			} else if m.Step == 3 && isBackend { // Framework
				if m.cursor > 0 {
					m.cursor--
				}
			} else if m.Step == 4 && isBackend { // ORM
				if m.cursor > 0 {
					m.cursor--
				}
			} else if m.Step == 5 && isBackend { // Database
				if m.cursor > 0 {
					m.cursor--
				}
			}

		case "down", "j":
			if m.Step == 0 { // Type
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			} else if m.Step == 2 { // Runtime (2)
				if m.cursor < 1 {
					m.cursor++
				}
			} else if m.Step == 3 && isBackend { // Framework (2)
				if m.cursor < 1 {
					m.cursor++
				}
			} else if m.Step == 4 && isBackend { // ORM (3)
				if m.cursor < 2 {
					m.cursor++
				}
			} else if m.Step == 5 && isBackend { // Database (3)
				if m.cursor < 2 {
					m.cursor++
				}
			}

		case "enter":
			if m.Step == 0 { // Select Type
				m.selected = m.cursor
				m.ProjectType = m.choices[m.cursor]
				m.Step = 1 // -> Name
				return m, nil

			} else if m.Step == 1 { // Input Name
				m.ProjectName = m.textInput.Value()
				if m.ProjectName == "" {
					m.ProjectName = "meu-projeto"
				}
				m.Step = 2 // -> Runtime
				m.cursor = 0
				return m, nil

			} else if m.Step == 2 { // Select Runtime
				if m.cursor == 0 {
					m.Runtime = "node"
				} else {
					m.Runtime = "bun"
				}

				// BRANCHING LOGIC
				if strings.Contains(m.ProjectType, "Backend") {
					m.Step = 3 // -> Framework
				} else {
					m.Step = 6 // -> Docker (Skip Framework/ORM/DB)
				}
				m.cursor = 0
				return m, nil

			} else if m.Step == 3 { // Select Framework
				if m.cursor == 0 {
					m.Framework = "fastify"
				} else {
					m.Framework = "nestjs"
				}
				m.Step = 4 // -> ORM
				m.cursor = 0
				return m, nil

			} else if m.Step == 4 { // Select ORM
				if m.cursor == 0 {
					m.ORM = "prisma"
				} else if m.cursor == 1 {
					m.ORM = "drizzle"
				} else {
					m.ORM = "none"
				}
				m.Step = 5 // -> DB
				m.cursor = 0
				return m, nil

			} else if m.Step == 5 { // Select DB
				if m.cursor == 0 {
					m.Database = "postgres"
				} else if m.cursor == 1 {
					m.Database = "sqlite"
				} else {
					m.Database = "mongo"
				}
				m.Step = 6 // -> Docker
				return m, nil
			}

		case "y", "Y":
			// Step 6: Docker
			if m.Step == 6 {
				m.UseDocker = true
				if strings.Contains(m.ProjectType, "Universal") || strings.Contains(m.ProjectType, "Monorepo") {
					m.Step = 7 // -> Turbo
				} else {
					return m, tea.Quit // Done for Backend
				}
				return m, nil
			}
			// Step 7: Turbo
			if m.Step == 7 {
				m.UseTurbo = true
				return m, tea.Quit
			}

		case "n", "N":
			if m.Step == 6 {
				m.UseDocker = false
				if strings.Contains(m.ProjectType, "Universal") || strings.Contains(m.ProjectType, "Monorepo") {
					m.Step = 7
				} else {
					return m, tea.Quit
				}
				return m, nil
			}
			if m.Step == 7 {
				m.UseTurbo = false
				return m, tea.Quit
			}
		}
	}

	if m.Step == 1 {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	s := "🧙‍♂️ Mestre Stack Wizard V2\n\n"

	// Helper to render choices
	renderChoices := func(options []string) string {
		out := ""
		for i, opt := range options {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			out += fmt.Sprintf("%s [ ] %s\n", cursor, opt)
			// checked logic omitted for simplicity since we move on enter
		}
		return out
	}

	switch m.Step {
	case 0:
		s += "1. Tipo de Projeto:\n\n"
		for i, choice := range m.choices {
			cursor, checked := " ", " "
			if m.cursor == i {
				cursor = ">"
			}
			if i == m.selected {
				checked = "x"
			}
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}
	case 1:
		s += "2. Nome do Projeto:\n\n" + m.textInput.View()
	case 2:
		s += "3. Runtime:\n\n" + renderChoices([]string{"Node.js (LTS)", "Bun (Ultra Fast)"})
	case 3:
		s += "4. Framework Backend:\n\n" + renderChoices([]string{"Fastify (Clean Arch)", "NestJS (Modular)"})
	case 4:
		s += "5. ORM (Banco de Dados):\n\n" + renderChoices([]string{"Prisma (Recomendado)", "Drizzle (Leve)", "Nenhum"})
	case 5:
		s += "6. Banco de Dados:\n\n" + renderChoices([]string{"PostgreSQL", "SQLite", "MongoDB"})
	case 6:
		s += "7. Configurar Docker/Compose? (y/n)\n"
	case 7:
		s += "8. Configurar TurboRepo? (y/n)\n"
	}

	s += "\n(Enter para confirmar, Ctrl+C para sair)\n"
	return s
}

func StartWizard() (Model, error) {
	p := tea.NewProgram(InitialModel())
	m, err := p.Run()
	if err != nil {
		return Model{}, err
	}
	return m.(Model), nil
}
