package piece

type Piece byte

func (this Piece) IsWhite() bool {
	switch this {
	case WhiteQueen, WhiteKing, WhiteBishop,
		WhiteRook, WhiteHorsie, WhitePawn:
		return true
	}
	return false
}

func (this Piece) IsBlack() bool {
	switch this {
	case BlackQueen, BlackKing, BlackBishop,
		BlackRook, BlackHorsie, BlackPawn:
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

func (this Piece) IsHorsieLike() bool {
	switch this {
	case BlackHorsie, WhiteHorsie:
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

func (this Piece) String() string {
	switch this {
	case InvalidPiece:
		return "?"
	case Empty:
		return " "

	case WhiteQueen, BlackQueen:
		return "Q"
	case WhiteKing, BlackKing:
		return "K"
	case WhiteBishop, BlackBishop:
		return "B"
	case WhiteRook, BlackRook:
		return "R"
	case WhiteHorsie, BlackHorsie:
		return "N"
	case WhitePawn, BlackPawn:
		return "P"
	}
	panic("should not be reached")
}

const (
	InvalidPiece Piece = iota
	Empty

	WhitePawn
	WhiteHorsie
	WhiteBishop
	WhiteRook
	WhiteQueen
	WhiteKing

	BlackPawn
	BlackHorsie
	BlackBishop
	BlackRook
	BlackQueen
	BlackKing
)
