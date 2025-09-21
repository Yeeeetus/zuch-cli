package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	railwayMap [][]Tile
	trains     []Train
	actions    []string // items on the to-do list#
	helpKeys   keyMap
	help       help.Model
	connected  bool
}

type Tile struct {
	Tracks      [4]bool
	Signals     [4]bool
	IsPlattform bool
	IsBlocked   bool
}

type Train struct {
	//irgendwie noch zusmamenfassung in einen Zug, selbstreferenz funktioniert nicht
	Position [3]int //x,y,track(1,2,3,4) ->
	Goal     [3]int //nur fürs testen
	MaxSpeed int

	//
	Size  int
	Cargo int
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func initialModel() model {
	return model{help: help.New(), helpKeys: keys, connected: false}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch {

		// These keys should exit the program.
		case key.Matches(msg, m.helpKeys.connect):
			m.connected = !m.connected
		case key.Matches(msg, m.helpKeys.help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.helpKeys.quit):
			return m, tea.Quit
		}

		// Return the updated model to the Bubble Tea runtime for processing.
		// Note that we're not returning a command.
	}
	return m, nil
}

func (m model) View() string {
	result := ""
	switch m.connected {
	case false:
		result += "Bitte mit einer instanz verbinden \n\n"
	}

	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))

	result += borderStyle.Render(convertMapToString(m))
	result += "\n\n\n\n\n\n"
	result += m.help.View(m.helpKeys)

	// Send the UI for rendering
	return result
}
