// segregated move generation between captures and quiet moves
package segregated

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

func ConsumeAllQuiet(mg *MoveGenerator) []*game.Move {
	output := []*game.Move{}
	mv := mg.NextQuiet()
	for mv != nil {
		output = append(output, mv)
		mg.g.UnMove()
		mv = mg.NextQuiet()
	}
	return output
}

func ConsumeAllCaptures(mg *MoveGenerator) []*game.Move {
	output := []*game.Move{}
	mv := mg.NextCapture()
	for mv != nil {
		output = append(output, mv)
		mg.g.UnMove()
		mv = mg.NextCapture()
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
	g *game.GameState

	quietPseudo     *MovesFor
	currQuietPseudo int
	currQuietSlot   int

	pseudoCapture     *MovesFor
	currPseudoCapture int
	currCaptureSlot   int

	slots *[]game.Slot // hopefully it's ordered
}

func (this *MoveGenerator) NextCapture() *game.Move {
	for this.currCaptureSlot < len(*this.slots) {
		slot := (*this.slots)[this.currCaptureSlot]
		piece := slot.Piece // Move() may alter the slot
		if slot.Piece == pc.Empty {
			this.pseudoCapture = nil
			this.currCaptureSlot++
			continue
		}
		if this.pseudoCapture == nil {
			moves := genMoves(this.g, slot.Pos, slot.Piece, false)
			this.pseudoCapture = &MovesFor{
				From: slot.Pos,
				To:   moves,
			}
			this.currPseudoCapture = 0
		}
		for this.currPseudoCapture < len(this.pseudoCapture.To) {
			to := this.pseudoCapture.To[this.currPseudoCapture]
			this.currPseudoCapture++
			lastCapt := this.g.MovesSinceLastCapture
			ok, capture := this.g.Move(this.pseudoCapture.From, to)
			if ok && capture != nil {
				move := &game.Move{
					Piece:   piece,
					From:    this.pseudoCapture.From,
					To:      to,
					Capture: capture,

					MovesSinceLastCapture: lastCapt,
				}
				return move
			} else if ok && capture == nil {
				panic("should have captured something")
			}
		}
		this.pseudoCapture = nil
		this.currCaptureSlot++
	}
	return nil
}

func (this *MoveGenerator) NextQuiet() *game.Move {
	for this.currQuietSlot < len(*this.slots) {
		slot := (*this.slots)[this.currQuietSlot]
		piece := slot.Piece // Move() may alter the slot
		if slot.Piece == pc.Empty {
			this.quietPseudo = nil
			this.currQuietSlot++
			continue
		}
		if this.quietPseudo == nil {
			moves := genMoves(this.g, slot.Pos, slot.Piece, true)
			this.quietPseudo = &MovesFor{
				From: slot.Pos,
				To:   moves,
			}
			this.currQuietPseudo = 0
		}
		for this.currQuietPseudo < len(this.quietPseudo.To) {
			to := this.quietPseudo.To[this.currQuietPseudo]
			this.currQuietPseudo++
			lastCapt := this.g.MovesSinceLastCapture
			ok, capture := this.g.Move(this.quietPseudo.From, to)
			if ok && capture == nil {
				move := &game.Move{
					Piece:   piece,
					From:    this.quietPseudo.From,
					To:      to,
					Capture: capture,

					MovesSinceLastCapture: lastCapt,
				}
				return move
			} else if capture != nil {
				panic("should not have captured something")
			}
		}
		this.quietPseudo = nil
		this.currQuietSlot++
	}
	return nil
}

func (this *MoveGenerator) Next() *game.Move {
	capt := this.NextCapture()
	if capt != nil {
		return capt
	}
	return this.NextQuiet()
}

func genMoves(g *game.GameState, Pos game.Point, piece pc.Piece, quiet bool) []game.Point {
	switch piece {
	case pc.BlackKing, pc.WhiteKing:
		return genKingMoves(g, Pos, quiet)
	case pc.BlackKnight, pc.WhiteKnight:
		return genHorsieMoves(g, Pos, quiet)
	case pc.BlackQueen, pc.WhiteQueen:
		return genQueenMoves(g, Pos, quiet)
	case pc.BlackBishop, pc.WhiteBishop:
		return genBishopMoves(g, Pos, quiet)
	case pc.BlackRook, pc.WhiteRook:
		return genRookMoves(g, Pos, quiet)
	case pc.BlackPawn:
		return genBlackPawnMoves(Pos, quiet)
	case pc.WhitePawn:
		return genWhitePawnMoves(Pos, quiet)
	}
	return nil
}

func genKingMoves(g *game.GameState, pos game.Point, quiet bool) []game.Point {
	output := []game.Point{}
	for _, offset := range game.KingOffsets {
		newpos := game.Point{
			Column: pos.Column + offset.Column,
			Row:    pos.Row + offset.Row,
		}
		if newpos.IsInvalid() {
			continue
		}
		if g.Board.AtPos(newpos).IsOccupied() != quiet {
			output = append(output, newpos)
		}
	}
	return output
}

func genHorsieMoves(g *game.GameState, pos game.Point, quiet bool) []game.Point {
	output := []game.Point{}
	for _, offset := range game.HorsieOffsets {
		newpos := game.Point{
			Column: pos.Column + offset.Column,
			Row:    pos.Row + offset.Row,
		}
		if newpos.IsInvalid() {
			continue
		}
		if g.Board.AtPos(newpos).IsOccupied() != quiet {
			output = append(output, newpos)
		}
	}
	return output
}

func genRookMoves(g *game.GameState, from game.Point, quiet bool) []game.Point {
	output := []game.Point{}
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
			if piece.IsOccupied() != quiet {
				output = append(output, pos)
			}
			if piece.IsOccupied() {
				break // something is blocking the way
			}
		}
	}
	return output
}

func genBishopMoves(g *game.GameState, from game.Point, quiet bool) []game.Point {
	output := []game.Point{}
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
			if piece.IsOccupied() != quiet {
				output = append(output, to)
			}
			if piece.IsOccupied() {
				break // something is blocking the way
			}
		}
	}

	return output
}

func genQueenMoves(g *game.GameState, pos game.Point, quiet bool) []game.Point {
	return append(genBishopMoves(g, pos, quiet), genRookMoves(g, pos, quiet)...)
}

func genBlackPawnMoves(pos game.Point, quiet bool) []game.Point {
	if quiet {
		return []game.Point{
			{Row: pos.Row + 1, Column: pos.Column},
		}
	}
	return []game.Point{
		{Row: pos.Row + 1, Column: pos.Column - 1},
		{Row: pos.Row + 1, Column: pos.Column + 1},
	}
}

func genWhitePawnMoves(pos game.Point, quiet bool) []game.Point {
	if quiet {
		return []game.Point{
			{Row: pos.Row - 1, Column: pos.Column},
		}
	}
	return []game.Point{
		{Row: pos.Row - 1, Column: pos.Column - 1},
		{Row: pos.Row - 1, Column: pos.Column + 1},
	}
}
