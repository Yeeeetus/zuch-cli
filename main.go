package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

type model struct {
	tiles     [][]Tile
	Trains    map[int]*Train
	helpKeys  keyMap
	conn      *websocket.Conn
	help      help.Model
	connected bool
	x         int
	y         int
	cor_x     int
	cor_y     int
	tile_x    int
	tile_y    int
	subtile   int
	isSignal  bool
	isTrack   bool
}

var p = tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())

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
	case tea.MouseMsg:
		// Prüfe den Typ des Mausereignisses
		switch msg.Button {
		case tea.MouseButtonLeft:
			m.x = msg.X
			m.y = msg.Y
			cor_x, cor_y := m.x-1, m.y-1
			m.tile_x, m.tile_y = cor_x/3, cor_y/3
			x_in_grid, y_in_grid := cor_x%3, cor_y%3

			m.subtile, m.isTrack, m.isSignal = calculateSubtile(x_in_grid, y_in_grid)
		}

	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch {

		// These keys should exit the program.
		case key.Matches(msg, m.helpKeys.connect):
			if len(m.tiles) == 0 {
				startListeningToBackend(&m)
			}
		case key.Matches(msg, m.helpKeys.help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.helpKeys.createTrack):
			sendTileUpdateMSG(&m, "rail.create")
		case key.Matches(msg, m.helpKeys.removeTrack):
			sendTileUpdateMSG(&m, "rail.remove")
		case key.Matches(msg, m.helpKeys.createSignal):
			sendTileUpdateMSG(&m, "signal.create")
		case key.Matches(msg, m.helpKeys.removeSignal):
			sendTileUpdateMSG(&m, "signal.remove")
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
	case trainRemoveMSG:
		delete(m.Trains, msg.id)
	case blockedTilesMSG:
		for _, v := range msg.Tiles {
			m.tiles[v[0]][v[1]].IsBlocked = true
		}
	case unblockedTilesMSG:
		for _, v := range msg.Tiles {
			m.tiles[v[0]][v[1]].IsBlocked = false
		}
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
	trainList := "Unsere Aktuellen Züge :D \n"
	if len(m.Trains) > 0 {
		for i, train := range m.Trains {
			for _, waggon := range train.Waggons {
				trainList += fmt.Sprint(i) + fmt.Sprintf(": %d;%d;%d", waggon.Position[0], waggon.Position[1], waggon.Position[2]) + "\n"

			}
		}
	}
	result = lipgloss.JoinHorizontal(lipgloss.Top, result, "                          ", borderStyle.Render(trainList)) // jaja die sind ein wichtiger platzhalter

	result = lipgloss.JoinVertical(lipgloss.Left, result, m.help.View(m.helpKeys))

	if !m.help.ShowAll {
		result += "\n\n\n"
	}

	location := fmt.Sprintf("Folgendes %3d:%3d, entspricht Feld %2d:%2d mit Subtile %2d, Mögliche Belegungen (Signal %t, Track %t)",
		m.cor_x, m.cor_y, m.tile_x, m.tile_y, m.subtile, m.isSignal, m.isTrack)

	result = lipgloss.JoinVertical(lipgloss.Left, result, location)
	// Send the UI for rendering
	return result
}
