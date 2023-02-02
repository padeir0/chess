package engine

import (
	"chess/game"
	pc "chess/game/piece"
	pr "chess/game/promotion"
)

type Move struct {
	From  game.Position
	To    game.Position
	State *game.GameState
}

func GenerateMoves(g *game.GameState) []*Move {
	slots := g.WhitePieces
	if g.BlackTurn {
		slots = g.BlackPieces
	}
	outStates := []*Move{}
	newGS := g.Copy()
	for _, slot := range slots {
		moves := pseudoLegalMoves(slot.Position, slot.Piece)
		for _, to := range moves {
			ok, _ := newGS.Move(slot.Position, to, pr.Queen)
			if ok {
				move := &Move{
					From:  slot.Position,
					To:    to,
					State: newGS,
				}
				outStates = append(outStates, move)
				newGS = g.Copy()
			}
		}
	}
	return outStates
}

func pseudoLegalMoves(Pos game.Position, piece pc.Piece) []game.Position {
	switch piece {
	case pc.BlackCastleKing, pc.WhiteCastleKing, pc.BlackKing, pc.WhiteKing:
		return genKingMoves(Pos)
	case pc.BlackHorsie, pc.WhiteHorsie:
		return genHorsieMoves(Pos)
	case pc.BlackQueen, pc.WhiteQueen:
		return genQueenMoves(Pos)
	case pc.BlackBishop, pc.WhiteBishop:
		return genBishopMoves(Pos)
	case pc.BlackRook, pc.WhiteRook,
		pc.BlackMovedRook, pc.WhiteMovedRook:
		return genRookMoves(Pos)
	case pc.BlackPassantPawn, pc.BlackMovedPawn, pc.BlackPawn:
		return genBlackPawnMoves(Pos)
	case pc.WhitePassantPawn, pc.WhiteMovedPawn, pc.WhitePawn:
		return genWhitePawnMoves(Pos)
	}
	return nil
}

var kingOffsets = []game.Position{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1} /*    */, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

func genKingMoves(pos game.Position) []game.Position {
	output := []game.Position{}
	for _, offset := range kingOffsets {
		newpos := game.Position{
			Column: pos.Column + offset.Column,
			Row:    pos.Row + offset.Row,
		}
		output = append(output, newpos)
	}
	return output
}

var horsieOffsets = []game.Position{
	{-2, -1}, {-2, +1},
	{-1, -2}, {-1, +2},
	{+1, -2}, {+1, +2},
	{+2, -1}, {+2, +1},
}

func genHorsieMoves(pos game.Position) []game.Position {
	output := []game.Position{}
	for _, offset := range horsieOffsets {
		newpos := game.Position{
			Column: pos.Column + offset.Column,
			Row:    pos.Row + offset.Row,
		}
		output = append(output, newpos)
	}
	return output
}

func genRookMoves(pos game.Position) []game.Position {
	output := []game.Position{}
	for i := 0; i <= 7; i++ {
		if i != pos.Row {
			newPos := game.Position{
				Row:    i,
				Column: pos.Column,
			}
			output = append(output, newPos)
		}
		if i != pos.Column {
			newPos := game.Position{
				Row:    pos.Row,
				Column: i,
			}
			output = append(output, newPos)
		}
	}
	return output
}

func genBishopMoves(pos game.Position) []game.Position {
	firstDiagPos := game.Position{Row: 0, Column: 0}
	diff := pos.Row - pos.Column
	if diff < 0 {
		firstDiagPos = game.Position{Row: -diff, Column: 0}
	} else {
		firstDiagPos = game.Position{Row: 0, Column: diff}
	}

	secDiagPos := game.Position{Row: 0, Column: 7}
	if diff < 0 {
		secDiagPos = game.Position{Row: 7 + diff, Column: 7}
	} else {
		secDiagPos = game.Position{Row: 7, Column: diff}
	}

	output := []game.Position{}
	for i := 0; i < 7; i++ {
		firstDiag := game.Position{
			Row:    firstDiagPos.Row - i,
			Column: firstDiagPos.Column + i,
		}
		if firstDiag.IsValid() && firstDiag != pos {
			output = append(output, firstDiag)
		}

		secDiag := game.Position{
			Row:    secDiagPos.Row - i,
			Column: secDiagPos.Column - i,
		}
		if secDiag.IsValid() && secDiag != pos {
			output = append(output, secDiag)
		}
	}

	return output
}

func genQueenMoves(pos game.Position) []game.Position {
	return append(genBishopMoves(pos), genRookMoves(pos)...)
}

func genBlackPawnMoves(pos game.Position) []game.Position {
	return []game.Position{
		{Row: pos.Row + 2, Column: pos.Column},
		{Row: pos.Row + 1, Column: pos.Column},
		{Row: pos.Row + 1, Column: pos.Column - 1},
		{Row: pos.Row + 1, Column: pos.Column + 1},
	}
}

func genWhitePawnMoves(pos game.Position) []game.Position {
	return []game.Position{
		{Row: pos.Row - 2, Column: pos.Column},
		{Row: pos.Row - 1, Column: pos.Column},
		{Row: pos.Row - 1, Column: pos.Column - 1},
		{Row: pos.Row - 1, Column: pos.Column + 1},
	}
}
