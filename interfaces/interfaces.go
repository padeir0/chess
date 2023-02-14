package interfaces

import (
	"chess/game"
)

type Engine interface {
	Play(g *game.GameState)
	String() string
}

// BasicEngine does evaluation only on leaf nodes
// and does not use any form of precomputation
type BasicEngine struct {
	Name   string
	Search Search
	Eval   Evaluator
	Depth  int
}

func (this *BasicEngine) Play(g *game.GameState) {
	bestMove := this.Search(g, this.Eval, this.Depth)
	ok, _ := g.Move(bestMove.From, bestMove.To)
	if !ok {
		panic("engine made ilegal move")
	}
}

func (this *BasicEngine) String() string {
	return this.Name
}

type Search func(g *game.GameState, eval Evaluator, depth int) *game.Move
type Evaluator func(*game.GameState) int
