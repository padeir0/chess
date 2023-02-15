package randcapt

import (
	"chess/game"
	ifaces "chess/interfaces"
	movegen "chess/movegen/segregated"
	"math/rand"
)

var _ ifaces.Search = BestMove

func BestMove(g *game.GameState, eval ifaces.Evaluator, depth int) *game.Move {
	newG := g.Copy()
	mvgen := movegen.NewMoveGenerator(newG)
	captures := movegen.ConsumeAllCaptures(mvgen)
	if len(captures) > 0 {
		i := rand.Intn(len(captures))
		return captures[i]
	}
	quiets := movegen.ConsumeAllQuiet(mvgen)
	if len(quiets) > 0 {
		i := rand.Intn(len(quiets))
		return quiets[i]
	}
	return game.NullMove
}
