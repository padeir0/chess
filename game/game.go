package game

import (
	ac "chess/asciicolors"
	pc "chess/game/piece"
	pr "chess/game/promotion"

	"sort"
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
		b.SetPos(slot.Pos, slot.Piece)
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
				output += ac.BackgroundYellow
			} else {
				output += ac.BackgroundRed
			}
			if (*this)[i*8+j].IsBlack() {
				output += ac.Black
			}
			output += " " + ac.Bold + (*this)[i*8+j].String() + " " + ac.Reset
		}
		output += row + "\n"
	}
	output += "    a  b  c  d  e  f  g  h  \n"
	return output
}

type Highlight struct {
	Pos   Position
	Color ac.Color
}

func (this *Board) Show(hls []Highlight) string {
	m := map[int]ac.Color{}
	for _, hl := range hls {
		m[hl.Pos.Column+8*hl.Pos.Row] = hl.Color
	}

	output := "    a  b  c  d  e  f  g  h  \n"
	for i := 0; i < 8; i++ {
		row := " " + strconv.Itoa(8-i) + " "
		output += row
		for j := 0; j < 8; j++ {
			if clr, ok := m[i*8+j]; ok {
				output += clr
			} else if (i+j)%2 == 0 {
				output += ac.BackgroundYellow
			} else {
				output += ac.BackgroundRed
			}
			if (*this)[i*8+j].IsBlack() {
				output += ac.Black
			}
			output += " " + ac.Bold + (*this)[i*8+j].String() + " " + ac.Reset
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
	Piece pc.Piece
	Pos   Position
}

func (this Slot) IsValid() bool {
	return this.Piece != pc.Empty && this.Piece != pc.InvalidPiece
}

func (this Slot) IsInvalid() bool {
	return this.Piece == pc.Empty || this.Piece == pc.InvalidPiece
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
		Board:     *InitialBoard(),

		BlackKingPosition: Position{Row: 0, Column: 4},
		WhiteKingPosition: Position{Row: 7, Column: 4},

		BlackPieces: []Slot{},
		WhitePieces: []Slot{},
	}

	for i, piece := range game.Board {
		if piece != pc.Empty {
			position := Position{Column: i % 8, Row: i / 8}
			slot := Slot{Piece: piece, Pos: position}
			if piece.IsWhite() {
				game.WhitePieces = append(game.WhitePieces, slot)
			} else {
				game.BlackPieces = append(game.BlackPieces, slot)
			}
			if !piece.IsPawnLike() {
				game.TotalValuablePieces += 1
			}
		}
	}

	orderByValue(game.WhitePieces)
	orderByValue(game.BlackPieces)

	return game
}

type GameState struct {
	BlackTurn bool
	Board     Board

	// to check for Checks
	BlackKingPosition Position
	WhiteKingPosition Position

	BlackPieces []Slot
	WhitePieces []Slot

	TotalValuablePieces int
}

// not to be used by the engine
func (this *GameState) IsOver() (bool, string) {
	if this.IsAttacked(this.BlackKingPosition, true) &&
		len(this.ValidMoves(true)) == 0 {
		return true, "white wins by checkmate"
	}
	if this.IsAttacked(this.WhiteKingPosition, false) &&
		len(this.ValidMoves(false)) == 0 {
		return true, "black wins by checkmate"
	}
	if len(this.ValidMoves(this.BlackTurn)) == 0 {
		return true, "draw by stalemate"
	}
	return false, ""
}

func (this *GameState) ShowMoves() string {
	moves := this.ValidMoves(this.BlackTurn)
	hls := MoveToHighlight(moves)
	return this.Board.Show(hls)
}

func MoveToHighlight(in []*Move) []Highlight {
	out := []Highlight{}
	for _, move := range in {
		hl := Highlight{
			Pos:   move.To,
			Color: ac.BackgroundGreen,
		}
		out = append(out, hl)
	}
	return out
}

func (this *GameState) Copy() *GameState {
	output := &GameState{
		BlackTurn:           this.BlackTurn,
		Board:               this.Board,
		BlackKingPosition:   this.BlackKingPosition,
		WhiteKingPosition:   this.WhiteKingPosition,
		BlackPieces:         make([]Slot, len(this.BlackPieces)),
		WhitePieces:         make([]Slot, len(this.WhitePieces)),
		TotalValuablePieces: this.TotalValuablePieces,
	}
	for i, slot := range this.BlackPieces {
		if slot.IsInvalid() {
			continue
		}
		output.BlackPieces[i] = slot
	}
	for i, slot := range this.WhitePieces {
		if slot.IsInvalid() {
			continue
		}
		output.WhitePieces[i] = slot
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

	if g.putsKingInCheck(from, to, capture) {
		return false, nil
	}

	piece := g.Board.AtPos(from)
	canPromote := canPromote(g, piece, to)
	if canPromote && prom == pr.InvalidPromotion {
		return false, nil
	}

	if capture != nil {
		g.Board.Pop(capture.Pos)
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

	g.updatePieceTable(newPiece, capture, from, to)

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

	case pc.BlackMovedPawn:
		return isValidBlackMovedPawnMove(g, from, to)
	case pc.BlackPawn:
		return isValidBlackPawnMove(g, from, to)

	case pc.WhiteMovedPawn:
		return isValidWhiteMovedPawnMove(g, from, to)
	case pc.WhitePawn:
		return isValidWhitePawnMove(g, from, to)
	}
	panic("oh no!")
}

func (this *GameState) ValidMoves(isBlack bool) []*Move {
	slots := this.WhitePieces
	if isBlack {
		slots = this.BlackPieces
	}
	outStates := []*Move{}
	newGS := this.Copy()
	for _, slot := range slots {
		if slot.IsInvalid() {
			continue
		}
		// Move() may alter the slot
		piece := slot.Piece
		moves := PseudoLegalMoves(slot.Pos, slot.Piece)
		for _, to := range moves {
			ok, capture := newGS.Move(slot.Pos, to, pr.Queen)
			if ok {
				move := &Move{
					Piece:   piece,
					From:    slot.Pos,
					To:      to,
					Capture: capture,
				}
				newGS.Unmake(piece, slot.Pos, to, capture)
				outStates = append(outStates, move)
			}
		}
	}
	return outStates
}

func (this *GameState) IsAttacked(pos Position, isBlack bool) bool {
	attackers := &this.BlackPieces
	if isBlack {
		attackers = &this.WhitePieces
	}

	for _, slot := range *attackers {
		if slot.IsInvalid() {
			continue
		}
		ok, capture := this.IsValidMove(slot.Pos, pos)
		if ok && capture != nil && capture.Pos == pos {
			return true
		}
	}
	return false
}

func (this *GameState) updatePieceTable(piece pc.Piece, capture *Slot, from, to Position) {
	if piece == pc.BlackCastleKing || piece == pc.BlackKing {
		this.BlackKingPosition = to
	}
	if piece == pc.WhiteCastleKing || piece == pc.WhiteKing {
		this.WhiteKingPosition = to
	}
	if capture != nil && !capture.Piece.IsPawnLike() {
		this.TotalValuablePieces -= 1
	}
	if piece.IsWhite() {
		// update moved piece
		for i, slot := range this.WhitePieces {
			if slot.IsValid() && slot.Pos == from {
				this.WhitePieces[i] = Slot{piece, to}
				break
			}
		}
		// update capture
		if capture != nil {
			for i, slot := range this.BlackPieces {
				if slot.IsValid() && slot.Pos == capture.Pos {
					this.BlackPieces[i] = Slot{pc.Empty, Position{0, 0}}
					break
				}
			}
		}
		return
	}
	if piece.IsBlack() {
		// update moved piece
		for i, slot := range this.BlackPieces {
			if slot.IsValid() && slot.Pos == from {
				this.BlackPieces[i] = Slot{piece, to}
				break
			}
		}
		// update capture
		if capture != nil {
			for i, slot := range this.WhitePieces {
				if slot.IsValid() && slot.Pos == capture.Pos {
					this.WhitePieces[i] = Slot{pc.Empty, Position{0, 0}}
					break
				}
			}
		}
		return
	}
}

func (this *GameState) unmakeTableUpdate(piece pc.Piece, capture *Slot, from, to Position) {
	if piece == pc.BlackCastleKing || piece == pc.BlackKing {
		this.BlackKingPosition = from
	}
	if piece == pc.WhiteCastleKing || piece == pc.WhiteKing {
		this.WhiteKingPosition = from
	}
	if capture != nil && !capture.Piece.IsPawnLike() {
		this.TotalValuablePieces += 1
	}
	if piece.IsWhite() {
		// update moved piece
		for i, slot := range this.WhitePieces {
			if slot.IsValid() && slot.Pos == to {
				this.WhitePieces[i] = Slot{piece, from}
				break
			}
		}
		// update capture
		if capture != nil {
			sl := Slot{capture.Piece, capture.Pos}
			for i, slot := range this.BlackPieces {
				if slot.IsInvalid() {
					this.BlackPieces[i] = sl
					return
				}
			}
			this.BlackPieces = append(this.BlackPieces, sl)
		}
		return
	}
	if piece.IsBlack() {
		// update moved piece
		for i, slot := range this.BlackPieces {
			if slot.IsValid() && slot.Pos == to {
				this.BlackPieces[i] = Slot{piece, from}
				break
			}
		}
		// update capture
		if capture != nil {
			sl := Slot{capture.Piece, capture.Pos}
			for i, slot := range this.WhitePieces {
				if slot.IsInvalid() {
					this.WhitePieces[i] = sl
					return
				}
			}
			this.WhitePieces = append(this.WhitePieces, sl)
		}
		return
	}
}

func (this *GameState) Unmake(piece pc.Piece, from, to Position, capture *Slot) {
	this.unmakeTableUpdate(piece, capture, from, to)
	this.Board.Pop(to)
	if capture != nil {
		this.Board.SetPos(capture.Pos, capture.Piece)
	}
	this.Board.SetPos(from, piece)
	this.BlackTurn = !this.BlackTurn
}

func (this *GameState) UnmakeMove(mv *Move) {
	this.unmakeTableUpdate(mv.Piece, mv.Capture, mv.From, mv.To)
	this.Board.Pop(mv.To)
	if mv.Capture != nil {
		this.Board.SetPos(mv.Capture.Pos, mv.Capture.Piece)
	}
	this.Board.SetPos(mv.From, mv.Piece)
	this.BlackTurn = !this.BlackTurn
}

func (this *GameState) putsKingInCheck(from, to Position, capture *Slot) bool {
	// perform move
	piece := this.Board.Pop(from)
	if capture != nil {
		this.Board.Pop(capture.Pos)
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
		this.Board.SetPos(capture.Pos, capture.Piece)
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
		_, ok := rook_ClosestPieceInWay(g, from, to)
		if ok {
			// can't capture forwards
			return false, nil
		}
	} else {
		if toPiece == pc.Empty {
			// no En Passant, sorry :)
			return false, nil
		} else {
			// capturing a piece
			if toPiece.IsBlack() != fromPiece.IsBlack() {
				capture = &Slot{
					Piece: toPiece,
					Pos:   to,
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
	posInWay, ok := rook_ClosestPieceInWay(g, from, to)

	// something in the way
	if ok && posInWay != to {
		return false, nil
	}

	fromPiece := g.Board.AtPos(from)
	toPiece := g.Board.AtPos(to)
	var capture *Slot = nil
	if ok && posInWay == to {
		if fromPiece.IsBlack() == toPiece.IsBlack() {
			// friendly piece in spot
			return false, nil
		}
		capture = &Slot{
			Piece: toPiece,
			Pos:   to,
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
			Piece: toPiece,
			Pos:   to,
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
			Piece: toPiece,
			Pos:   to,
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
			Piece: toPiece,
			Pos:   to,
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
func rook_ClosestPieceInWay(g *GameState, from, to Position) (Position, bool) {
	if from.Column != to.Column {
		quant := 1
		if from.Column > to.Column {
			quant = -1
		}
		for i := from.Column + quant; i != to.Column+quant; i += quant {
			if g.Board.At(from.Row, i) != pc.Empty {
				return Position{
					Row:    from.Row,
					Column: i,
				}, true
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
				return Position{
					Row:    i,
					Column: from.Column,
				}, true
			}
		}
	}
	return Position{}, false
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
	if piece == pc.BlackPawn {
		return pc.BlackMovedPawn
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

type Move struct {
	Piece   pc.Piece
	From    Position
	To      Position
	Capture *Slot
}

func (this *Move) String() string {
	if this.Capture != nil {
		return this.Piece.String() +
			" " + this.From.String() + this.To.String() +
			this.Capture.Piece.String()
	}
	return this.Piece.String() + " " +
		this.From.String() + this.To.String()
}

func PseudoLegalMoves(Pos Position, piece pc.Piece) []Position {
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
	case pc.BlackMovedPawn, pc.BlackPawn:
		return genBlackPawnMoves(Pos)
	case pc.WhiteMovedPawn, pc.WhitePawn:
		return genWhitePawnMoves(Pos)
	}
	return nil
}

var KingOffsets = []Position{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1} /*    */, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

func genKingMoves(pos Position) []Position {
	output := []Position{}
	for _, offset := range KingOffsets {
		newpos := Position{
			Column: pos.Column + offset.Column,
			Row:    pos.Row + offset.Row,
		}
		output = append(output, newpos)
	}
	return output
}

var BishopOffsets = []Position{
	{-1, -1} /*   */, {-1, 1},

	{1, -1} /*    */, {1, 1},
}

var RookOffsets = []Position{
	/*    */ {-1, 0}, /*    */
	{0, -1} /*    */, {0, 1},
	/*    */ {1, 0}, /*    */
}

var HorsieOffsets = []Position{
	{-2, -1}, {-2, +1},
	{-1, -2}, {-1, +2},
	{+1, -2}, {+1, +2},
	{+2, -1}, {+2, +1},
}

func genHorsieMoves(pos Position) []Position {
	output := []Position{}
	for _, offset := range HorsieOffsets {
		newpos := Position{
			Column: pos.Column + offset.Column,
			Row:    pos.Row + offset.Row,
		}
		output = append(output, newpos)
	}
	return output
}

func genRookMoves(pos Position) []Position {
	output := []Position{}
	for i := 0; i <= 7; i++ {
		if i != pos.Row {
			newPos := Position{
				Row:    i,
				Column: pos.Column,
			}
			output = append(output, newPos)
		}
		if i != pos.Column {
			newPos := Position{
				Row:    pos.Row,
				Column: i,
			}
			output = append(output, newPos)
		}
	}
	return output
}

func genBishopMoves(pos Position) []Position {
	firstDiagPos := Position{Row: 0, Column: 0}
	diff := pos.Row - pos.Column
	if diff < 0 {
		firstDiagPos = Position{Row: -diff, Column: 0}
	} else {
		firstDiagPos = Position{Row: 0, Column: diff}
	}

	secDiagPos := Position{Row: 0, Column: 7}
	if diff < 0 {
		secDiagPos = Position{Row: 7 + diff, Column: 7}
	} else {
		secDiagPos = Position{Row: 7, Column: diff}
	}

	output := []Position{}
	for i := 0; i < 7; i++ {
		firstDiag := Position{
			Row:    firstDiagPos.Row - i,
			Column: firstDiagPos.Column + i,
		}
		if firstDiag.IsValid() && firstDiag != pos {
			output = append(output, firstDiag)
		}

		secDiag := Position{
			Row:    secDiagPos.Row - i,
			Column: secDiagPos.Column - i,
		}
		if secDiag.IsValid() && secDiag != pos {
			output = append(output, secDiag)
		}
	}

	return output
}

func genQueenMoves(pos Position) []Position {
	return append(genBishopMoves(pos), genRookMoves(pos)...)
}

func genBlackPawnMoves(pos Position) []Position {
	return []Position{
		{Row: pos.Row + 2, Column: pos.Column},
		{Row: pos.Row + 1, Column: pos.Column},
		{Row: pos.Row + 1, Column: pos.Column - 1},
		{Row: pos.Row + 1, Column: pos.Column + 1},
	}
}

func genWhitePawnMoves(pos Position) []Position {
	return []Position{
		{Row: pos.Row - 2, Column: pos.Column},
		{Row: pos.Row - 1, Column: pos.Column},
		{Row: pos.Row - 1, Column: pos.Column - 1},
		{Row: pos.Row - 1, Column: pos.Column + 1},
	}
}

func orderByValue(a []Slot) {
	sort.Slice(a, func(i, j int) bool {
		isl := a[i]
		jsl := a[j]
		if isl.Piece == pc.Empty || jsl.Piece == pc.Empty {
			return false
		}
		if isl.Piece > jsl.Piece {
			return true
		}
		return false
	})
}
