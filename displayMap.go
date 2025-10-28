package main

import (
	"fmt"
	"slices"
)

func convertMapToString(model *model) string {
	// ich wil die map eigentlich gerne mit 3x3 Tiles beschreiben#
	intermidiate := [][]rune{}
	if len(model.tiles) == 0 {
		return ""
	}

	// wilken adressiert die mit [x][y] => x gibt das einzelne an, y die zeile

	for y := range model.tiles {
		lineRow1 := []rune{}
		lineRow2 := []rune{}
		lineRow3 := []rune{}
		for x := range model.tiles[y] {

			stringTile := [3][]rune{[]rune{' ', ' ', ' '}, []rune{' ', ' ', ' '}, []rune{' ', ' ', ' '}}
			// links nach rechts
			// Soll wohl mit x und y laufen, andersrum als normalerweise,
			stringTile = displayTracks(model.tiles[x][y], stringTile)
			lineRow1 = append(lineRow1, stringTile[0]...)
			lineRow2 = append(lineRow2, stringTile[1]...)
			lineRow3 = append(lineRow3, stringTile[2]...)
		}
		intermidiate = append(intermidiate, lineRow1)
		intermidiate = append(intermidiate, lineRow2)
		intermidiate = append(intermidiate, lineRow3)
	}
	// Display Trains
	var indexes []int
	for i := range model.Trains {
		indexes = append(indexes, i)
	}
	slices.Sort(indexes)
	var trains []*Train
	for i := range indexes {
		trains = append(trains, model.Trains[i])
	}

	for i, train := range trains {
		for _, wagon := range train.Waggons {
			// 0 => lins 1 oben 2 rechts 3 unten
			actualX := -1
			actualY := -1
			switch wagon.Position[2] {
			case 1:
				actualX = (wagon.Position[0] * 3)
				actualY = (wagon.Position[1] * 3) + 1
			case 2:
				actualX = (wagon.Position[0] * 3) + 1
				actualY = (wagon.Position[1] * 3)
			case 3:
				actualX = (wagon.Position[0] * 3) + 2
				actualY = (wagon.Position[1] * 3) + 1
			case 4:
				actualX = (wagon.Position[0] * 3) + 1
				actualY = (wagon.Position[1] * 3) + 2
			}
			if !(actualX == -1 || actualY == -1) {
				intermidiate[actualY+1][actualX] = []rune(fmt.Sprint(i))[0]
			}

		}
	}

	result := ""
	for _, line := range intermidiate {
		result += string(line) + "\n"
	}
	return result
}

func displayTracks(tile Tile, lines [3][]rune) [3][]rune {
	tracks := tile.Tracks
	// Gleise anzeigen
	if tile.IsBlocked {
		lines[0][0] = '*'
		lines[0][2] = '*'
		lines[2][0] = '*'
		lines[2][2] = '*'
	}

	if tile.IsPlattform {
		lines[0][0] = '#'
		lines[0][2] = '#'
		lines[2][0] = '#'
		lines[2][2] = '#'
	}

	if tracks[0] {
		lines[1][0] = '-'
	}
	if tracks[1] {
		lines[0][1] = '|'
	}
	if tracks[2] {
		lines[1][2] = '-'
	}
	if tracks[3] {
		lines[2][1] = '|'
	}
	if tracks[0] && tracks[2] {
		lines[1][1] = '-'
	}
	if tracks[1] && tracks[3] {
		lines[1][1] = '|'
	}
	if tracks[0] && tracks[2] && tracks[1] && tracks[3] {
		lines[1][1] = '+'
	}
	//  Signale Anzeigen
	if tile.Signals[0] {
		lines[0][0] = 'S'
	}
	if tile.Signals[1] {
		lines[0][2] = 'S'
	}
	if tile.Signals[2] {
		lines[2][2] = 'S'
	}
	if tile.Signals[3] {
		lines[2][0] = 'S'
	}

	return lines
}

func calculateSubtile(x_in_grid int, y_in_grid int) (int, bool, bool) {
	isSignal, isTrack := false, false // mitte geht ja auhc noch und das ist keins von beidem da man die nicht direkt bauen kann !
	var subtile int
	switch {
	// Gleise
	case x_in_grid == 0 && y_in_grid == 1:
		subtile = 1
		isTrack = true
	case x_in_grid == 1 && y_in_grid == 0:
		subtile = 2
		isTrack = true
	case x_in_grid == 2 && y_in_grid == 1:
		subtile = 3
		isTrack = true
	case x_in_grid == 1 && y_in_grid == 2:
		subtile = 4
		isTrack = true
		// Signale
	case x_in_grid == 0 && y_in_grid == 1:
		subtile = 1
		isSignal = true
	case x_in_grid == 1 && y_in_grid == 0:
		subtile = 2
		isSignal = true
	case x_in_grid == 2 && y_in_grid == 1:
		subtile = 3
		isSignal = true
	case x_in_grid == 1 && y_in_grid == 2:
		subtile = 4
		isSignal = true
	default:
		subtile = -1
	}
	return subtile, isTrack, isSignal
}
