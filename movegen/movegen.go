package movegen

import (
	"chess/game"
	pc "chess/game/piece"
)

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

type MovesFor struct {
	From game.Position
	To   []game.Position
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
			moves := game.PseudoLegalMoves(this.g, slot.Pos, slot.Piece)
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
