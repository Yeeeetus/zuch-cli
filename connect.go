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
				var state gamestateTemp
				err := json.Unmarshal(envelope.Msg, &state)
				if err != nil {
					fmt.Println("EROOR", err)
				}
				p.Send(state)

			case "tile.update":
				var update tileUpdateMSG
				err := json.Unmarshal(envelope.Msg, &update)
				if err != nil {
					fmt.Println("EROOR", err)
				}
				p.Send(update)
			case "train.move":
				var recievedMSG trainMoveMSG
				err := json.Unmarshal(envelope.Msg, &recievedMSG)
				if err != nil {
					fmt.Println("EROOR", err)
				}
				p.Send(recievedMSG)
			case "train.create":
				var recievedMSG Train
				err := json.Unmarshal(envelope.Msg, &recievedMSG)
				if err != nil {
					fmt.Println("EROOR", err)
				}
				p.Send(recievedMSG)
			}

		}
	}()
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
	X       int
	Y       int
	Subtile int // 1 => links, 2 => oben, 3 => rechts, 4 => unten
	// Wilken hat sich entschlossen immer wenn ein subtile als int gespeichert wird bei 1 anzufangen und wenn es ein bool[4] ist bei 0, also kann sein das es sich irgendwo verschiebt aber das kriegen wir sicher noch behoben Bei schienen auch analog
	Subject string //
	Action  string // remove, build
}
