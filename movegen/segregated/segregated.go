// segregated move generation between captures and quiet moves
package segregated

import (
	"chess/game"
	pc "chess/game/piece"
	. "chess/movegen/common"

	"fmt"
)

var _ Generator = &MoveGenerator{}

func ConsumeAll(mg *MoveGenerator) []game.Move {
	output := []game.Move{}
	mv, ok := mg.Next()
	for ok {
		output = append(output, mv)
		mg.g.UnMove()
		mv, ok = mg.Next()
	}
	return output
}

func ConsumeAllQuiet(mg *MoveGenerator) []game.Move {
	output := []game.Move{}
	mv, ok := mg.NextQuiet()
	for ok {
		output = append(output, mv)
		mg.g.UnMove()
		mv, ok = mg.NextQuiet()
	}
	return output
}

func ConsumeAllCaptures(mg *MoveGenerator) []game.Move {
	output := []game.Move{}
	mv, ok := mg.NextCapture()
	for ok {
		output = append(output, mv)
		mg.g.UnMove()
		mv, ok = mg.NextCapture()
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

	currQuietOffset int
	currQuietPseudo int
	currQuietSlot   int

	currCaptureOffset int
	currPseudoCapture int
	currCaptureSlot   int

	slots *[]game.Slot // hopefully it's ordered
}

func (this *MoveGenerator) NextCapture() (game.Move, bool) {
	for this.currCaptureSlot < len(*this.slots) {
		slot := (*this.slots)[this.currCaptureSlot]
		if slot.Piece == pc.Empty {
			this.currCaptureSlot++
			this.currPseudoCapture = 0
			this.currCaptureOffset = 0
			continue
		}
		to, hasMove := this.nextMove(slot.Pos, slot.Piece, false)
		for hasMove {
			lastCapt := this.g.MovesSinceLastCapture
			ok, capture := this.g.Move(slot.Pos, to)
			if ok && capture != nil {
				move := game.Move{
					Piece:   slot.Piece,
					From:    slot.Pos,
					To:      to,
					Capture: capture,

					MovesSinceLastCapture: lastCapt,
				}
				return move, true
			} else if ok && capture == nil {
				fmt.Println(this.g.Board.String())
				fmt.Println(this.g.Moves)
				fmt.Println(this.currCaptureSlot, this.currCaptureOffset, this.currPseudoCapture)
				fmt.Println(slot.Piece, slot.Pos, to, capture)
				fmt.Println(slot)
				err := this.g.CheckInvalid()
				if err != "" {
					panic(err)
				}
				panic("should have captured something")
			}
			to, hasMove = this.nextMove(slot.Pos, slot.Piece, false)
		}
		this.currPseudoCapture = 0
		this.currCaptureOffset = 0
		this.currCaptureSlot++
	}
	return game.Move{}, false
}

func (this *MoveGenerator) NextQuiet() (game.Move, bool) {
	for this.currQuietSlot < len(*this.slots) {
		slot := (*this.slots)[this.currQuietSlot]
		if slot.Piece == pc.Empty {
			this.currQuietSlot++
			this.currQuietPseudo = 0
			this.currQuietOffset = 0
			continue
		}
		to, hasMove := this.nextMove(slot.Pos, slot.Piece, true)
		for hasMove {
			lastCapt := this.g.MovesSinceLastCapture
			ok, capture := this.g.Move(slot.Pos, to)
			if ok && capture == nil {
				move := game.Move{
					Piece:   slot.Piece,
					From:    slot.Pos,
					To:      to,
					Capture: capture,

					MovesSinceLastCapture: lastCapt,
				}
				return move, true
			} else if capture != nil {
				fmt.Println(this.currCaptureSlot, this.currPseudoCapture)
				fmt.Println(slot.Piece, slot.Pos, to, capture)
				err := this.g.CheckInvalid()
				if err != "" {
					panic(err)
				}
				panic("should not have captured something")
			}
			to, hasMove = this.nextMove(slot.Pos, slot.Piece, true)
		}
		this.currQuietPseudo = 0
		this.currQuietOffset = 0
		this.currQuietSlot++
	}
	return game.Move{}, false
}

func (this *MoveGenerator) Next() (game.Move, bool) {
	capt, ok := this.NextCapture()
	if ok {
		return capt, ok
	}
	return this.NextQuiet()
}

// generates pseudolegal moves
func (this *MoveGenerator) nextMove(Pos game.Point, piece pc.Piece, quiet bool) (game.Point, bool) {
	switch piece {
	case pc.BlackKing, pc.WhiteKing:
		return this.nextSimpleMove(Pos, quiet, game.KingOffsets)
	case pc.BlackKnight, pc.WhiteKnight:
		return this.nextSimpleMove(Pos, quiet, game.HorsieOffsets)
	case pc.BlackQueen, pc.WhiteQueen:
		return this.nextSlideMove(Pos, quiet, game.QueenOffsets)
	case pc.BlackBishop, pc.WhiteBishop:
		return this.nextSlideMove(Pos, quiet, game.BishopOffsets)
	case pc.BlackRook, pc.WhiteRook:
		return this.nextSlideMove(Pos, quiet, game.RookOffsets)
	case pc.BlackPawn:
		return this.nextBlackPawnMove(Pos, quiet)
	case pc.WhitePawn:
		return this.nextWhitePawnMove(Pos, quiet)
	}
	return game.Point{}, false
}

// for pieces that can move to fixed squares (King, Knight)
func (this *MoveGenerator) nextSimpleMove(pos game.Point, quiet bool, offsets []game.Point) (game.Point, bool) {
	currOffset := &this.currCaptureOffset
	if quiet {
		currOffset = &this.currQuietOffset
	}
	for ; *currOffset < len(offsets); *currOffset++ {
		offset := offsets[*currOffset]
		newpos := game.Point{
			Column: pos.Column + offset.Column,
			Row:    pos.Row + offset.Row,
		}
		if newpos == pos {
			panic("null move")
		}
		if newpos.IsInvalid() {
			continue
		}
		piece := this.g.Board.AtPos(newpos)
		if piece.IsOccupied() != quiet {
			*currOffset += 1
			return newpos, true
		}
	}
	*currOffset = 0
	return game.Point{}, false
}

// for pieces that can slide (Queen, Rook, Bishop)
func (this *MoveGenerator) nextSlideMove(from game.Point, quiet bool, offsets []game.Point) (game.Point, bool) {
	currOffset := &this.currCaptureOffset
	if quiet {
		currOffset = &this.currQuietOffset
	}
	currPseudo := &this.currPseudoCapture
	if quiet {
		currPseudo = &this.currQuietPseudo
	}

	for *currOffset < len(offsets) {
		offset := offsets[*currOffset]
		for *currPseudo < 7 {
			newpos := game.Point{
				Column: from.Column + (offset.Column * (*currPseudo + 1)),
				Row:    from.Row + (offset.Row * (*currPseudo + 1)),
			}
			if newpos.IsInvalid() {
				*currPseudo = 0
				break
			}
			piece := this.g.Board.AtPos(newpos)
			isOccupied := piece.IsOccupied()
			if isOccupied != quiet {
				// if it is occupied we want to stop sliding even after the return
				if isOccupied {
					*currPseudo = 0
					*currOffset += 1
				} else {
					*currPseudo += 1
				}
				return newpos, true
			}
			// if something is blocking the way, we stop sliding
			if isOccupied {
				*currPseudo = 0
				break
			}
			*currPseudo += 1
		}
		*currPseudo = 0
		*currOffset += 1
	}
	return game.Point{}, false
}

func (this *MoveGenerator) nextBlackPawnMove(pos game.Point, quiet bool) (game.Point, bool) {
	if quiet {
		if this.currQuietOffset == 0 {
			this.currQuietOffset++
			return game.Point{Row: pos.Row + 1, Column: pos.Column}, true
		}
		return game.Point{}, false
	}
	if this.currCaptureOffset == 0 {
		this.currCaptureOffset++
		return game.Point{Row: pos.Row + 1, Column: pos.Column - 1}, true
	}
	if this.currCaptureOffset == 1 {
		this.currCaptureOffset++
		return game.Point{Row: pos.Row + 1, Column: pos.Column + 1}, true
	}
	return game.Point{}, false
}

func (this *MoveGenerator) nextWhitePawnMove(pos game.Point, quiet bool) (game.Point, bool) {
	if quiet {
		if this.currQuietOffset == 0 {
			this.currQuietOffset++
			return game.Point{Row: pos.Row - 1, Column: pos.Column}, true
		}
		return game.Point{}, false
	}
	if this.currCaptureOffset == 0 {
		this.currCaptureOffset++
		return game.Point{Row: pos.Row - 1, Column: pos.Column - 1}, true
	}
	if this.currCaptureOffset == 1 {
		this.currCaptureOffset++
		return game.Point{Row: pos.Row - 1, Column: pos.Column + 1}, true
	}
	return game.Point{}, false
}
