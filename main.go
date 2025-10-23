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
	tiles     [][]Tile
	Trains    map[int]*Train
	actions   []string // items on the to-do list#
	helpKeys  keyMap
	help      help.Model
	connected bool
}

type Tile struct {
	Tracks      [4]bool
	Signals     [4]bool
	IsPlattform bool
	IsBlocked   bool
	ActiveTile  *ActiveTile
}

type ActiveTile struct {
	Category   *ActiveTileCategory
	Name       string
	Level      int
	Stations   []*Station //Stationen, die in der Nähe sind. wird mit changeStationTile verwaltet
	Storage    map[string]int
	maxStorage int //maximum Lager pro Gut -> sonst kann es zu unwiederruflichen auffüllen kommen
}
type ActiveTileCategory struct {
	Productioncycles []Produktionszyklus `json:"Produktionszyklen"`
}

type Produktionszyklus struct {
	Consumtion                 map[string]int `json:"Verbrauch"`
	Produktion                 map[string]int `json:"Produktion"`
	VerfuegbareLevelUndScaling []int          `json:"verfügbareLevelUndScaling"`
}

type Train struct {
	Waggons            []TrainType //Alle müssen nebeneinander spawnen
	Schedule           Schedule
	NextStop           Stop     //nur fürs testen
	currentPath        [][3]int //neu berechnen bei laden
	currentPathSignals [][3]int
	Name               string
	waiting            bool
	Id                 int
}

type TrainType struct {
	Position [3]int //x,y,sub
	MaxSpeed int
	Id       int
	Size     int
	Cargo    int
}

var p = tea.NewProgram(initialModel())

func main() {
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
	return model{help: help.New(), helpKeys: keys, connected: false, Trains: make(map[int]*Train)}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch {

		// These keys should exit the program.
		case key.Matches(msg, m.helpKeys.connect):
			if len(m.tiles) == 0 {
				startListeningToBackend()
			}
		case key.Matches(msg, m.helpKeys.help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.helpKeys.quit):
			return m, tea.Quit
		}

		// Return the updated model to the Bubble Tea runtime for processing.
		// Note that we're not returning a command.
	case gamestateTemp:
		for _, train := range msg.Trains {
			m.Trains[train.Id] = &train
		}
		m.tiles = msg.Tiles
		m.connected = true
	case trainMoveMSG:
		m.Trains[msg.Id].Waggons = msg.Waggons
	case Train: // train.create
		m.Trains[msg.Id] = &msg

	case signalCreateMSG:
		m.tiles[msg.Position[0]][msg.Position[1]].Signals[(msg.Position[2])-1] = true
	case signalRemoveMSG:
		m.tiles[msg.Position[0]][msg.Position[1]].Signals[(msg.Position[2])-1] = false
	case railCreateMSG:
		m.tiles[msg.Position[0]][msg.Position[1]].Tracks[(msg.Position[2])-1] = true
	case railRemoveMSG:
		m.tiles[msg.Position[0]][msg.Position[1]].Tracks[(msg.Position[2])-1] = false
	}

	return m, nil
}

func (m model) View() string {
	result := ""
	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))
	switch m.connected {
	case false:
		result += "Bitte mit einer instanz verbinden \n\n"
	case true:
		result += borderStyle.Render(convertMapToString(&m))
		result += "\n"

	}

	result += m.help.View(m.helpKeys)
	if !m.help.ShowAll {
		result += "\n\n\n"
	}

	trainList := "Unsere Aktuellen Züge :D \n"
	if len(m.Trains) > 0 {
		for i, train := range m.Trains {
			for _, waggon := range train.Waggons {
				trainList += fmt.Sprint(i) + fmt.Sprintf(": %d;%d;%d", waggon.Position[0], waggon.Position[1], waggon.Position[2]) + "\n"

			}
		}
	}
	// Send the UI for rendering
	return lipgloss.JoinHorizontal(0.2, result, borderStyle.Render(trainList))
}
