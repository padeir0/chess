package result

type Result int

const (
	InvalidResult Result = iota
	Draw
	WhiteWins
	BlackWins
)
