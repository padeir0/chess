package evaluation

import (
	"chess/game"
	pc "chess/game/piece"
)

// maximize for white
func Evaluate(g *game.GameState) float64 {
	var total float64 = 0
	for _, slot := range g.WhitePieces {
		if slot.IsInvalid() {
			continue
		}
		pinfo := &pieceInfo{
			Piece:   slot.Piece,
			Pos:     slot.Pos,
			IsBlack: false,
		}
		total += getPieceWeight(g, pinfo)
	}
	for _, slot := range g.BlackPieces {
		if slot.IsInvalid() {
			continue
		}
		pinfo := &pieceInfo{
			Piece:   slot.Piece,
			Pos:     slot.Pos,
			IsBlack: true,
		}
		total -= getPieceWeight(g, pinfo)
	}
	return total
}

type pieceInfo struct {
	Piece   pc.Piece
	Pos     game.Position
	IsBlack bool
}

const kingWeight float64 = 100

func getPieceWeight(g *game.GameState, pinfo *pieceInfo) float64 {
	if pinfo.Piece.IsKingLike() {
		return kingWeight + protectionWeight(g, pinfo)
	}
	var pieceWeight float64 = 0
	if pinfo.Piece.IsQueenLike() {
		pieceWeight = 15
	} else if pinfo.Piece.IsRookLike() {
		pieceWeight = 7
	} else if pinfo.Piece.IsBishopLike() {
		pieceWeight = bishopWeight(g, pinfo)
	} else if pinfo.Piece.IsHorsieLike() {
		pieceWeight = 3 + horsieMobility(g, pinfo)
	} else if pinfo.Piece.IsPawnLike() {
		pieceWeight = pawnWeight(g, pinfo)
	} else {
		panic("invalid piece: " + pinfo.Piece.String())
	}
	defmod := defMod(g, pinfo.Pos, pinfo.IsBlack)
	return pieceWeight + defmod
}

func isEndgame(g *game.GameState) bool {
	return g.TotalValuablePieces <= 8
}

const defmod float64 = 0.025

func defMod(g *game.GameState, pos game.Position, isBlack bool) float64 {
	// see if position is attacked by their own pieces
	defenders := &g.BlackPieces
	if !isBlack {
		defenders = &g.WhitePieces
	}
	var output float64 = 0
	for _, slot := range *defenders {
		if slot.Piece == pc.Empty {
			continue
		}
		switch slot.Piece {
		case pc.WhiteRook, pc.BlackRook,
			pc.BlackMovedRook, pc.WhiteMovedRook:
			if pos.Column == slot.Pos.Column ||
				pos.Row == slot.Pos.Row {
				output += defmod
			}
		case pc.WhiteBishop, pc.BlackBishop:
			if Abs(pos.Column-slot.Pos.Column) ==
				Abs(pos.Row-slot.Pos.Row) {
				output += defmod
			}
		case pc.BlackQueen, pc.WhiteQueen:
			if Abs(pos.Column-slot.Pos.Column) ==
				Abs(pos.Row-slot.Pos.Row) ||
				pos.Column == slot.Pos.Column ||
				pos.Row == slot.Pos.Row {
				output += defmod
			}
		case pc.BlackPawn, pc.BlackPassantPawn, pc.BlackMovedPawn:
			if pos.Row-slot.Pos.Row == 1 &&
				Abs(pos.Column-slot.Pos.Column) == 1 {
				output += defmod
			}
		case pc.WhitePawn, pc.WhitePassantPawn, pc.WhiteMovedPawn:
			if pos.Row-slot.Pos.Row == -1 &&
				Abs(pos.Column-slot.Pos.Column) == 1 {
				output += defmod
			}
		case pc.WhiteHorsie, pc.BlackHorsie:
			if !((Abs(slot.Pos.Column-pos.Column) == 2 &&
				Abs(slot.Pos.Row-pos.Row) == 1) ||
				(Abs(slot.Pos.Column-pos.Column) == 1 &&
					Abs(slot.Pos.Row-pos.Row) == 2)) {
				output += defmod
			}
		case pc.BlackKing, pc.WhiteKing, pc.BlackCastleKing, pc.WhiteCastleKing:
			ColDiff := slot.Pos.Column - pos.Column
			RowDiff := slot.Pos.Row - pos.Row
			if !(((ColDiff == 1) || (ColDiff == 0) || (ColDiff == -1)) &&
				((RowDiff == 1) || (RowDiff == 0) || (RowDiff == -1))) {
				output += defmod
			}
		}
	}
	return output
}

func protectionWeight(g *game.GameState, pinfo *pieceInfo) float64 {
	var weight float64 = 1
	for _, offset := range game.KingOffsets {
		pos := game.Position{
			Column: pinfo.Pos.Column + offset.Column,
			Row:    pinfo.Pos.Row + offset.Row,
		}
		if pos.IsInvalid() {
			continue
		}
		piece := g.Board.AtPos(pos)
		if piece.IsQueenLike() {
			weight += 0.2
		} else if piece.IsRookLike() {
			weight += 0.15
		} else {
			weight += 0.1
		}
	}
	return weight
}

func horsieMobility(g *game.GameState, pinfo *pieceInfo) float64 {
	var mobMod float64 = 0
	for _, offset := range game.HorsieOffsets {
		pos := game.Position{
			Column: pinfo.Pos.Column + offset.Column,
			Row:    pinfo.Pos.Row + offset.Row,
		}
		if pos.IsInvalid() {
			continue
		}
		piece := g.Board.AtPos(pos)
		if piece == pc.Empty {
			mobMod += 0.05
		} else if piece.IsBlack() != pinfo.IsBlack {
			mobMod += 0.1
		}
	}
	return mobMod
}

func bishopWeight(g *game.GameState, pinfo *pieceInfo) float64 {
	bishopMob := bishopMobility(g, pinfo)
	if hasBishopPair(g, pinfo) {
		return 3.5 + bishopMob
	}
	return 3 + bishopMob
}

func bishopMobility(g *game.GameState, pinfo *pieceInfo) float64 {
	var mobMod float64 = 0
	for _, offset := range game.BishopOffsets {
		pos := game.Position{
			Column: pinfo.Pos.Column + offset.Column,
			Row:    pinfo.Pos.Row + offset.Row,
		}
		if pos.IsInvalid() {
			continue
		}
		if g.Board.AtPos(pos) == pc.Empty {
			mobMod += 0.01
		}
	}
	return mobMod
}

// inneficient but will do for now
func hasBishopPair(g *game.GameState, pinfo *pieceInfo) bool {
	lightBishop := false
	darkBishop := false
	collection := &g.WhitePieces
	if pinfo.IsBlack {
		collection = &g.BlackPieces
	}
	for _, slot := range *collection {
		if slot.Piece == pc.Empty {
			continue
		}
		if slot.Piece == pc.BlackBishop {
			if (slot.Pos.Column+slot.Pos.Row*8)%2 == 0 {
				lightBishop = true
			} else {
				darkBishop = true
			}
		}
	}
	return lightBishop && darkBishop
}

func pawnWeight(g *game.GameState, pinfo *pieceInfo) float64 {
	var out float64 = pawnColValue(g, pinfo) * pawnRowMultiplier(g, pinfo)
	if isPassedPawn(g, pinfo) {
		out += 0.25
	}
	return out
}

func isPassedPawn(g *game.GameState, pinfo *pieceInfo) bool {
	quant := -1
	if pinfo.IsBlack {
		quant = 1
	}
	for i := pinfo.Pos.Row - quant; i >= 0; i += quant {
		if g.Board.At(i, pinfo.Pos.Column).IsPawnLike() {
			return false
		}
		col := pinfo.Pos.Column - 1
		if col >= 0 && g.Board.At(i, col).IsPawnLike() {
			return false
		}
		col = pinfo.Pos.Column + 1
		if col <= 7 && g.Board.At(i, col).IsPawnLike() {
			return false
		}
	}
	return true
}

func pawnColValue(g *game.GameState, pinfo *pieceInfo) float64 {
	if isEndgame(g) {
		switch pinfo.Pos.Column {
		case 7, 0:
			return 1.2
		case 6, 1:
			return 1.1
		case 5, 2:
			return 1
		case 4, 3:
			return 0.95
		}
	}
	switch pinfo.Pos.Column {
	case 7, 0:
		return 0.8
	case 6, 1:
		return 0.95
	case 5, 2:
		return 1
	case 4, 3:
		return 1.1
	}
	return 1
}

func pawnRowMultiplier(g *game.GameState, pinfo *pieceInfo) float64 {
	if pinfo.IsBlack {
		return 1.25 * (float64(pinfo.Pos.Row) / 4)
	}
	return 1.25 * (float64(7-pinfo.Pos.Row) / 4)
}

func Abs32(a int32) int32 {
	y := a >> 31
	return (a ^ y) - y
}

func Abs(a int) int {
	y := int32(a) >> 31
	return int((int32(a) ^ y) - y)
}
