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
	Search BasicSearch
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

// IntermediateEngine prefers to do evaluation on leaf
// nodes that are quiet, but may fall short if necessary.
// Does not use any form of precomputation
type IntermediateEngine struct {
	Name     string
	Search   ExtendedSearch
	Eval     Evaluator
	Depth    int
	ExtDepth int
}

func (this *IntermediateEngine) Play(g *game.GameState) {
	bestMove := this.Search(g, this.Eval, this.ExtDepth, this.Depth)
	ok, _ := g.Move(bestMove.From, bestMove.To)
	if !ok {
		panic("engine made ilegal move")
	}
}

func (this *IntermediateEngine) String() string {
	return this.Name
}

type BasicSearch func(g *game.GameState, eval Evaluator, depth int) *game.Move
type ExtendedSearch func(g *game.GameState, eval Evaluator, extdepth, depth int) *game.Move
type Evaluator func(*game.GameState) int
