package piece

type Piece byte

func (this Piece) IsWhite() bool {
	switch this {
	case WhiteQueen, WhiteKing, WhiteCastleKing, WhiteBishop,
		WhiteRook, WhiteMovedRook, WhiteHorsie, WhitePawn,
		WhiteMovedPawn:
		return true
	}
	return false
}

func (this Piece) IsBlack() bool {
	switch this {
	case BlackQueen, BlackKing, BlackCastleKing, BlackBishop,
		BlackRook, BlackMovedRook, BlackHorsie, BlackPawn,
		BlackMovedPawn:
		return true
	}
	return false
}

func (this Piece) IsPawnLike() bool {
	switch this {
	case BlackPawn, WhitePawn,
		BlackMovedPawn, WhiteMovedPawn:
		return true
	}
	return false
}

func (this Piece) IsKingLike() bool {
	switch this {
	case BlackKing, WhiteKing,
		BlackCastleKing, WhiteCastleKing:
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
	case BlackRook, WhiteRook,
		BlackMovedRook, WhiteMovedRook:
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
		return "W"
	case WhiteCastleKing, BlackCastleKing:
		return "K"
	case WhiteKing, BlackKing:
		return "Ḱ"
	case WhiteBishop, BlackBishop:
		return "B"
	case WhiteRook, BlackRook:
		return "R"
	case WhiteMovedRook, BlackMovedRook:
		return "Ŕ"
	case WhiteHorsie, BlackHorsie:
		return "H"
	case WhitePawn, BlackPawn:
		return "P"
	case WhiteMovedPawn, BlackMovedPawn:
		return "Ṕ"
	}
	panic("should not be reached")
}

const (
	InvalidPiece Piece = iota
	Empty

	WhiteMovedPawn
	WhitePawn
	WhiteHorsie
	WhiteBishop
	WhiteRook
	WhiteMovedRook
	WhiteQueen
	WhiteKing
	WhiteCastleKing

	BlackMovedPawn
	BlackPawn
	BlackHorsie
	BlackBishop
	BlackRook
	BlackMovedRook
	BlackQueen
	BlackKing
	BlackCastleKing
)
