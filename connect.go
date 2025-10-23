package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

func startListeningToBackend() {
	conn, _, _ := websocket.DefaultDialer.Dial("ws://localhost:8080/ws?username=asd", nil)
	go func() {

		for {
			var envelope wsEnvelope
			err := conn.ReadJSON(&envelope)
			if err != nil {
				fmt.Println("Error beim Ws", err.Error())
				return
			}

			switch envelope.Type {
			case "game.initialLoad":
				unMarshalAndSend[gamestateTemp](envelope)
			case "rail.create":
				unMarshalAndSend[railCreateMSG](envelope)
			case "rail.remove":
				unMarshalAndSend[railRemoveMSG](envelope)
			case "signal.create":
				unMarshalAndSend[signalCreateMSG](envelope)
			case "signal.remove":
				unMarshalAndSend[signalRemoveMSG](envelope)
			case "train.move":
				unMarshalAndSend[trainMoveMSG](envelope)
			case "train.create":
				unMarshalAndSend[Train](envelope)
			}

		}
	}()
}

func unMarshalAndSend[T any](envelope wsEnvelope) T {
	var recievedMSG T
	err := json.Unmarshal(envelope.Msg, &recievedMSG)
	if err != nil {
		fmt.Println("EROOR", err)
	}
	p.Send(recievedMSG)
	return recievedMSG
}

type wsEnvelope struct {
	Type string
	Msg  json.RawMessage
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
	Trains    []Train
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
