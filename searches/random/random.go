package random

import (
	"chess/game"
	ifaces "chess/interfaces"
	movegen "chess/movegen/basic"
	"math/rand"
)

var _ ifaces.BasicSearch = BestMove

func BestMove(g *game.GameState, eval ifaces.Evaluator, depth int) game.Move {
	newG := g.Copy()
	mvgen := movegen.NewMoveGenerator(newG)
	moves := movegen.ConsumeAll(mvgen)
	i := rand.Intn(len(moves))
	return moves[i]
}
