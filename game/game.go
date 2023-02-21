package game

import (
	ac "chess/asciicolors"
	pc "chess/game/piece"
	rs "chess/game/result"

	"math/rand"
	"sort"
	"strconv"
)

func debugBoard(slots []*Slot) *Board {
	b := Board{}
	for i := 0; i < 64; i++ {
		b[i] = pc.Empty
	}
	b.SetPos(Point{Row: 0, Column: 4}, pc.BlackKing)
	b.SetPos(Point{Row: 7, Column: 4}, pc.WhiteKing)
	for _, slot := range slots {
		b.SetPos(slot.Pos, slot.Piece)
	}
	return &b
}

func ShuffledBoard() *Board {
	bag := []pc.Piece{pc.BlackRook, pc.BlackKnight, pc.BlackBishop, pc.BlackQueen, pc.BlackKing, pc.BlackBishop, pc.BlackKnight, pc.BlackRook}
	rand.Shuffle(len(bag), func(i, j int) {
		a := bag[i]
		bag[i] = bag[j]
		bag[j] = a
	})
	output := InitialBoard()
	for i, blackPiece := range bag {
		output.Set(0, i, blackPiece)
		whitePiece := pc.Empty
		switch blackPiece {
		case pc.BlackRook:
			whitePiece = pc.WhiteRook
		case pc.BlackKnight:
			whitePiece = pc.WhiteKnight
		case pc.BlackBishop:
			whitePiece = pc.WhiteBishop
		case pc.BlackQueen:
			whitePiece = pc.WhiteQueen
		case pc.BlackKing:
			whitePiece = pc.WhiteKing
		}
		output.Set(7, i, whitePiece)
	}
	return output
}

func InitialBoard() *Board {
	return &Board{
		pc.BlackRook, pc.BlackKnight, pc.BlackBishop, pc.BlackQueen, pc.BlackKing, pc.BlackBishop, pc.BlackKnight, pc.BlackRook,
		pc.BlackPawn, pc.BlackPawn, pc.BlackPawn, pc.BlackPawn, pc.BlackPawn, pc.BlackPawn, pc.BlackPawn, pc.BlackPawn,
		pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty,
		pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty,
		pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty,
		pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty, pc.Empty,
		pc.WhitePawn, pc.WhitePawn, pc.WhitePawn, pc.WhitePawn, pc.WhitePawn, pc.WhitePawn, pc.WhitePawn, pc.WhitePawn,
		pc.WhiteRook, pc.WhiteKnight, pc.WhiteBishop, pc.WhiteQueen, pc.WhiteKing, pc.WhiteBishop, pc.WhiteKnight, pc.WhiteRook,
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

func (this *Board) AtPos(pos Point) pc.Piece {
	return (*this)[pos.Column+8*pos.Row]
}

func (this *Board) At(row, column int) pc.Piece {
	return (*this)[column+8*row]
}

func (this *Board) SetPos(pos Point, s pc.Piece) {
	(*this)[pos.Column+8*pos.Row] = s
}

func (this *Board) Set(row, column int, s pc.Piece) {
	(*this)[column+8*row] = s
}

func (this *Board) Pop(pos Point) pc.Piece {
	i := pos.Column + 8*pos.Row
	ret := (*this)[i]
	(*this)[i] = pc.Empty
	return ret
}

type Highlight struct {
	Pos   Point
	Color ac.Color
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

type Slot struct {
	Piece pc.Piece
	Pos   Point
}

func (this Slot) IsValid() bool {
	return this.Piece != pc.Empty && this.Piece != pc.InvalidPiece
}

func (this Slot) IsInvalid() bool {
	return this.Piece == pc.Empty || this.Piece == pc.InvalidPiece
}

func (this Slot) String() string {
	return this.Piece.String() + this.Pos.String()
}

func NewPosition(s string) Point {
	col := s[0]
	row := s[1]
	if col >= 'a' && col <= 'h' &&
		row >= '1' && row <= '8' {
		pos := Point{
			Column: int(col - 'a'),
			Row:    7 - int(row-'1'),
		}
		return pos
	}
	panic("invalid position string")
}

type Point struct {
	Row    int // 1 2 3 4 5 6 7 8
	Column int // a b c d e f g h
}

func (this Point) IsValid() bool {
	return this.Column >= 0 && this.Column <= 7 &&
		this.Row >= 0 && this.Row <= 7
}

func (this Point) IsInvalid() bool {
	return this.Column < 0 || this.Column > 7 ||
		this.Row < 0 || this.Row > 7
}

func (this Point) String() string {
	col := rune(this.Column) + 'a'
	row := (7 - rune(this.Row)) + '1'
	return string(col) + string(row)
}

func InitialGame(board *Board) *GameState {
	game := &GameState{
		BlackTurn: false,
		Board:     *board,

		BlackPieces: []Slot{},
		WhitePieces: []Slot{},

		Moves: NewMoveStack(),

		IsOver: false,
		Result: rs.InvalidResult,

		TotalValuablePieces:   0,
		MovesSinceLastCapture: 0,
	}

	for i, piece := range game.Board {
		if piece != pc.Empty {
			position := Point{Column: i % 8, Row: i / 8}
			slot := Slot{Piece: piece, Pos: position}
			if piece.IsWhite() {
				game.WhitePieces = append(game.WhitePieces, slot)
			} else {
				game.BlackPieces = append(game.BlackPieces, slot)
			}
			if piece == pc.WhiteKing {
				game.WhiteKingPosition = slot.Pos
			}
			if piece == pc.BlackKing {
				game.BlackKingPosition = slot.Pos
			}
			if !piece.IsPawnLike() {
				game.TotalValuablePieces += 1
			}
		}
	}

	OrderSlots(game.WhitePieces)
	OrderSlots(game.BlackPieces)

	return game
}

func NewMoveStack() *MoveStack {
	return &MoveStack{
		top:  0,
		data: make([]Move, 64),
	}
}

type MoveStack struct {
	top  int
	data []Move
}

func (this *MoveStack) Push(mv Move) {
	if this.data == nil {
		this.data = make([]Move, 64)
	}
	if this.top >= len(this.data) {
		this.data = append(this.data, make([]Move, 64)...)
	}
	this.data[this.top] = mv
	this.top++
}

func (this *MoveStack) Pop() (Move, bool) {
	if this.top <= 0 {
		return Move{}, false
	}
	this.top--
	return this.data[this.top], true
}

func (this *MoveStack) Top() (Move, bool) {
	if this.top <= 0 {
		return Move{}, false
	}
	return this.data[this.top-1], false
}

func (this *MoveStack) Copy() *MoveStack {
	a := &MoveStack{
		top:  this.top,
		data: make([]Move, len(this.data)),
	}
	for i, mv := range this.data {
		a.data[i] = mv
	}
	return a
}

func (this *MoveStack) String() string {
	output := ""
	for i, move := range this.data[:this.top] {
		if i%2 == 0 {
			output += strconv.Itoa(i+1) + ". " + move.String()
		} else {
			output += " " + move.String() + ", "
		}
	}
	return output
}

type GameState struct {
	BlackTurn bool
	Board     Board

	// to check for Checks
	BlackKingPosition Point
	WhiteKingPosition Point

	BlackPieces []Slot
	WhitePieces []Slot

	Moves *MoveStack

	IsOver bool
	Result rs.Result
	Reason string

	TotalValuablePieces   int
	MovesSinceLastCapture int
}

// for debugging
func (this *GameState) CheckPieces() bool {
	unique := map[Point]pc.Piece{}
	for _, slot := range this.BlackPieces {
		if slot.IsInvalid() {
			continue
		}
		_, ok := unique[slot.Pos]
		if ok {
			return true
		}
		unique[slot.Pos] = slot.Piece
	}
	return false
}

func (this *GameState) Copy() *GameState {
	output := &GameState{
		BlackTurn:             this.BlackTurn,
		Board:                 this.Board,
		BlackKingPosition:     this.BlackKingPosition,
		WhiteKingPosition:     this.WhiteKingPosition,
		BlackPieces:           make([]Slot, len(this.BlackPieces)),
		WhitePieces:           make([]Slot, len(this.WhitePieces)),
		TotalValuablePieces:   this.TotalValuablePieces,
		MovesSinceLastCapture: this.MovesSinceLastCapture,
		Moves:                 this.Moves.Copy(),
		IsOver:                this.IsOver,
		Result:                this.Result,
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

var NullMove = &Move{
	Piece: pc.InvalidPiece,
	From: Point{
		Row:    0,
		Column: 0,
	},
	To: Point{
		Row:    0,
		Column: 0,
	},
	Capture: nil,
}

// returns if the move was sucessful
// passing your turn is represented by from == to
func (this *GameState) Move(from, to Point) (bool, *Slot) {
	if this.IsOver {
		return false, nil
	}
	// so we can restore this counter on UnMoves
	movesSinceLastCapt := this.MovesSinceLastCapture
	if from == to { // passing turn (null move)
		previous, ok := this.Moves.Top()
		if ok && previous.IsPass() {
			this.IsOver = true
			this.Result = rs.Draw
			this.Reason = "Both players passed turn"
		}
		null := *NullMove
		this.MovesSinceLastCapture++
		null.MovesSinceLastCapture = movesSinceLastCapt

		this.Moves.Push(null)
		this.BlackTurn = !this.BlackTurn
		return true, nil
	}
	fromPiece := this.Board.AtPos(from)
	if fromPiece.IsWhite() == this.BlackTurn {
		return false, nil
	}

	ok, capture := this.IsValidMove(from, to)
	if !ok {
		return false, nil
	}

	piece := this.Board.AtPos(from)

	if capture != nil {
		this.Board.Pop(capture.Pos)
		if capture.Piece == pc.BlackKing {
			this.IsOver = true
			this.Result = rs.WhiteWins
			this.Reason = "Black king was captured"
		} else if capture.Piece == pc.WhiteKing {
			this.IsOver = true
			this.Result = rs.BlackWins
			this.Reason = "White king was captured"
		}
		this.MovesSinceLastCapture = 0
	} else {
		this.MovesSinceLastCapture++
		if this.MovesSinceLastCapture == 50 {
			this.IsOver = true
			this.Result = rs.Draw
			this.Reason = "50 move limit exceeded"
		}
	}

	oldpiece := piece
	if canPromote(this, piece, to) {
		piece = promote(piece.IsBlack())
	}
	this.Board.SetPos(from, pc.Empty)
	this.Board.SetPos(to, piece)

	this.updatePieceTable(piece, capture, from, to)

	this.Moves.Push(Move{
		Piece:   oldpiece,
		From:    from,
		To:      to,
		Capture: capture,

		MovesSinceLastCapture: movesSinceLastCapt,
	})
	this.BlackTurn = !this.BlackTurn
	return true, capture
}

// returns if its valid and the position of the captured piece, if any
func (g *GameState) IsValidMove(from, to Point) (bool, *Slot) {
	if from.IsInvalid() || to.IsInvalid() {
		return false, nil
	}
	fromPiece := g.Board.AtPos(from)
	if fromPiece == pc.Empty {
		return false, nil
	}

	switch fromPiece {
	case pc.BlackKing, pc.WhiteKing:
		return isValidMovedKingMove(g, from, to)

	case pc.BlackKnight, pc.WhiteKnight:
		return isValidHorsieMove(g, from, to)

	case pc.BlackQueen, pc.WhiteQueen:
		return isValidQueenMove(g, from, to)

	case pc.BlackBishop, pc.WhiteBishop:
		return isValidBishopMove(g, from, to)

	case pc.BlackRook, pc.WhiteRook:
		return isValidRookMove(g, from, to)

	case pc.BlackPawn:
		return isValidBlackPawnMove(g, from, to)
	case pc.WhitePawn:
		return isValidWhitePawnMove(g, from, to)
	}
	panic("oh no!")
}

func (this *GameState) IsAttacked(pos Point, isBlack bool) bool {
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

func (this *GameState) updatePieceTable(piece pc.Piece, capture *Slot, from, to Point) {
	if piece == pc.BlackKing {
		this.BlackKingPosition = to
	}
	if piece == pc.WhiteKing {
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
			if capture.Piece == pc.BlackKing {
				this.BlackKingPosition = Point{-1, -1}
			}
			for i, slot := range this.BlackPieces {
				if slot.IsValid() && slot.Pos == capture.Pos {
					this.BlackPieces[i] = Slot{pc.Empty, Point{0, 0}}
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
			if capture.Piece == pc.WhiteKing {
				this.WhiteKingPosition = Point{-1, -1}
			}
			for i, slot := range this.WhitePieces {
				if slot.IsValid() && slot.Pos == capture.Pos {
					this.WhitePieces[i] = Slot{pc.Empty, Point{0, 0}}
					break
				}
			}
		}
		return
	}
}

func (this *GameState) unmakeTableUpdate(piece pc.Piece, capture *Slot, from, to Point) {
	if piece == pc.BlackKing {
		this.BlackKingPosition = from
	}
	if piece == pc.WhiteKing {
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
			if capture.Piece == pc.BlackKing {
				this.BlackKingPosition = capture.Pos
			}
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
			if capture.Piece == pc.WhiteKing {
				this.WhiteKingPosition = capture.Pos
			}
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

func (this *GameState) UnMove() {
	mv, ok := this.Moves.Pop()
	if !ok {
		return
	}
	this.unmakeTableUpdate(mv.Piece, mv.Capture, mv.From, mv.To)
	this.Board.Pop(mv.To)
	if mv.Capture != nil {
		this.Board.SetPos(mv.Capture.Pos, mv.Capture.Piece)
	}
	this.Board.SetPos(mv.From, mv.Piece)
	this.BlackTurn = !this.BlackTurn
	this.MovesSinceLastCapture = mv.MovesSinceLastCapture
	if this.IsOver {
		this.IsOver = false
		this.Reason = ""
		this.Result = rs.InvalidResult
	}
}

/*
 P
###
*/
//
func isValidBlackPawnMove(g *GameState, from, to Point) (bool, *Slot) {
	// check shape of movement
	if to.Row-from.Row != 1 ||
		to.Column-from.Column > 1 || to.Column-from.Column < -1 {
		return false, nil
	}
	return isValidPawnMove(g, from, to)
}

/*
###
 P
*/
//
func isValidWhitePawnMove(g *GameState, from, to Point) (bool, *Slot) {
	// check shape of movement
	if to.Row-from.Row != -1 ||
		to.Column-from.Column > 1 || to.Column-from.Column < -1 {
		return false, nil
	}
	return isValidPawnMove(g, from, to)
}

func isValidPawnMove(g *GameState, from, to Point) (bool, *Slot) {
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

func isValidQueenMove(g *GameState, from, to Point) (bool, *Slot) {
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

func isValidRookMove(g *GameState, from, to Point) (bool, *Slot) {
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

func isValidBishopMove(g *GameState, from, to Point) (bool, *Slot) {
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

func isValidKingMove(g *GameState, from, to Point) (bool, *Slot) {
	// TODO: validate castling
	return isValidMovedKingMove(g, from, to)
}

/*
### (-1, -1) (-1, 0) (-1, 1)
#K# (0,  -1)         (0,  1)
### (1,  -1) (1,  0) (1,  1)
*/
func isValidMovedKingMove(g *GameState, from, to Point) (bool, *Slot) {
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
//if you remove this comment, gofmt will destroy the formatting
func isValidHorsieMove(g *GameState, from, to Point) (bool, *Slot) {
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
func bishop_ClosestPieceInWay(g *GameState, from, to Point) *Point {
	rowQuant := 1
	if from.Row > to.Row {
		rowQuant = -1
	}
	colQuant := 1
	if from.Column > to.Column {
		colQuant = -1
	}

	currPos := Point{
		Row:    from.Row + rowQuant,
		Column: from.Column + colQuant,
	}
	destPos := Point{
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
func rook_ClosestPieceInWay(g *GameState, from, to Point) (Point, bool) {
	if from.Column != to.Column {
		quant := 1
		if from.Column > to.Column {
			quant = -1
		}
		for i := from.Column + quant; i != to.Column+quant; i += quant {
			if g.Board.At(from.Row, i) != pc.Empty {
				return Point{
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
				return Point{
					Row:    i,
					Column: from.Column,
				}, true
			}
		}
	}
	return Point{}, false
}

func canPromote(g *GameState, piece pc.Piece, to Point) bool {
	return (piece == pc.BlackPawn && to.Row == 7) ||
		(piece == pc.WhitePawn && to.Row == 0)
}

func promote(isBlack bool) pc.Piece {
	if isBlack {
		return pc.BlackQueen
	}
	return pc.WhiteQueen
}

func Abs(a int32) int32 {
	y := a >> 31
	return (a ^ y) - y
}

type Move struct {
	Piece   pc.Piece
	From    Point
	To      Point
	Capture *Slot

	MovesSinceLastCapture int
}

func (this *Move) String() string {
	if this.Capture != nil {
		return this.Piece.String() +
			this.From.String() + this.To.String() +
			"x" + this.Capture.Piece.String()
	}
	return this.Piece.String() +
		this.From.String() + this.To.String()
}

func (this *Move) IsPass() bool {
	return this.From == this.To
}

var KingOffsets = []Point{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1} /*    */, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

var BishopOffsets = []Point{
	{-1, -1} /*   */, {-1, 1},

	{1, -1} /*    */, {1, 1},
}

var RookOffsets = []Point{
	/*    */ {-1, 0}, /*    */
	{0, -1} /*    */, {0, 1},
	/*    */ {1, 0}, /*    */
}

var HorsieOffsets = []Point{
	{-2, -1}, {-2, +1},
	{-1, -2}, {-1, +2},
	{+1, -2}, {+1, +2},
	{+2, -1}, {+2, +1},
}

func OrderSlots(a []Slot) {
	sort.SliceStable(a, func(i, j int) bool {
		isl := a[i]
		jsl := a[j]
		if isl.Piece > jsl.Piece {
			return false // order by value
		}
		// make sure we have deterministic ordering
		if isl.Piece == jsl.Piece {
			if isl.Pos.Column > jsl.Pos.Column {
				return false
			}
			if isl.Pos.Column == jsl.Pos.Column &&
				isl.Pos.Row > jsl.Pos.Row {
				return false
			}
		}
		return true // should never happen, but ok
	})
}
