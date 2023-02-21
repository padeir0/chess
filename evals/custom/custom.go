package custom

import (
	"chess/evals/common"
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
		total += getPieceWeight(g, pinfo) + getPositionalWeight(isEndgame(g), pinfo)
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
		total -= getPieceWeight(g, pinfo) + getPositionalWeight(isEndgame(g), pinfo)
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
					weight += 5
				} else if piece.IsRookLike() {
					weight += 3
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
		if piece == pc.Empty {
			if pinfo.IsBlack {
				if inKingRegion(g.WhiteKingPosition, pos) {
					mobMod += 25
				}
			} else {
				if inKingRegion(g.BlackKingPosition, pos) {
					mobMod += 25
				}
			}
		}
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
				if pinfo.IsBlack {
					if inKingRegion(g.WhiteKingPosition, pos) {
						mobMod += 15
					}
				} else {
					if inKingRegion(g.BlackKingPosition, pos) {
						mobMod += 15
					}
				}
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
				if pinfo.IsBlack {
					if inKingRegion(g.WhiteKingPosition, pos) {
						mobMod += 35
					}
				} else {
					if inKingRegion(g.BlackKingPosition, pos) {
						mobMod += 35
					}
				}
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

// favor pawns close to promotion
// prefer pawn blockade on king side
var pawn_md_psqt = common.PieceSquareTable{
	0, 0, 0, 0, 0, 0, 0, 0,
	100, 125, 90, 70, 70, 90, 125, 100,
	80, 95, 60, 40, 40, 60, 95, 80,
	50, 13, 6, 21, 23, 12, 17, -23,
	20, -2, -5, 20, 20, 6, 10, -25,
	0, -4, -4, -10, 3, 3, 33, -12,
	-10, -1, -20, -23, -15, 24, 38, -22,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// favor pawns close to promotion, preferably on edges
var pawn_ed_psqt = common.PieceSquareTable{
	0, 0, 0, 0, 0, 0, 0, 0,
	178, 173, 158, 134, 147, 132, 165, 187,
	94, 100, 85, 67, 56, 53, 82, 84,
	32, 24, 13, 5, -2, 4, 17, 17,
	13, 9, -3, -7, -7, -8, 3, -1,
	4, 7, -6, 1, 0, -5, -1, -8,
	13, 8, 8, 10, 13, 0, 2, -7,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// avoid puting knight on edge
// protect center squares
var knight_psqt = common.PieceSquareTable{
	-30, -5, -5, -5, -5, -5, -5, -30,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-15, 0, 0, 0, 0, 0, 0, -15,
	-20, 0, 0, 0, 0, 0, 0, -20,
	-20, 0, 0, 0, 0, 0, 0, -20,
	-20, 0, 20, 0, 0, 20, 0, -20,
	-20, 0, 0, 0, 0, 0, 0, -20,
	-50, -10, -10, -10, -10, -10, -10, -50,
}

// protect king behind pawns
var king_md_psqt = common.PieceSquareTable{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	-20, -10, -7, -5, 5, 20, 30, 20,
}

// centralize king on endgame
var king_ed_psqt = common.PieceSquareTable{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 10, 30, 30, 10, 0, 0,
	0, 0, 10, 30, 30, 10, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// favor controling the center and flanks
var rook_psqt = common.PieceSquareTable{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	10, 20, 20, 10, 10, 20, 20, 10,
	10, 20, 20, 10, 10, 20, 20, 10,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// favor controling the center from far
var bishop_psqt = common.PieceSquareTable{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 5, 5, 0, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	10, 20, 20, 10, 10, 20, 20, 10,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

func getPositionalWeight(isEndgame bool, pinfo *PieceInfo) int {
	pos := pinfo.Pos
	if pinfo.IsBlack {
		pos = common.Mirror(pinfo.Pos)
	}
	if isEndgame {
		switch pinfo.Piece {
		case pc.WhitePawn, pc.BlackPawn:
			return pawn_ed_psqt.AtPos(pos)
		case pc.BlackKing, pc.WhiteKing:
			return king_ed_psqt.AtPos(pos)
		case pc.BlackKnight, pc.WhiteKnight:
			return knight_psqt.AtPos(pos)
		case pc.BlackRook, pc.WhiteRook:
			return rook_psqt.AtPos(pos)
		case pc.BlackBishop, pc.WhiteBishop:
			return bishop_psqt.AtPos(pos)
		}
	}
	switch pinfo.Piece {
	case pc.WhitePawn, pc.BlackPawn:
		return pawn_md_psqt.AtPos(pos)
	case pc.BlackKing, pc.WhiteKing:
		return king_md_psqt.AtPos(pos)
	case pc.BlackKnight, pc.WhiteKnight:
		return knight_psqt.AtPos(pos)
	case pc.BlackRook, pc.WhiteRook:
		return rook_psqt.AtPos(pos)
	case pc.BlackBishop, pc.WhiteBishop:
		return bishop_psqt.AtPos(pos)
	}
	return 0
}

func inKingRegion(kingPos, otherPos game.Point) bool {
	for _, offset := range game.KingOffsets {
		pos := game.Point{
			Column: kingPos.Column + offset.Column,
			Row:    kingPos.Row + offset.Row,
		}
		if pos == otherPos {
			return true
		}
	}
	return false
}
