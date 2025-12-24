package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	choices     []string
	cursor      int
	selected    int
	step        int
	ProjectType string
}

func InitialModel() Model {
	return Model{
		choices:  []string{"Backend API (Clean Arch)", "Universal App (Web/Mobile/Desktop)", "Monorepo Genérico"},
		selected: -1,
		step:     0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.selected = m.cursor
			m.ProjectType = m.choices[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := "Qual tipo de projeto você quer criar?\n\n"

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

	s += "\nUse as setas para mover, Enter para selecionar.\n"
	return s
}

func StartWizard() (Model, error) {
	p := tea.NewProgram(InitialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Erro no Wizard: %v", err)
		os.Exit(1)
	}
	// Aqui chamariamos o scaffold baseado na escolha
	// Por enquanto, apenas printamos.
	// fmt.Println("Você escolheu:", m.projectType)
	// (Precisariamos retornar o model para acessar o dado, mas para o teste inicial está ok)
	return m.(Model), nil
}
