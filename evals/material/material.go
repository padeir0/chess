package material

import (
	"chess/game"
	pc "chess/game/piece"
	rs "chess/game/result"
	ifaces "chess/interfaces"
)

var _ ifaces.Evaluator = Evaluate

func Evaluate(g *game.GameState) int {
	if g.IsOver {
		switch g.Result {
		case rs.WhiteWins:
			return 10000
		case rs.BlackWins:
			return -10000
		case rs.Draw:
			return 0
		}
	}
	var total int = 0
	for _, slot := range g.WhitePieces {
		if slot.IsInvalid() {
			continue
		}
		total += getPieceWeight(slot.Piece)
	}
	for _, slot := range g.BlackPieces {
		if slot.IsInvalid() {
			continue
		}
		total -= getPieceWeight(slot.Piece)
	}
	return total
}

func getPieceWeight(p pc.Piece) int {
	switch p {
	case pc.InvalidPiece:
		return 0
	case pc.Empty:
		return 0

	case pc.WhiteKing, pc.BlackKing:
		return 10000
	case pc.WhiteQueen, pc.BlackQueen:
		return 900
	case pc.WhiteRook, pc.BlackRook:
		return 500
	case pc.WhiteBishop, pc.BlackBishop:
		return 330
	case pc.WhiteKnight, pc.BlackKnight:
		return 320
	case pc.WhitePawn, pc.BlackPawn:
		return 100
	}
	return 0
}
