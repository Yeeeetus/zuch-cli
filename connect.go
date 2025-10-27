package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

func startListeningToBackend(m *model) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws?username=asd", nil)
	if err != nil {
		return
	}
	m.conn = conn
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
			case "train.remove":
				unMarshalAndSend[trainRemoveMSG](envelope)
			case "tiles.block":
				unMarshalAndSend[blockedTilesMSG](envelope)
			case "tiles.unblock":
				unMarshalAndSend[unblockedTilesMSG](envelope)
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

func sendTileUpdateMSG(m *model, msgType string) {
	if m.conn != nil {
		pos := tileUpdateMSG{Position: [3]int{m.tile_x, m.tile_y, m.subtile}}
		envelope := wsEnvelopeSend{Type: msgType, Msg: pos}
		err := m.conn.WriteJSON(envelope)
		if err != nil {
			fmt.Println(err)
		}
	}
}
