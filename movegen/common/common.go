package common

import "chess/game"

type Generator interface {
	Next() (game.Move, bool)
}

type MovesFor struct {
	From game.Point
	To   []game.Point
}
