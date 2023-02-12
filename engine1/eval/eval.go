package eval

import (
	"chess/game"
	pc "chess/game/piece"
	rs "chess/game/result"
)

func Evaluate(g *game.GameState) int {
	if g.IsOver {
		switch g.Result {
		case rs.Draw:
			return 0
		case rs.WhiteWins:
			return 10000
		case rs.BlackWins:
			return -10000
		}
	}
	var total int = 0
	for _, slot := range g.WhitePieces {
		if slot.IsInvalid() {
			continue
		}
		total += getPieceWeight(slot.Piece) + centerTable[slot.Pos.Row*8+slot.Pos.Column]
	}
	for _, slot := range g.BlackPieces {
		if slot.IsInvalid() {
			continue
		}
		total -= getPieceWeight(slot.Piece) + centerTable[slot.Pos.Row*8+slot.Pos.Column]
	}
	return total
}

var centerTable = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 3, 3, 3, 3, 0, 0,
	0, 0, 5, 15, 15, 5, 0, 0,
	0, 0, 5, 15, 15, 5, 0, 0,
	0, 0, 3, 3, 3, 3, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

func getPieceWeight(piece pc.Piece) int {
	switch piece {
	case pc.WhiteQueen, pc.BlackQueen:
		return 1500
	case pc.WhiteKing, pc.BlackKing:
		return 10000
	case pc.WhiteBishop, pc.BlackBishop:
		return 300
	case pc.WhiteRook, pc.BlackRook:
		return 700
	case pc.WhiteHorsie, pc.BlackHorsie:
		return 300
	case pc.WhitePawn, pc.BlackPawn:
		return 100
	}
	return 0
}
