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

		case "up", "k":
			if m.Step == 0 {
				if m.cursor > 0 {
					m.cursor--
				}
			}
		case "down", "j":
			if m.Step == 0 {
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			}
		case "enter":
			if m.Step == 0 {
				m.selected = m.cursor
				m.ProjectType = m.choices[m.cursor]
				m.Step = 1 // Avança para o input
				return m, nil
			} else if m.Step == 1 {
				m.ProjectName = m.textInput.Value()
				if m.ProjectName == "" {
					m.ProjectName = "meu-projeto" // Fallback
				}
				return m, tea.Quit
			}
		}
	}

	// Atualiza input apenas no passo 1
	if m.Step == 1 {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	s := ""

	if m.Step == 0 {
		s += "🧙‍♂️ Mestre Stack Wizard\n\n"
		s += "Escolha o tipo de projeto:\n\n"

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
		s += "\n(Use setas para navegar e Enter para confirmar)\n"
	} else if m.Step == 1 {
		s += "Nome do Projeto:\n\n"
		s += m.textInput.View()
		s += "\n\n(Digite o nome e aperte Enter para criar)\n"
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
