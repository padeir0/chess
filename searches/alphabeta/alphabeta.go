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
	mv := mg.Next()
	for mv != nil {
		leaf := &Node{Move: mv}
		alphabeta(g, leaf, alpha, beta, depth-1, eval)
		g.UnMove()

		bestMove = Max(bestMove, leaf)
		if leaf.Score > alpha {
			alpha = bestMove.Score
		}
		if beta+1 < alpha {
			break
		}
		mv = mg.Next()
	}
	n.Score = reduce(bestMove.Score)
	return bestMove
}

func minimizingPlayer(g *game.GameState, n *Node, alpha, beta int, depth int, eval ifaces.Evaluator) *Node {
	mg := movegen.NewMoveGenerator(g)
	var bestMove *Node
	mv := mg.Next()
	for mv != nil {
		leaf := &Node{Move: mv}
		alphabeta(g, leaf, alpha, beta, depth-1, eval)
		g.UnMove()

		bestMove = Min(bestMove, leaf)
		if leaf.Score < beta {
			beta = bestMove.Score
		}
		if beta+1 < alpha {
			break
		}
		mv = mg.Next()
	}
	n.Score = reduce(bestMove.Score)
	return bestMove
}

// we use this because (for example)
// a checkmate in 3 is worse than checkmate in 2
// !!!
//     it's very important to tune this on AlphaBeta
//     since it may lead to bad pruning
//     and it makes AlphaBeta perform WORSE
//     the deeper you go
// !!!
func reduce(a int) int {
	return (a * 1023) / 1024
}
