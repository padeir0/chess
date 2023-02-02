package game

import (
	pc "chess/game/piece"
	pr "chess/game/promotion"
	"strconv"
)

func debugBoard(slots []*Slot) *Board {
	b := Board{}
	for i := 0; i < 64; i++ {
		b[i] = pc.Empty
	}
	b.SetPos(Position{Row: 0, Column: 4}, pc.BlackCastleKing)
	b.SetPos(Position{Row: 7, Column: 4}, pc.WhiteCastleKing)
	for _, slot := range slots {
		b.SetPos(slot.Position, slot.Piece)
	}
	return &b
}

func InitialBoard() *Board {
	return &Board{
		pc.BlackRook, pc.BlackHorsie, pc.BlackBishop, pc.BlackQueen, pc.BlackCastleKing, pc.BlackBishop, pc.BlackHorsie, pc.BlackRook,
		pc.BlackPawn, pc.BlackPawn, pc.BlackPawn, pc.BlackPawn, pc.BlackPawn, pc.BlackPawn, pc.BlackPawn, pc.BlackPawn,
		pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty,
		pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty,
		pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty,
		pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty,
		pc.WhitePawn, pc.WhitePawn, pc.WhitePawn, pc.WhitePawn, pc.WhitePawn, pc.WhitePawn, pc.WhitePawn, pc.WhitePawn,
		pc.WhiteRook, pc.WhiteHorsie, pc.WhiteBishop, pc.WhiteQueen, pc.WhiteCastleKing, pc.WhiteBishop, pc.WhiteHorsie, pc.WhiteRook,
	}
}

type Board [64]pc.Piece

func (this *Board) String() string {
	output := "    a  b  c  d  e  f  g  h  \n"
	for i := 0; i < 8; i++ {
		row := " " + strconv.Itoa(8-i) + " "
		output += row
		for j := 0; j < 8; j++ {
			if (i+j)%2 == 0 {
				output += "\u001b[43m"
			} else {
				output += "\u001b[41m"
			}
			if (*this)[i*8+j].IsBlack() {
				output += "\u001b[30m"
			}
			output += " \u001b[1m" + (*this)[i*8+j].String() + " \u001b[0m"
		}
		output += row + "\n"
	}
	output += "    a  b  c  d  e  f  g  h  \n"
	return output
}

func (this *Board) AtPos(pos Position) pc.Piece {
	return (*this)[pos.Column+8*pos.Row]
}

func (this *Board) At(row, column int) pc.Piece {
	return (*this)[column+8*row]
}

func (this *Board) SetPos(pos Position, s pc.Piece) {
	(*this)[pos.Column+8*pos.Row] = s
}
func (this *Board) Pop(pos Position) pc.Piece {
	i := pos.Column + 8*pos.Row
	ret := (*this)[i]
	(*this)[i] = pc.Empty
	return ret
}

type Slot struct {
	Piece    pc.Piece
	Position Position
}

type Position struct {
	Row    int // 1 2 3 4 5 6 7 8
	Column int // a b c d e f g h
}

func (this Position) IsValid() bool {
	return this.Column >= 0 && this.Column <= 7 &&
		this.Row >= 0 && this.Row <= 7
}

func (this Position) IsInvalid() bool {
	return this.Column < 0 || this.Column > 7 ||
		this.Row < 0 || this.Row > 7
}

func (this Position) String() string {
	col := rune(this.Column) + 'a'
	row := (7 - rune(this.Row)) + '1'
	return string(col) + string(row)
}

func InitialGame() *GameState {
	game := &GameState{
		BlackTurn: false,
		Board: *debugBoard([]*Slot{
			{Position: Position{7, 1}, Piece: pc.WhiteRook},
		}),

		BlackKingPosition: Position{Row: 0, Column: 4},
		WhiteKingPosition: Position{Row: 7, Column: 4},

		BlackPieces: []*Slot{},
		WhitePieces: []*Slot{},
	}

	for i, piece := range game.Board {
		if piece != pc.Empty {
			position := Position{Column: i % 8, Row: i / 8}
			slot := &Slot{Piece: piece, Position: position}
			if piece.IsWhite() {
				game.WhitePieces = append(game.WhitePieces, slot)
			} else {
				game.BlackPieces = append(game.BlackPieces, slot)
			}
		}
	}

	return game
}

type GameState struct {
	BlackTurn bool
	Board     Board

	// to check for Checks
	BlackKingPosition Position
	WhiteKingPosition Position

	BlackPieces []*Slot
	WhitePieces []*Slot
}

func (this *GameState) Copy() *GameState {
	output := &GameState{
		BlackTurn:         this.BlackTurn,
		Board:             this.Board,
		BlackKingPosition: this.BlackKingPosition,
		WhiteKingPosition: this.WhiteKingPosition,
		BlackPieces:       make([]*Slot, len(this.BlackPieces)),
		WhitePieces:       make([]*Slot, len(this.WhitePieces)),
	}
	for i, slot := range this.BlackPieces {
		new := *slot
		output.BlackPieces[i] = &new
	}
	for i, slot := range this.WhitePieces {
		new := *slot
		output.WhitePieces[i] = &new
	}
	return output
}

// returns if the move was sucessful
func (g *GameState) Move(from, to Position, prom pr.Promotion) (bool, *Slot) {
	fromPiece := g.Board.AtPos(from)
	if fromPiece.IsWhite() == g.BlackTurn {
		return false, nil
	}

	ok, capture := g.IsValidMove(from, to)
	if !ok {
		return false, nil
	}

	if g.PutsKingInCheck(from, to, capture) {
		return false, nil
	}

	piece := g.Board.AtPos(from)
	canPromote := canPromote(g, piece, to)
	if canPromote && prom == pr.InvalidPromotion {
		return false, nil
	}

	// now we know the move is valid
	g.RemovePassantPawns()

	if capture != nil {
		g.Board.Pop(capture.Position)
	}

	newPiece := pc.InvalidPiece
	if canPromote {
		newPiece = promotionToPiece(piece.IsBlack(), prom)
		g.Board.SetPos(from, pc.Empty)
		g.Board.SetPos(to, newPiece)
	} else {
		newPiece = alterPieceState(piece, from, to)
		g.Board.SetPos(from, pc.Empty)
		g.Board.SetPos(to, newPiece)
	}

	g.UpdatePieceTable(newPiece, capture, from, to)

	g.BlackTurn = !g.BlackTurn
	return true, capture
}

// returns if its valid and the position of the captured piece, if any
func (g *GameState) IsValidMove(from, to Position) (bool, *Slot) {
	if from.IsInvalid() || to.IsInvalid() {
		return false, nil
	}
	fromPiece := g.Board.AtPos(from)
	if fromPiece == pc.Empty {
		return false, nil
	}

	switch fromPiece {
	case pc.BlackCastleKing, pc.WhiteCastleKing: // can castle
		return isValidKingMove(g, from, to)
	case pc.BlackKing, pc.WhiteKing: // can't castle
		return isValidMovedKingMove(g, from, to)

	case pc.BlackHorsie, pc.WhiteHorsie:
		return isValidHorsieMove(g, from, to)

	case pc.BlackQueen, pc.WhiteQueen:
		return isValidQueenMove(g, from, to)

	case pc.BlackBishop, pc.WhiteBishop:
		return isValidBishopMove(g, from, to)

	case pc.BlackRook, pc.WhiteRook,
		pc.BlackMovedRook, pc.WhiteMovedRook:
		return isValidRookMove(g, from, to)

	case pc.BlackPassantPawn, pc.BlackMovedPawn:
		return isValidBlackMovedPawnMove(g, from, to)
	case pc.BlackPawn:
		return isValidBlackPawnMove(g, from, to)

	case pc.WhitePassantPawn, pc.WhiteMovedPawn:
		return isValidWhiteMovedPawnMove(g, from, to)
	case pc.WhitePawn:
		return isValidWhitePawnMove(g, from, to)
	}
	panic("oh no!")
}

func (this *GameState) UpdatePieceTable(piece pc.Piece, capture *Slot, from, to Position) {
	if piece == pc.BlackCastleKing || piece == pc.BlackKing {
		this.BlackKingPosition = to
	}
	if piece == pc.WhiteCastleKing || piece == pc.WhiteKing {
		this.WhiteKingPosition = to
	}
	if piece.IsWhite() {
		// update moved piece
		for i, slot := range this.WhitePieces {
			if slot != nil && slot.Position == from {
				this.WhitePieces[i] = &Slot{piece, to}
				break
			}
		}
		// update capture
		if capture != nil {
			for i, slot := range this.BlackPieces {
				if slot != nil && slot.Position == capture.Position {
					this.BlackPieces[i] = nil
					break
				}
			}
		}
		return
	}
	if piece.IsBlack() {
		// update moved piece
		for i, slot := range this.BlackPieces {
			if slot != nil && slot.Position == from {
				this.BlackPieces[i] = &Slot{piece, to}
				break
			}
		}
		// update capture
		if capture != nil {
			for i, slot := range this.WhitePieces {
				if slot != nil && slot.Position == capture.Position {
					this.WhitePieces[i] = nil
					break
				}
			}
		}
		return
	}
}

func (this *GameState) RemovePassantPawns() {
	for i, slot := range this.WhitePieces {
		if slot != nil && slot.Piece == pc.WhitePassantPawn {
			this.WhitePieces[i].Piece = pc.WhiteMovedPawn
			this.Board.SetPos(slot.Position, pc.WhiteMovedPawn)
		}
	}
	for i, slot := range this.BlackPieces {
		if slot != nil && slot.Piece == pc.BlackPassantPawn {
			this.BlackPieces[i].Piece = pc.BlackMovedPawn
			this.Board.SetPos(slot.Position, pc.BlackMovedPawn)
		}
	}
}

func (this *GameState) IsAttacked(pos Position, isBlack bool) bool {
	attackers := &this.BlackPieces
	if isBlack {
		attackers = &this.WhitePieces
	}

	for _, slot := range *attackers {
		if slot == nil {
			continue
		}
		ok, capture := this.IsValidMove(slot.Position, pos)
		if ok && capture != nil && capture.Position == pos {
			return true
		}
	}
	return false
}

func (this *GameState) PutsKingInCheck(from, to Position, capture *Slot) bool {
	// perform move
	piece := this.Board.Pop(from)
	if capture != nil {
		this.Board.Pop(capture.Position)
	}
	this.Board.SetPos(to, piece)

	answer := false
	if piece.IsKingLike() {
		answer = this.IsAttacked(to, this.BlackTurn)
	} else if this.BlackTurn {
		answer = this.IsAttacked(this.BlackKingPosition, true)
	} else {
		answer = this.IsAttacked(this.WhiteKingPosition, false)
	}

	// undo move
	this.Board.Pop(to)
	if capture != nil {
		this.Board.SetPos(capture.Position, capture.Piece)
	}
	this.Board.SetPos(from, piece)
	return answer
}

func isValidBlackMovedPawnMove(g *GameState, from, to Position) (bool, *Slot) {
	// can only move 2 steps foward one time
	if to.Row-from.Row != 1 {
		return false, nil
	}
	return isValidBlackPawnMove(g, from, to)
}

/*
  P
 ###
  #
*/
func isValidBlackPawnMove(g *GameState, from, to Position) (bool, *Slot) {
	// check shape of movement
	if to.Row-from.Row > 2 || to.Row <= from.Row ||
		to.Column-from.Column > 1 || to.Column-from.Column < -1 {
		return false, nil
	}
	return isValidPawnMove(g, from, to)
}

func isValidWhiteMovedPawnMove(g *GameState, from, to Position) (bool, *Slot) {
	// can only move 2 steps foward one time
	if to.Row-from.Row != -1 {
		return false, nil
	}
	return isValidWhitePawnMove(g, from, to)
}

/*
  #
 ###
  P
*/
func isValidWhitePawnMove(g *GameState, from, to Position) (bool, *Slot) {
	// check shape of movement
	if to.Row-from.Row < -2 || to.Row >= from.Row ||
		to.Column-from.Column > 1 || to.Column-from.Column < -1 {
		return false, nil
	}
	return isValidPawnMove(g, from, to)
}

func isValidPawnMove(g *GameState, from, to Position) (bool, *Slot) {
	var capture *Slot = nil
	fromPiece := g.Board.AtPos(from)
	toPiece := g.Board.AtPos(to)
	// rook-like move
	if to.Column == from.Column {
		pos := rook_ClosestPieceInWay(g, from, to)
		if pos != nil {
			// can't capture forwards
			return false, nil
		}
	} else {
		if toPiece == pc.Empty {
			// check for en passant
			sidePiece := g.Board.At(from.Row, to.Column)
			if (sidePiece != pc.BlackPassantPawn &&
				sidePiece != pc.WhitePassantPawn) ||
				sidePiece.IsBlack() == fromPiece.IsBlack() {
				// no en passant on side
				return false, nil
			}
			capture = &Slot{
				Piece:    sidePiece,
				Position: Position{Row: from.Row, Column: to.Column},
			}
		} else {
			// capturing a piece
			if toPiece.IsBlack() != fromPiece.IsBlack() {
				capture = &Slot{
					Piece:    toPiece,
					Position: to,
				}
			} else {
				// capturing friend
				return false, nil
			}
		}
	}

	return true, capture
}

func isValidQueenMove(g *GameState, from, to Position) (bool, *Slot) {
	ok, slot := isValidBishopMove(g, from, to)
	if ok {
		return ok, slot
	}
	ok, slot = isValidRookMove(g, from, to)
	if ok {
		return ok, slot
	}
	return false, nil
}

func isValidRookMove(g *GameState, from, to Position) (bool, *Slot) {
	// check if shape is rook-like
	if from.Column != to.Column && from.Row != to.Row {
		return false, nil
	}
	posInWay := rook_ClosestPieceInWay(g, from, to)

	// something in the way
	if posInWay != nil && *posInWay != to {
		return false, nil
	}

	fromPiece := g.Board.AtPos(from)
	toPiece := g.Board.AtPos(to)
	var capture *Slot = nil
	if posInWay != nil && *posInWay == to {
		if fromPiece.IsBlack() == toPiece.IsBlack() {
			// friendly piece in spot
			return false, nil
		}
		capture = &Slot{
			Piece:    toPiece,
			Position: to,
		}
	}
	// move is valid ---

	return true, capture
}

func isValidBishopMove(g *GameState, from, to Position) (bool, *Slot) {
	// check if shape is bishop-like
	if Abs(int32(from.Column-to.Column)) != Abs(int32(from.Row-to.Row)) {
		return false, nil
	}

	posInWay := bishop_ClosestPieceInWay(g, from, to)
	if posInWay != nil && *posInWay != to {
		// piece in way
		return false, nil
	}

	fromPiece := g.Board.AtPos(from)
	toPiece := g.Board.AtPos(to)
	var capture *Slot = nil
	if posInWay != nil && *posInWay == to {
		if fromPiece.IsBlack() == toPiece.IsBlack() {
			// friendly piece in spot
			return false, nil
		}
		capture = &Slot{
			Piece:    toPiece,
			Position: to,
		}
	}

	return true, capture
}

func isValidKingMove(g *GameState, from, to Position) (bool, *Slot) {
	// TODO: validate castling
	return isValidMovedKingMove(g, from, to)
}

/*
 ### (-1, -1) (-1, 0) (-1, 1)
 #K# (0,  -1)         (0,  1)
 ### (1,  -1) (1,  0) (1,  1)
*/
func isValidMovedKingMove(g *GameState, from, to Position) (bool, *Slot) {
	ColDiff := from.Column - to.Column
	RowDiff := from.Row - to.Row
	if !(((ColDiff == 1) || (ColDiff == 0) || (ColDiff == -1)) &&
		((RowDiff == 1) || (RowDiff == 0) || (RowDiff == -1))) {
		return false, nil
	}

	fromPiece := g.Board.AtPos(from)
	toPiece := g.Board.AtPos(to)
	var capture *Slot = nil
	if toPiece != pc.Empty {
		if fromPiece.IsBlack() == toPiece.IsBlack() {
			// friendly piece in way
			return false, nil
		}
		capture = &Slot{
			Piece:    toPiece,
			Position: to,
		}
	}

	return true, capture
}

/*
 # #      (r-2, c-1) (r-2, c+1)
#   #  (r-1, c-2)      (r-1, c+2)
  H              (r, c)
#   #  (r+1, c-2)     (r+1, c+2)
 # #     (r+2, c-1) (r+2, c+1)
*/
func isValidHorsieMove(g *GameState, from, to Position) (bool, *Slot) {
	// check if shape is horsie-like
	if !((Abs(int32(from.Column-to.Column)) == 2 &&
		Abs(int32(from.Row-to.Row)) == 1) ||
		(Abs(int32(from.Column-to.Column)) == 1 &&
			Abs(int32(from.Row-to.Row)) == 2)) {
		return false, nil
	}

	fromPiece := g.Board.AtPos(from)
	toPiece := g.Board.AtPos(to)
	var capture *Slot = nil
	if toPiece != pc.Empty {
		if fromPiece.IsBlack() == toPiece.IsBlack() {
			// friendly piece in way
			return false, nil
		}
		capture = &Slot{
			Piece:    toPiece,
			Position: to,
		}
	}

	return true, capture
}

/*
diagonals only (bishop-like)
 #   #  (0, 1)        (0, 5)
  # #      (1, 2) (1, 4)
   B          (2, 3)
  # #      (3, 2) (3, 4)
 #   #  (4, 1)        (4, 5)
*/
func bishop_ClosestPieceInWay(g *GameState, from, to Position) *Position {
	rowQuant := 1
	if from.Row > to.Row {
		rowQuant = -1
	}
	colQuant := 1
	if from.Column > to.Column {
		colQuant = -1
	}

	currPos := Position{
		Row:    from.Row + rowQuant,
		Column: from.Column + colQuant,
	}
	destPos := Position{
		Row:    to.Row + rowQuant,
		Column: to.Column + colQuant,
	}
	for currPos != destPos {
		piece := g.Board.AtPos(currPos)
		if piece != pc.Empty {
			return &currPos
		}
		currPos.Row += rowQuant
		currPos.Column += colQuant
	}
	return nil
}

/*
col + row only (rook-like)
.   #
.   #
. ##R##
.   #
.   #
*/
func rook_ClosestPieceInWay(g *GameState, from, to Position) *Position {
	if from.Column != to.Column {
		quant := 1
		if from.Column > to.Column {
			quant = -1
		}
		for i := from.Column + quant; i != to.Column+quant; i += quant {
			if g.Board.At(from.Row, i) != pc.Empty {
				return &Position{
					Row:    from.Row,
					Column: i,
				}
			}
		}
	}
	if from.Row != to.Row {
		quant := 1
		if from.Row > to.Row {
			quant = -1
		}
		for i := from.Row + quant; i != to.Row+quant; i += quant {
			if g.Board.At(i, from.Column) != pc.Empty {
				return &Position{
					Row:    i,
					Column: from.Column,
				}
			}
		}
	}
	return nil
}

func canPromote(g *GameState, piece pc.Piece, to Position) bool {
	return (piece == pc.BlackPawn && to.Row == 7) ||
		(piece == pc.WhitePawn && to.Row == 0)
}

func promotionToPiece(isBlack bool, prom pr.Promotion) pc.Piece {
	if isBlack {
		switch prom {
		case pr.Horsie:
			return pc.BlackHorsie
		case pr.Rook:
			return pc.BlackRook
		case pr.Queen:
			return pc.BlackQueen
		case pr.Bishop:
			return pc.BlackBishop
		}
	}
	switch prom {
	case pr.Horsie:
		return pc.WhiteHorsie
	case pr.Rook:
		return pc.WhiteRook
	case pr.Queen:
		return pc.WhiteQueen
	case pr.Bishop:
		return pc.WhiteBishop
	}
	panic("no horsequeens")
}

// move must be valid
func alterPieceState(piece pc.Piece, from, to Position) pc.Piece {
	if piece == pc.BlackPawn && from.Row-to.Row == -2 {
		return pc.BlackPassantPawn
	}
	if piece == pc.BlackPawn {
		return pc.BlackMovedPawn
	}
	if piece == pc.WhitePawn && from.Row-to.Row == 2 {
		return pc.WhitePassantPawn
	}
	if piece == pc.WhitePawn {
		return pc.WhiteMovedPawn
	}
	if piece == pc.WhiteRook {
		return pc.WhiteMovedRook
	}
	if piece == pc.BlackRook {
		return pc.BlackMovedRook
	}
	return piece
}

func Abs(a int32) int32 {
	y := a >> 31
	return (a ^ y) - y
}

/*
Horsequeen
#  #  #
 #####
 #####
###P###
 #####
 #####
#  #  #
*/
