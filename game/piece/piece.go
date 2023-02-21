package piece

type Piece byte

func (this Piece) IsWhite() bool {
	switch this {
	case WhiteQueen, WhiteKing, WhiteBishop,
		WhiteRook, WhiteKnight, WhitePawn:
		return true
	}
	return false
}

func (this Piece) IsBlack() bool {
	switch this {
	case BlackQueen, BlackKing, BlackBishop,
		BlackRook, BlackKnight, BlackPawn:
		return true
	}
	return false
}

func (this Piece) IsPawnLike() bool {
	switch this {
	case BlackPawn, WhitePawn:
		return true
	}
	return false
}

func (this Piece) IsKingLike() bool {
	switch this {
	case BlackKing, WhiteKing:
		return true
	}
	return false
}

func (this Piece) IsKnightLike() bool {
	switch this {
	case BlackKnight, WhiteKnight:
		return true
	}
	return false
}

func (this Piece) IsBishopLike() bool {
	switch this {
	case BlackBishop, WhiteBishop:
		return true
	}
	return false
}

func (this Piece) IsQueenLike() bool {
	switch this {
	case BlackQueen, WhiteQueen:
		return true
	}
	return false
}

func (this Piece) IsRookLike() bool {
	switch this {
	case BlackRook, WhiteRook:
		return true
	}
	return false
}

func (this Piece) IsOccupied() bool {
	return this != Empty && this != InvalidPiece
}

func (this Piece) String() string {
	switch this {
	case InvalidPiece:
		return "?"
	case Empty:
		return " "

	case WhiteQueen:
		return "Q"
	case BlackQueen:
		return "q"
	case WhiteKing:
		return "K"
	case BlackKing:
		return "k"
	case WhiteBishop:
		return "B"
	case BlackBishop:
		return "b"
	case WhiteRook:
		return "R"
	case BlackRook:
		return "r"
	case WhiteKnight:
		return "N"
	case BlackKnight:
		return "n"
	case WhitePawn:
		return "P"
	case BlackPawn:
		return "p"
	}
	panic("should not be reached")
}

const (
	InvalidPiece Piece = iota
	Empty

	WhitePawn
	WhiteKnight
	WhiteBishop
	WhiteRook
	WhiteQueen
	WhiteKing

	BlackPawn
	BlackKnight
	BlackBishop
	BlackRook
	BlackQueen
	BlackKing
)
