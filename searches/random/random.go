package random

import (
	"chess/game"
	"chess/movegen"
	"math/rand"
)

func BestMove(g *game.GameState, eval game.Evaluator) *game.Move {
	newG := g.Copy()
	mvgen := movegen.NewMoveGenerator(newG)
	moves := movegen.ConsumeAll(mvgen)
	i := rand.Intn(len(moves))
	return moves[i]
}
