package engines

import (
	"chess/game"

	"chess/evals/matposmob"
	"chess/searches/minimax"
	"chess/searches/random"
)

var AllEngines = map[string]game.Engine{
	"default": Default,
	"random":  Random,
}

var Default game.Engine = &game.BasicEngine{
	Name:   "default",
	Search: minimax.BestMove,
	Eval:   matposmob.Evaluate,
}

var Random game.Engine = &game.BasicEngine{
	Name:   "random",
	Search: random.BestMove,
	Eval:   func(*game.GameState) int { return 0 },
}
