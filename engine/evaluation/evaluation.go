package evaluation

import (
	"chess/game"
	pc "chess/game/piece"
	"fmt"
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

func EvaluatePrint(g *game.GameState) float64 {
	var total float64 = 0
	pinfos := []*pieceInfo{}
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
		pinfos = append(pinfos, pinfo)
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
		pinfos = append(pinfos, pinfo)
	}
	metrics(pinfos)
	return total
}

func metrics(pinfos []*pieceInfo) {
	for i, pinfo := range pinfos {
		fmt.Printf("%v %v: %0.3f\t", pinfo.Piece, pinfo.Pos, pinfo.Weight)
		if i%4 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

type pieceInfo struct {
	Piece   pc.Piece
	Pos     game.Position
	IsBlack bool
	Weight  float64
}

const kingWeight float64 = 100
const checkPenalty float64 = 3
const earlyGameQueenPenalty float64 = 2

func getPieceWeight(g *game.GameState, pinfo *pieceInfo) float64 {
	if pinfo.Piece.IsKingLike() {
		pinfo.Weight = kingWeight + protectionWeight(g, pinfo)
		if g.IsAttacked(pinfo.Pos, pinfo.IsBlack) {
			pinfo.Weight -= checkPenalty
		}
		return pinfo.Weight
	}
	var pieceWeight float64 = 0
	if pinfo.Piece.IsQueenLike() {
		pieceWeight = 30 + rookMobility(g, pinfo) + bishopMobility(g, pinfo)
		if !isEndgame(g) {
			if pinfo.IsBlack {
				if pinfo.Pos.Row >= 4 {
					pieceWeight -= earlyGameQueenPenalty
				}
			} else {
				if pinfo.Pos.Row <= 4 {
					pieceWeight -= earlyGameQueenPenalty
				}
			}
		}
	} else if pinfo.Piece.IsRookLike() {
		pieceWeight = 7 + rookMobility(g, pinfo)
	} else if pinfo.Piece.IsBishopLike() {
		pieceWeight = bishopWeight(g, pinfo)
	} else if pinfo.Piece.IsHorsieLike() {
		pieceWeight = 3 + horsieMobility(g, pinfo)
	} else if pinfo.Piece.IsPawnLike() {
		pieceWeight = pawnWeight(g, pinfo)
	} else {
		panic("invalid piece: " + pinfo.Piece.String())
	}
	pinfo.Weight = pieceWeight
	return pinfo.Weight
}

func isEndgame(g *game.GameState) bool {
	return g.TotalValuablePieces <= 8
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
			weight += 0.15
		} else if piece.IsRookLike() {
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
			mobMod -= 0.05
			continue
		}
		piece := g.Board.AtPos(pos)
		if piece == pc.Empty {
			mobMod += 0.05
		} else if piece.IsBlack() != pinfo.IsBlack {
			mobMod += 0.08
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

func rookMobility(g *game.GameState, pinfo *pieceInfo) float64 {
	var mobMod float64 = 0
	for _, offset := range game.RookOffsets {
		for i := 0; i < 7; i++ {
			pos := game.Position{
				Column: pinfo.Pos.Column + (offset.Column * i),
				Row:    pinfo.Pos.Row + (offset.Row * i),
			}
			if pos.IsInvalid() {
				break
			}
			piece := g.Board.AtPos(pos)
			if piece == pc.Empty {
				mobMod += 0.05
			} else if piece.IsBlack() != pinfo.IsBlack {
				if piece.IsRookLike() {
					mobMod -= 0.2
				} else {
					mobMod += 0.15
				}
				break
			} else {
				mobMod += 0.1
				break
			}
		}
	}
	return mobMod
}

func bishopMobility(g *game.GameState, pinfo *pieceInfo) float64 {
	var mobMod float64 = 0
	for _, offset := range game.BishopOffsets {
		for i := 0; i < 7; i++ {
			pos := game.Position{
				Column: pinfo.Pos.Column + (offset.Column * i),
				Row:    pinfo.Pos.Row + (offset.Row * i),
			}
			if pos.IsInvalid() {
				break
			}
			piece := g.Board.AtPos(pos)
			if piece == pc.Empty {
				mobMod += 0.05
			} else if piece.IsBlack() != pinfo.IsBlack {
				if piece.IsBishopLike() {
					mobMod -= 0.2
				} else {
					mobMod += 0.15
				}
				break
			} else {
				mobMod += 0.1
				break
			}
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
	var out float64 = pawnColValue(g, pinfo) + pawnRowMod(g, pinfo)
	if hasFreeFront(g, pinfo) {
		out += 0.25
	}
	if isConnectedPawn(g, pinfo) {
		out += 0.1
	}
	return out
}

func isConnectedPawn(g *game.GameState, pinfo *pieceInfo) bool {
	var left, right game.Position
	var pawn, passantPawn, movedPawn pc.Piece
	if pinfo.IsBlack {
		left = game.Position{
			Column: pinfo.Pos.Column - 1,
			Row:    pinfo.Pos.Row - 1,
		}
		right = game.Position{
			Column: pinfo.Pos.Column - 1,
			Row:    pinfo.Pos.Row + 1,
		}
		pawn = pc.BlackPawn
		passantPawn = pc.BlackPassantPawn
		passantPawn = pc.BlackMovedPawn
	} else {
		left = game.Position{
			Column: pinfo.Pos.Column + 1,
			Row:    pinfo.Pos.Row - 1,
		}
		right = game.Position{
			Column: pinfo.Pos.Column + 1,
			Row:    pinfo.Pos.Row + 1,
		}
		pawn = pc.WhitePawn
		passantPawn = pc.WhitePassantPawn
		passantPawn = pc.WhiteMovedPawn
	}
	if left.IsValid() {
		leftpiece := g.Board.AtPos(left)
		if leftpiece == pawn ||
			leftpiece == passantPawn ||
			leftpiece == movedPawn {
			return true
		}
	}
	if right.IsValid() {
		rightpiece := g.Board.AtPos(right)
		if rightpiece == pawn ||
			rightpiece == passantPawn ||
			rightpiece == movedPawn {
			return true
		}
	}
	return false
}

func hasFreeFront(g *game.GameState, pinfo *pieceInfo) bool {
	if pinfo.IsBlack {
		for i := pinfo.Pos.Row + 1; i <= 7; i++ {
			if g.Board.At(i, pinfo.Pos.Column) != pc.Empty {
				return false
			}
		}
		return true
	}
	for i := pinfo.Pos.Row - 1; i >= 0; i-- {
		if g.Board.At(i, pinfo.Pos.Column) != pc.Empty {
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
		return 1.1
	case 4, 3:
		return 1.2
	}
	return 1
}

func pawnRowMod(g *game.GameState, pinfo *pieceInfo) float64 {
	if pinfo.IsBlack {
		switch pinfo.Pos.Column {
		case 1, 0:
			return -0.1
		case 3, 2:
			return -0.05
		case 5, 4:
			return 0.1
		case 7, 6:
			return 0.2
		}
	}
	switch pinfo.Pos.Column {
	case 1, 0:
		return 0.2
	case 3, 2:
		return 0.01
	case 5, 4:
		return -0.05
	case 7, 6:
		return -0.1
	}
	return 0
}

func Abs32(a int32) int32 {
	y := a >> 31
	return (a ^ y) - y
}

func Abs(a int) int {
	y := int32(a) >> 31
	return int((int32(a) ^ y) - y)
}
