package naivealphabeta

import (
	"chess/game"
	ifaces "chess/interfaces"
	"chess/movegen"
	. "chess/searches/common"
	"fmt"
)

var _ ifaces.Search = BestMove
var _ = fmt.Sprintf("please stop bothering me, Go")

func BestMove(g *game.GameState, eval ifaces.Evaluator, depth int) *game.Move {
	n := &Node{
		Move:  game.NullMove,
		Score: 314159,
	}
	newG := g.Copy()
	bestMove := alphabeta(newG, n, MinusInf, PlusInf, depth, eval)

	//fmt.Println(n.NextMoves(g.BlackTurn))
	//fmt.Println("Best Move: ", bestMove.Move)
	//fmt.Println("Best Score: ", bestMove.Score)

	return bestMove.Move
}

func alphabeta(g *game.GameState, n *Node, alpha, beta int, depth int, eval ifaces.Evaluator) *Node {
	if depth == 0 || g.IsOver {
		n.Score = eval(g)
		return n
	}
	if g.BlackTurn {
		return minimizingPlayer(g, n, alpha, beta, depth, eval)
	}
	return maximizingPlayer(g, n, alpha, beta, depth, eval)
}

func maximizingPlayer(g *game.GameState, n *Node, alpha, beta int, depth int, eval ifaces.Evaluator) *Node {
	mg := movegen.NewMoveGenerator(g)
	var bestMove *Node
	for {
		mv := mg.Next()
		if mv == nil {
			break
		}
		leaf := &Node{Move: mv}
		alphabeta(g, leaf, alpha, beta, depth-1, eval)
		g.UnMove()
		// for debugging
		// n.AddLeaf(leaf)

		bestMove = Max(bestMove, leaf)
		if bestMove.Score > alpha {
			alpha = bestMove.Score
		}
		if beta < alpha {
			break
		}
	}
	n.Score = reduce(bestMove.Score)
	return bestMove
}

func minimizingPlayer(g *game.GameState, n *Node, alpha, beta int, depth int, eval ifaces.Evaluator) *Node {
	mg := movegen.NewMoveGenerator(g)
	var bestMove *Node
	for {
		mv := mg.Next()
		if mv == nil {
			break
		}
		leaf := &Node{Move: mv}
		alphabeta(g, leaf, alpha, beta, depth-1, eval)
		g.UnMove()
		// for debugging
		// n.AddLeaf(leaf)

		bestMove = Min(bestMove, leaf)
		if bestMove.Score < beta {
			beta = bestMove.Score
		}
		if beta < alpha {
			break
		}
	}
	n.Score = reduce(bestMove.Score)
	return bestMove
}

// we use this because (for example)
// a checkmate in 3 is worse than checkmate in 2
func reduce(a int) int {
	return (a * 7) / 8
}
