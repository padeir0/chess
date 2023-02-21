package result

type Result int

func (this Result) String() string {
	switch this {
	case Draw:
		return "Draw"
	case WhiteWins:
		return "White Wins"
	case BlackWins:
		return "Black Wins"
	}
	return "?"
}

const (
	InvalidResult Result = iota
	Draw
	WhiteWins
	BlackWins
)
