package custom

import (
	. "chess/evals/common"
	"chess/game"
	pc "chess/game/piece"
	rs "chess/game/result"
	ifaces "chess/interfaces"
)

var _ ifaces.Evaluator = Evaluate

// evaluates material, position and mobility:
//     pawn structure
//         connected pawns
//     mobility
//         horsie, queen, bishop, rook, pawn
//     piece value
//         bishop pair
//     king safety
//         pawn and piece blockade
//
// all values on centipawns, using integers

// maximize for white
func Evaluate(g *game.GameState, depth int) int {
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
		pinfo := &PieceInfo{
			Piece:   slot.Piece,
			Pos:     slot.Pos,
			IsBlack: false,
		}
		total += getPieceWeight(g, pinfo) + GetPositionalWeight(isEndgame(g), false, slot.Piece, slot.Pos)
	}
	for _, slot := range g.BlackPieces {
		if slot.IsInvalid() {
			continue
		}
		pinfo := &PieceInfo{
			Piece:   slot.Piece,
			Pos:     slot.Pos,
			IsBlack: true,
		}
		total -= getPieceWeight(g, pinfo) + GetPositionalWeight(isEndgame(g), true, slot.Piece, slot.Pos)
	}
	return total
}

type PieceInfo struct {
	Piece   pc.Piece
	Pos     game.Point
	IsBlack bool
	Weight  int
}

const kingWeight int = 10000

func getPieceWeight(g *game.GameState, pinfo *PieceInfo) int {
	if pinfo.Piece.IsKingLike() {
		if isEndgame(g) {
			pinfo.Weight = kingWeight
			return pinfo.Weight
		}
		pinfo.Weight = kingWeight + protectionWeight(g, pinfo)
		return pinfo.Weight
	}
	var pieceWeight int = 0
	if pinfo.Piece.IsQueenLike() {
		pieceWeight = 1500 + queenMobility(g, pinfo)
	} else if pinfo.Piece.IsRookLike() {
		pieceWeight = 700 + rookMobility(g, pinfo)
	} else if pinfo.Piece.IsBishopLike() {
		pieceWeight = 300 + bishopMobility(g, pinfo)
	} else if pinfo.Piece.IsKnightLike() {
		pieceWeight = 300 + horsieMobility(g, pinfo)
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

func protectionWeight(g *game.GameState, pinfo *PieceInfo) int {
	var weight int = 1
	for _, offset := range game.KingOffsets {
		pos := game.Point{
			Column: pinfo.Pos.Column + offset.Column,
			Row:    pinfo.Pos.Row + offset.Row,
		}
		if pos.IsInvalid() {
			continue
		}
		piece := g.Board.AtPos(pos)
		if piece != pc.Empty {
			if piece.IsBlack() == pinfo.IsBlack {
				if piece.IsPawnLike() {
					weight += 3
				} else {
					weight += 1
				}
			} else {
				weight -= 5
			}
		}
	}
	return weight
}

func horsieMobility(g *game.GameState, pinfo *PieceInfo) int {
	var mobMod int = 0
	for _, offset := range game.HorsieOffsets {
		pos := game.Point{
			Column: pinfo.Pos.Column + offset.Column,
			Row:    pinfo.Pos.Row + offset.Row,
		}
		if pos.IsInvalid() {
			mobMod -= 5
			continue
		}
		piece := g.Board.AtPos(pos)
		if piece != pc.Empty {
			// this paints horsies as support pieces
			if piece.IsBlack() != pinfo.IsBlack { // attacking
				mobMod += 8
			} else { // defending
				mobMod += 4
			}
		}
	}
	return mobMod
}

func queenMobility(g *game.GameState, pinfo *PieceInfo) int {
	return rookMobility(g, pinfo)/2 + bishopMobility(g, pinfo)/2
}

func rookMobility(g *game.GameState, pinfo *PieceInfo) int {
	var mobMod int = 0
	for _, offset := range game.RookOffsets {
		for i := 1; i < 7; i++ {
			pos := game.Point{
				Column: pinfo.Pos.Column + (offset.Column * i),
				Row:    pinfo.Pos.Row + (offset.Row * i),
			}
			if pos.IsInvalid() {
				break
			}
			piece := g.Board.AtPos(pos)
			if piece == pc.Empty {
				mobMod += 5
			} else if piece.IsBlack() != pinfo.IsBlack {
				if piece.IsRookLike() {
					mobMod -= 20
				} else {
					mobMod += 15
				}
				break
			} else {
				break
			}
		}
	}
	return mobMod
}

func bishopMobility(g *game.GameState, pinfo *PieceInfo) int {
	var mobMod int = 0
	for _, offset := range game.BishopOffsets {
		for i := 1; i < 7; i++ {
			pos := game.Point{
				Column: pinfo.Pos.Column + (offset.Column * i),
				Row:    pinfo.Pos.Row + (offset.Row * i),
			}
			if pos.IsInvalid() {
				break
			}
			piece := g.Board.AtPos(pos)
			if piece == pc.Empty {
				mobMod += 5
			} else if piece.IsBlack() != pinfo.IsBlack {
				if piece.IsBishopLike() {
					mobMod -= 20
				} else {
					mobMod += 15
				}
				break
			} else {
				break
			}
		}
	}
	return mobMod
}

// inneficient but will do for now
func hasBishopPair(g *game.GameState, pinfo *PieceInfo) bool {
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

func pawnWeight(g *game.GameState, pinfo *PieceInfo) int {
	var out int = 100
	if hasFreeFront(g, pinfo) {
		out += 25
	}
	if !isConnectedPawn(g, pinfo) {
		out -= 10
	}
	return out
}

func isConnectedPawn(g *game.GameState, pinfo *PieceInfo) bool {
	var left, right game.Point
	var pawn pc.Piece
	if pinfo.IsBlack {
		left = game.Point{
			Column: pinfo.Pos.Column - 1,
			Row:    pinfo.Pos.Row - 1,
		}
		right = game.Point{
			Column: pinfo.Pos.Column - 1,
			Row:    pinfo.Pos.Row + 1,
		}
		pawn = pc.BlackPawn
	} else {
		left = game.Point{
			Column: pinfo.Pos.Column + 1,
			Row:    pinfo.Pos.Row - 1,
		}
		right = game.Point{
			Column: pinfo.Pos.Column + 1,
			Row:    pinfo.Pos.Row + 1,
		}
		pawn = pc.WhitePawn
	}
	if left.IsValid() {
		leftpiece := g.Board.AtPos(left)
		if leftpiece == pawn {
			return true
		}
	}
	if right.IsValid() {
		rightpiece := g.Board.AtPos(right)
		if rightpiece == pawn {
			return true
		}
	}
	return false
}

func hasFreeFront(g *game.GameState, pinfo *PieceInfo) bool {
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
