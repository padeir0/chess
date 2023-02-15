// naive move generation utilizing piece ordering
package basic

import (
	"chess/game"
	pc "chess/game/piece"
	. "chess/movegen/common"
)

var _ Generator = &MoveGenerator{}

func ConsumeAll(mg *MoveGenerator) []*game.Move {
	output := []*game.Move{}
	mv := mg.Next()
	for mv != nil {
		output = append(output, mv)
		mg.g.UnMove()
		mv = mg.Next()
	}
	return output
}

func NewMoveGenerator(g *game.GameState) *MoveGenerator {
	slots := &g.WhitePieces
	if g.BlackTurn {
		slots = &g.BlackPieces
	}
	mg := &MoveGenerator{
		g:     g,
		slots: slots,
	}
	return mg
}

type MoveGenerator struct {
	g           *game.GameState
	pseudoLegal *MovesFor
	currPseudo  int

	slots    *[]game.Slot
	currSlot int
}

func (this *MoveGenerator) Next() *game.Move {
	for this.currSlot < len(*this.slots) {
		slot := (*this.slots)[this.currSlot]
		piece := slot.Piece // Move() may alter the slot
		if slot.Piece == pc.Empty {
			this.pseudoLegal = nil
			this.currSlot++
			continue
		}
		if this.pseudoLegal == nil {
			moves := PseudoLegalMoves(this.g, slot.Pos, slot.Piece)
			this.pseudoLegal = &MovesFor{
				From: slot.Pos,
				To:   moves,
			}
			this.currPseudo = 0
		}
		for this.currPseudo < len(this.pseudoLegal.To) {
			to := this.pseudoLegal.To[this.currPseudo]
			this.currPseudo++
			lastCapt := this.g.MovesSinceLastCapture
			ok, capture := this.g.Move(this.pseudoLegal.From, to)
			if ok {
				move := &game.Move{
					Piece:   piece,
					From:    this.pseudoLegal.From,
					To:      to,
					Capture: capture,

					MovesSinceLastCapture: lastCapt,
				}
				return move
			}
		}
		this.pseudoLegal = nil
		this.currSlot++
	}
	return nil
}

func PseudoLegalMoves(g *game.GameState, Pos game.Point, piece pc.Piece) []game.Point {
	switch piece {
	case pc.BlackKing, pc.WhiteKing:
		return genKingMoves(Pos)
	case pc.BlackKnight, pc.WhiteKnight:
		return genHorsieMoves(Pos)
	case pc.BlackQueen, pc.WhiteQueen:
		return genQueenMoves(g, Pos)
	case pc.BlackBishop, pc.WhiteBishop:
		return genBishopMoves(g, Pos)
	case pc.BlackRook, pc.WhiteRook:
		return genRookMoves(g, Pos)
	case pc.BlackPawn:
		return genBlackPawnMoves(Pos)
	case pc.WhitePawn:
		return genWhitePawnMoves(Pos)
	}
	return nil
}

func genKingMoves(pos game.Point) []game.Point {
	output := []game.Point{}
	for _, offset := range game.KingOffsets {
		newpos := game.Point{
			Column: pos.Column + offset.Column,
			Row:    pos.Row + offset.Row,
		}
		output = append(output, newpos)
	}
	return output
}

func genHorsieMoves(pos game.Point) []game.Point {
	output := []game.Point{}
	for _, offset := range game.HorsieOffsets {
		newpos := game.Point{
			Column: pos.Column + offset.Column,
			Row:    pos.Row + offset.Row,
		}
		output = append(output, newpos)
	}
	return output
}

func genRookMoves(g *game.GameState, from game.Point) []game.Point {
	output := []game.Point{}
	fromPiece := g.Board.AtPos(from)
	for _, offset := range game.RookOffsets {
		for i := 1; i < 7; i++ {
			pos := game.Point{
				Column: from.Column + (offset.Column * i),
				Row:    from.Row + (offset.Row * i),
			}
			if pos.IsInvalid() {
				break
			}
			piece := g.Board.AtPos(pos)
			if piece == pc.Empty {
				output = append(output, pos)
			} else if piece.IsBlack() != fromPiece.IsBlack() {
				output = append(output, pos)
				break
			} else {
				break
			}
		}
	}
	return output
}

func genBishopMoves(g *game.GameState, from game.Point) []game.Point {
	output := []game.Point{}

	fromPiece := g.Board.AtPos(from)
	for _, offset := range game.BishopOffsets {
		for i := 1; i < 7; i++ {
			to := game.Point{
				Column: from.Column + (offset.Column * i),
				Row:    from.Row + (offset.Row * i),
			}
			if to.IsInvalid() {
				break
			}
			piece := g.Board.AtPos(to)
			if piece == pc.Empty {
				output = append(output, to)
			} else if piece.IsBlack() != fromPiece.IsBlack() {
				output = append(output, to)
				break
			} else {
				break
			}
		}
	}

	return output
}

func genQueenMoves(g *game.GameState, pos game.Point) []game.Point {
	return append(genBishopMoves(g, pos), genRookMoves(g, pos)...)
}
func genBlackPawnMoves(pos game.Point) []game.Point {
	return []game.Point{
		{Row: pos.Row + 1, Column: pos.Column},
		{Row: pos.Row + 1, Column: pos.Column - 1},
		{Row: pos.Row + 1, Column: pos.Column + 1},
	}
}

func genWhitePawnMoves(pos game.Point) []game.Point {
	return []game.Point{
		{Row: pos.Row - 1, Column: pos.Column},
		{Row: pos.Row - 1, Column: pos.Column - 1},
		{Row: pos.Row - 1, Column: pos.Column + 1},
	}
}
