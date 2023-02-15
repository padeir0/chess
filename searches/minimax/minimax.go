package minimax

import (
	"chess/game"
	ifaces "chess/interfaces"
	movegen "chess/movegen/basic"
	. "chess/searches/common"
	// "fmt"
)

var _ ifaces.Search = BestMove

func BestMove(g *game.GameState, eval ifaces.Evaluator, depth int) *game.Move {
	n := &Node{
		Move:  game.NullMove,
		Score: 314159,
	}
	newG := g.Copy()
	bestMove := miniMax(newG, n, depth, eval)

	//fmt.Sprintln(n.NextMoves(g.BlackTurn))
	//fmt.Sprintln("Best Move: ", bestMove.Move)
	//fmt.Sprintln("Best Score: ", bestMove.Score)

	return bestMove.Move
}

func miniMax(g *game.GameState, n *Node, depth int, eval ifaces.Evaluator) *Node {
	if depth == 0 || g.IsOver {
		n.Score = eval(g)
		return n
	}
	if g.BlackTurn {
		return minimizingPlayer(g, n, depth, eval)
	}
	return maximizingPlayer(g, n, depth, eval)
}

func maximizingPlayer(g *game.GameState, n *Node, depth int, eval ifaces.Evaluator) *Node {
	mg := movegen.NewMoveGenerator(g)
	var bestMove *Node
	for {
		mv := mg.Next()
		if mv == nil {
			break
		}
		leaf := &Node{Move: mv}
		miniMax(g, leaf, depth-1, eval)
		g.UnMove()
		// for debugging
		// n.AddLeaf(leaf)

		bestMove = Max(bestMove, leaf)
	}
	n.Score = reduce(bestMove.Score)
	return bestMove
}

func minimizingPlayer(g *game.GameState, n *Node, depth int, eval ifaces.Evaluator) *Node {
	mg := movegen.NewMoveGenerator(g)
	var bestMove *Node
	for {
		mv := mg.Next()
		if mv == nil {
			break
		}
		leaf := &Node{Move: mv}
		miniMax(g, leaf, depth-1, eval)
		g.UnMove()
		// for debugging
		// n.AddLeaf(leaf)

		bestMove = Min(bestMove, leaf)
	}
	n.Score = reduce(bestMove.Score)
	return bestMove
}

// we use this because (for example)
// a checkmate in 3 is worse than checkmate in 2
func reduce(a int) int {
	return (a * 7) / 8
}
