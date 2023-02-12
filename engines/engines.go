package engines

import (
	"chess/game"

	"chess/evals/matposmob"
	"chess/searches/minimax"
)

var Default = game.Engine{
	Search: minimax.BestMove,
	Eval:   matposmob.Evaluate,
}
