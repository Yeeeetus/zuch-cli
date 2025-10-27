package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

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

type wsEnvelope struct {
	Type string
	Msg  json.RawMessage
}

type wsEnvelopeSend struct {
	Type string
	Msg  any
}

type trainMoveMSG struct {
	Id      int
	Waggons []TrainType
}

type gamestateTemp struct {
	Users     []User
	Schedules []Schedule
	Stations  []Station
	Tiles     [][]Tile
	Trains    map[int]Train
}

// blockiert mehrere, nicht nur ein tile
type blockedTilesMSG struct {
	Tiles [][2]int
}

// blockiert mehrere, nicht nur ein tile
type unblockedTilesMSG struct {
	Tiles [][2]int
}

type User struct {
	username    string
	isConnected bool
	connection  *websocket.Conn
}

type UserInput struct {
	action    string
	username  string
	parameter any
}

type Schedule struct {
	name  string
	user  User
	stops []Stop
}

type Stop struct {
	id   int
	goal [3]int //wird Station, bei der man sich ein Tile abholt
}

type Station struct {
	plattforms []Plattform
}

type Plattform struct {
	tiles []Tile
}

type tileUpdateMSG struct {
	Position [3]int // 0 => links, 1 => oben, 2 => rechts, 3 => unten
	// Wilken hat sich entschlossen immer wenn ein subtile als int gespeichert wird bei 1 anzufangen und wenn es ein bool[4] ist bei 0, also kann sein das es sich irgendwo verschiebt aber das kriegen wir sicher noch behoben Bei schienen auch analog
}

// nicht im backend vorhanden, nur hier da bubbletea das will
type signalCreateMSG struct {
	Position [3]int // 0 => links, 1 => oben, 2 => rechts, 3 => unten
}
type signalRemoveMSG struct {
	Position [3]int // 0 => links, 1 => oben, 2 => rechts, 3 => unten
}
type railCreateMSG struct {
	Position [3]int // 0 => links, 1 => oben, 2 => rechts, 3 => unten
}
type railRemoveMSG struct {
	Position [3]int // 0 => links, 1 => oben, 2 => rechts, 3 => unten
}

type trainRemoveMSG struct {
	id int
}
