package alphabeta

import (
	"chess/game"
	ifaces "chess/interfaces"
	movegen "chess/movegen/segregated"
	. "chess/searches/common"
	"fmt"
)

var _ ifaces.BasicSearch = BestMove
var _ = fmt.Sprintf("please stop bothering me, Go")

func BestMove(g *game.GameState, eval ifaces.Evaluator, depth int) *game.Move {
	nodes = 0
	n := &Node{
		Move:  game.NullMove,
		Score: 314159,
	}
	newG := g.Copy()
	bestMove := alphabeta(newG, n, MinusInf, PlusInf, depth, eval)

	//fmt.Println("nodes: ", nodes)

	//fmt.Println(n.NextMoves(g.BlackTurn))
	//fmt.Println("Best Move: ", bestMove.Move)
	//fmt.Println("Best Score: ", bestMove.Score)

	return bestMove.Move
}

var nodes = 0

func alphabeta(g *game.GameState, n *Node, alpha, beta int, depth int, eval ifaces.Evaluator) *Node {
	nodes++
	if depth == 0 || g.IsOver {
		n.Score = eval(g, depth)
		return n
	}
	if g.BlackTurn {
		return minimizingPlayer(g, n, alpha, beta, depth, eval)
	}
	return maximizingPlayer(g, n, alpha, beta, depth, eval)
}

func maximizingPlayer(g *game.GameState, n *Node, alpha, beta int, depth int, eval ifaces.Evaluator) *Node {
	mg := movegen.NewMoveGenerator(g)
	mv := mg.Next()
	if mv == nil {
		panic("nil move!!")
	}
	var alphaMove *Node
	for mv != nil {
		leaf := &Node{Move: mv}
		alphabeta(g, leaf, alpha, beta, depth-1, eval)
		g.UnMove()

		if leaf.Score >= beta {
			n.Score = beta
			return leaf
		}
		if leaf.Score > alpha {
			alpha = leaf.Score
			alphaMove = leaf
		}
		mv = mg.Next()
	}
	n.Score = alpha
	return alphaMove
}

func minimizingPlayer(g *game.GameState, n *Node, alpha, beta int, depth int, eval ifaces.Evaluator) *Node {
	mg := movegen.NewMoveGenerator(g)
	mv := mg.Next()
	if mv == nil {
		panic("nil move!!")
	}
	var betaMove *Node
	for mv != nil {
		leaf := &Node{Move: mv}
		alphabeta(g, leaf, alpha, beta, depth-1, eval)
		g.UnMove()

		if leaf.Score <= alpha {
			n.Score = alpha
			return leaf
		}
		if leaf.Score < beta {
			beta = leaf.Score
			betaMove = leaf
		}
		mv = mg.Next()
	}
	n.Score = beta
	return betaMove
}
