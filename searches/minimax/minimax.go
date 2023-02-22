package minimax

import (
	"chess/game"
	ifaces "chess/interfaces"
	movegen "chess/movegen/basic"
	. "chess/searches/common"
	"fmt"
)

var _ ifaces.BasicSearch = BestMove
var _ = fmt.Sprintf(":)")

func BestMove(g *game.GameState, eval ifaces.Evaluator, depth int) game.Move {
	n := &Node{
		Move:  *game.NullMove,
		Score: 314159,
	}
	newG := g.Copy()
	bestMove := miniMax(newG, n, depth, eval)

	//fmt.Println(n.NextMoves(g.BlackTurn))
	//fmt.Println("Best Move: ", bestMove.Move)
	//fmt.Println("Best Score: ", bestMove.Score)

	return bestMove.Move
}

func miniMax(g *game.GameState, n *Node, depth int, eval ifaces.Evaluator) *Node {
	if depth == 0 || g.IsOver {
		n.Score = eval(g, depth)
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
	mv, ok := mg.Next()
	for ok {
		leaf := &Node{Move: mv}
		miniMax(g, leaf, depth-1, eval)
		g.UnMove()
		// for debugging
		// n.AddLeaf(leaf)

		bestMove = Max(bestMove, leaf)
		mv, ok = mg.Next()
	}
	n.Score = bestMove.Score
	return bestMove
}

func minimizingPlayer(g *game.GameState, n *Node, depth int, eval ifaces.Evaluator) *Node {
	mg := movegen.NewMoveGenerator(g)
	var bestMove *Node
	mv, ok := mg.Next()
	for ok {
		leaf := &Node{Move: mv}
		miniMax(g, leaf, depth-1, eval)
		g.UnMove()
		// for debugging
		// n.AddLeaf(leaf)

		bestMove = Min(bestMove, leaf)
		mv, ok = mg.Next()
	}
	n.Score = bestMove.Score
	return bestMove
}

func reduce(a int) int {
	return (a * 1023) / 1024
}
