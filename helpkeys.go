package main

import "github.com/charmbracelet/bubbles/key"

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	connect      key.Binding
	disconnect   key.Binding
	save         key.Binding
	pause        key.Binding
	unpause      key.Binding
	quit         key.Binding
	help         key.Binding
	createTrack  key.Binding
	removeTrack  key.Binding
	createSignal key.Binding
	removeSignal key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.connect, k.help, k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.connect, k.disconnect, k.pause, k.unpause}, // first column
		{k.save, k.quit}, // second column
		{k.createTrack, k.removeTrack, k.createSignal, k.removeSignal},
	}
}

var keys = keyMap{
	connect: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "connect to Server"),
	),
	disconnect: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "disconnect"),
	),
	save: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "save game"),
	),
	pause: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Pause game"),
	),
	unpause: key.NewBinding(
		key.WithKeys("u", "r"),
		key.WithHelp("u/r", "Resume game"),
	),
	quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	help: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "Hilfe"),
	),
	createTrack: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "Gleis bauen"),
	),
	removeTrack: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "gleis entfernen"),
	),
	createSignal: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "Signal bauen"),
	),
	removeSignal: key.NewBinding(
		key.WithKeys("4"),
		key.WithHelp("4", "Signal entfernen"),
	),
}
