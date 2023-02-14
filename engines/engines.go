package engines

import (
	"chess/game"
	. "chess/interfaces"

	"chess/evals/custom"
	"chess/evals/material"
	"chess/evals/psqt"

	"chess/searches/minimax"
	"chess/searches/naivealphabeta"
	"chess/searches/negamax"
	"chess/searches/random"
)

var AllEngines = map[string]Engine{
	"random": Random,

	"minimax":        Minimax,
	"minimax_mat":    Minimax_Mat,
	"minimax_psqt":   Minimax_Psqt,
	"minimaxII":      MinimaxII,
	"minimaxII_mat":  MinimaxII_Mat,
	"minimaxII_psqt": MinimaxII_Psqt,

	"negamax":   NegaMax,
	"negamaxII": NegaMaxII,

	"alphabeta":      AlphaBeta,
	"alphabeta_mat":  AlphaBeta_Mat,
	"alphabeta_psqt": AlphaBeta_Psqt,
}

var Random Engine = &BasicEngine{
	Name:   "random",
	Search: random.BestMove,
	Eval:   func(*game.GameState) int { return 0 },
	Depth:  0,
}

var Minimax Engine = &BasicEngine{
	Name:   "minimax",
	Search: minimax.BestMove,
	Eval:   custom.Evaluate,
	Depth:  2,
}

var Minimax_Mat Engine = &BasicEngine{
	Name:   "minimax_mat",
	Search: minimax.BestMove,
	Eval:   material.Evaluate,
	Depth:  2,
}

var Minimax_Psqt Engine = &BasicEngine{
	Name:   "minimax_psqt",
	Search: minimax.BestMove,
	Eval:   psqt.Evaluate,
	Depth:  2,
}

var MinimaxII Engine = &BasicEngine{
	Name:   "minimaxII",
	Search: minimax.BestMove,
	Eval:   custom.Evaluate,
	Depth:  3,
}

var MinimaxII_Mat Engine = &BasicEngine{
	Name:   "minimaxII_mat",
	Search: minimax.BestMove,
	Eval:   material.Evaluate,
	Depth:  3,
}

var MinimaxII_Psqt Engine = &BasicEngine{
	Name:   "minimaxII_psqt",
	Search: minimax.BestMove,
	Eval:   psqt.Evaluate,
	Depth:  3,
}

var NegaMax Engine = &BasicEngine{
	Name:   "negamax",
	Search: negamax.BestMove,
	Eval:   custom.Evaluate,
	Depth:  2,
}

var NegaMaxII Engine = &BasicEngine{
	Name:   "negamaxII",
	Search: negamax.BestMove,
	Eval:   custom.Evaluate,
	Depth:  3,
}

var AlphaBeta Engine = &BasicEngine{
	Name:   "alphabeta",
	Search: naivealphabeta.BestMove,
	Eval:   custom.Evaluate,
	Depth:  3,
}

var AlphaBeta_Mat Engine = &BasicEngine{
	Name:   "alphabeta_mat",
	Search: naivealphabeta.BestMove,
	Eval:   material.Evaluate,
	Depth:  3,
}

var AlphaBeta_Psqt Engine = &BasicEngine{
	Name:   "alphabeta_psqt",
	Search: naivealphabeta.BestMove,
	Eval:   psqt.Evaluate,
	Depth:  3,
}
