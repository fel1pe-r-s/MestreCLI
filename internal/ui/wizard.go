package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	choices     []string
	cursor      int
	selected    int
	Step        int
	ProjectType string
	ProjectName string
	Runtime     string // "node" or "bun"
	UseDocker   bool
	UseTurbo    bool
	textInput   textinput.Model
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
		Runtime:   "node", // Default
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnterAltScreen)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		// NAVIGATION
		case "up", "k":
			if m.Step == 0 {
				if m.cursor > 0 {
					m.cursor--
				}
			} else if m.Step == 2 { // Runtime (2 options)
				if m.cursor > 0 {
					m.cursor--
				}
			}
		case "down", "j":
			if m.Step == 0 {
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			} else if m.Step == 2 {
				if m.cursor < 1 {
					m.cursor++
				}
			}

		// CONFIRMATION
		case "enter":
			if m.Step == 0 { // Project Type
				m.selected = m.cursor
				m.ProjectType = m.choices[m.cursor]
				m.Step = 1 // Go to Name
				return m, nil
			} else if m.Step == 1 { // Project Name
				m.ProjectName = m.textInput.Value()
				if m.ProjectName == "" {
					m.ProjectName = "meu-projeto"
				}
				m.Step = 2 // Go to Runtime
				m.cursor = 0
				return m, nil
			} else if m.Step == 2 { // Runtime
				if m.cursor == 0 {
					m.Runtime = "node"
				} else {
					m.Runtime = "bun"
				}
				m.Step = 3 // Go to Docker
				return m, nil
			}

		// YES/NO Selection for Steps 3 & 4
		case "y", "Y":
			if m.Step == 3 { // Docker
				m.UseDocker = true
				m.Step = 4 // Go to Turbo
				return m, nil
			} else if m.Step == 4 { // Turbo
				m.UseTurbo = true
				return m, tea.Quit // FINISH
			}
		case "n", "N":
			if m.Step == 3 {
				m.UseDocker = false
				m.Step = 4
				return m, nil
			} else if m.Step == 4 {
				m.UseTurbo = false
				return m, tea.Quit // FINISH
			}
		}
	}

	if m.Step == 1 {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	s := "🧙‍♂️ Mestre Stack Wizard\n\n"

	if m.Step == 0 {
		s += "1. Escolha o tipo de projeto:\n\n"
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			checked := " "
			if i == m.selected {
				checked = "x"
			}
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}
		s += "\n(Use setas e Enter)"
	} else if m.Step == 1 {
		s += "2. Qual o nome do projeto?\n\n"
		s += m.textInput.View()
		s += "\n\n(Digite e aperte Enter)"
	} else if m.Step == 2 {
		s += "3. Qual Runtime usar?\n\n"
		c1, c2 := " ", " "
		if m.cursor == 0 {
			c1 = ">"
		} else {
			c2 = ">"
		}
		s += fmt.Sprintf("%s [ ] Node.js (Padrão)\n", c1)
		s += fmt.Sprintf("%s [ ] Bun (Rápido)\n", c2)
		s += "\n(Use setas e Enter)"
	} else if m.Step == 3 {
		s += "4. Configurar Docker Compose? (y/n)\n"
	} else if m.Step == 4 {
		s += "5. Configurar TurboRepo? (y/n)\n"
	}

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
