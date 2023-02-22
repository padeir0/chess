package engines

import (
	"chess/game"
	. "chess/interfaces"

	"chess/evals/custom"
	"chess/evals/material"
	"chess/evals/psqt"

	"chess/searches/alphabeta"
	"chess/searches/minimax"
	"chess/searches/negamax"
	"chess/searches/quiescence"
	"chess/searches/randcapt"
	"chess/searches/random"
	"chess/searches/typeB"
)

var AllEngines = map[string]Engine{
	"random":   Random,
	"randcapt": RandCapt,

	"minimax":         Minimax,
	"minimax_mat":     Minimax_Mat,
	"minimax_psqt":    Minimax_Psqt,
	"minimaxII":       MinimaxII,
	"minimaxII_mat":   MinimaxII_Mat,
	"minimaxII_psqt":  MinimaxII_Psqt,
	"minimaxIII_mat":  MinimaxIII_Mat,
	"minimaxIII_psqt": MinimaxIII_Psqt,

	"negamax":     NegaMax,
	"negamax_mat": NegaMax_Mat,
	"negamaxII":   NegaMaxII,

	"alphabeta":         AlphaBeta,
	"alphabeta_mat":     AlphaBeta_Mat,
	"alphabeta_psqt":    AlphaBeta_Psqt,
	"alphabetaII":       AlphaBetaII,
	"alphabetaII_mat":   AlphaBetaII_Mat,
	"alphabetaII_psqt":  AlphaBetaII_Psqt,
	"alphabetaIII":      AlphaBetaIII,
	"alphabetaIII_psqt": AlphaBetaIII_Psqt,
	"alphabetaIII_mat":  AlphaBetaIII_Mat,
	"alphabetaIV_mat":   AlphaBetaIV_Mat,
	"alphabetaIV_psqt":  AlphaBetaIV_Psqt,
	"alphabetaV_mat":    AlphaBetaV_Mat,
	"alphabetaV_psqt":   AlphaBetaV_Psqt,

	"quiescence":         Quiescence,
	"quiescence_mat":     Quiescence_Mat,
	"quiescence_psqt":    Quiescence_Psqt,
	"quiescenceII":       QuiescenceII,
	"quiescenceII_mat":   QuiescenceII_Mat,
	"quiescenceII_psqt":  QuiescenceII_Psqt,
	"quiescenceIII":      QuiescenceIII,
	"quiescenceIII_psqt": QuiescenceIII_Psqt,
	"quiescenceIII_mat":  QuiescenceIII_Mat,

	"typeb":      TypeB,
	"typeb_mat":  TypeB_Mat,
	"typeb_psqt": TypeB_Psqt,
}

var Random Engine = &BasicEngine{
	Name:   "random",
	Search: random.BestMove,
	Eval:   func(*game.GameState, int) int { return 0 },
	Depth:  0,
}

var RandCapt Engine = &BasicEngine{
	Name:   "randcapt",
	Search: randcapt.BestMove,
	Eval:   func(*game.GameState, int) int { return 0 },
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

var MinimaxIII_Mat Engine = &BasicEngine{
	Name:   "minimaxIII_mat",
	Search: minimax.BestMove,
	Eval:   material.Evaluate,
	Depth:  4,
}

var MinimaxIII_Psqt Engine = &BasicEngine{
	Name:   "minimaxIII_psqt",
	Search: minimax.BestMove,
	Eval:   psqt.Evaluate,
	Depth:  4,
}

var NegaMax Engine = &BasicEngine{
	Name:   "negamax",
	Search: negamax.BestMove,
	Eval:   custom.Evaluate,
	Depth:  2,
}

var NegaMax_Mat Engine = &BasicEngine{
	Name:   "negamax",
	Search: negamax.BestMove,
	Eval:   material.Evaluate,
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
	Search: alphabeta.BestMove,
	Eval:   custom.Evaluate,
	Depth:  2,
}

var AlphaBeta_Mat Engine = &BasicEngine{
	Name:   "alphabeta_mat",
	Search: alphabeta.BestMove,
	Eval:   material.Evaluate,
	Depth:  2,
}

var AlphaBeta_Psqt Engine = &BasicEngine{
	Name:   "alphabeta_psqt",
	Search: alphabeta.BestMove,
	Eval:   psqt.Evaluate,
	Depth:  2,
}

var AlphaBetaII Engine = &BasicEngine{
	Name:   "alphabetaII",
	Search: alphabeta.BestMove,
	Eval:   custom.Evaluate,
	Depth:  3,
}

var AlphaBetaII_Mat Engine = &BasicEngine{
	Name:   "alphabetaII_mat",
	Search: alphabeta.BestMove,
	Eval:   material.Evaluate,
	Depth:  3,
}

var AlphaBetaII_Psqt Engine = &BasicEngine{
	Name:   "alphabetaII_psqt",
	Search: alphabeta.BestMove,
	Eval:   psqt.Evaluate,
	Depth:  3,
}

var AlphaBetaIII Engine = &BasicEngine{
	Name:   "alphabetaIII",
	Search: alphabeta.BestMove,
	Eval:   custom.Evaluate,
	Depth:  4,
}

var AlphaBetaIII_Mat Engine = &BasicEngine{
	Name:   "alphabetaIII_mat",
	Search: alphabeta.BestMove,
	Eval:   material.Evaluate,
	Depth:  4,
}

var AlphaBetaIII_Psqt Engine = &BasicEngine{
	Name:   "alphabetaIII_psqt",
	Search: alphabeta.BestMove,
	Eval:   psqt.Evaluate,
	Depth:  4,
}

var AlphaBetaIV_Mat Engine = &BasicEngine{
	Name:   "alphabetaIV_mat",
	Search: alphabeta.BestMove,
	Eval:   material.Evaluate,
	Depth:  5,
}

var AlphaBetaIV_Psqt Engine = &BasicEngine{
	Name:   "alphabetaIV_psqt",
	Search: alphabeta.BestMove,
	Eval:   psqt.Evaluate,
	Depth:  5,
}

var AlphaBetaV_Mat Engine = &BasicEngine{
	Name:   "alphabetaV_mat",
	Search: alphabeta.BestMove,
	Eval:   material.Evaluate,
	Depth:  6,
}

var AlphaBetaV_Psqt Engine = &BasicEngine{
	Name:   "alphabetaV_psqt",
	Search: alphabeta.BestMove,
	Eval:   psqt.Evaluate,
	Depth:  6,
}

var Quiescence Engine = &IntermediateEngine{
	Name:     "quiescence",
	Search:   quiescence.BestMove,
	Eval:     custom.Evaluate,
	Depth:    2,
	ExtDepth: 10,
}

var Quiescence_Mat Engine = &IntermediateEngine{
	Name:     "quiescence_mat",
	Search:   quiescence.BestMove,
	Eval:     material.Evaluate,
	Depth:    2,
	ExtDepth: 10,
}

var Quiescence_Psqt Engine = &IntermediateEngine{
	Name:     "quiescence_psqt",
	Search:   quiescence.BestMove,
	Eval:     psqt.Evaluate,
	Depth:    2,
	ExtDepth: 10,
}

var QuiescenceII Engine = &IntermediateEngine{
	Name:     "quiescenceII",
	Search:   quiescence.BestMove,
	Eval:     custom.Evaluate,
	Depth:    3,
	ExtDepth: 10,
}

var QuiescenceII_Mat Engine = &IntermediateEngine{
	Name:     "quiescenceII_mat",
	Search:   quiescence.BestMove,
	Eval:     material.Evaluate,
	Depth:    3,
	ExtDepth: 10,
}

var QuiescenceII_Psqt Engine = &IntermediateEngine{
	Name:     "quiescenceII_psqt",
	Search:   quiescence.BestMove,
	Eval:     psqt.Evaluate,
	Depth:    3,
	ExtDepth: 10,
}

var QuiescenceIII Engine = &IntermediateEngine{
	Name:     "quiescenceIII",
	Search:   quiescence.BestMove,
	Eval:     custom.Evaluate,
	Depth:    4,
	ExtDepth: 10,
}

var QuiescenceIII_Mat Engine = &IntermediateEngine{
	Name:     "quiescenceIII_mat",
	Search:   quiescence.BestMove,
	Eval:     material.Evaluate,
	Depth:    4,
	ExtDepth: 10,
}

var QuiescenceIII_Psqt Engine = &IntermediateEngine{
	Name:     "quiescenceIII_psqt",
	Search:   quiescence.BestMove,
	Eval:     psqt.Evaluate,
	Depth:    4,
	ExtDepth: 10,
}

var TypeB Engine = &TypeBEngine{
	Name:    "typeb",
	Search:  typeB.BestMove,
	Eval:    custom.Evaluate,
	Depth:   5,
	Breadth: []int{5, 7, 9, 9, 15, 15},
}

var TypeB_Mat Engine = &TypeBEngine{
	Name:    "typeb_mat",
	Search:  typeB.BestMove,
	Eval:    material.Evaluate,
	Depth:   5,
	Breadth: []int{5, 7, 9, 9, 15, 15},
}

var TypeB_Psqt Engine = &TypeBEngine{
	Name:    "typeb_psqt",
	Search:  typeB.BestMove,
	Eval:    psqt.Evaluate,
	Depth:   5,
	Breadth: []int{5, 7, 9, 9, 15, 15},
}
