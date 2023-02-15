package quiescence

import (
	"chess/game"
	ifaces "chess/interfaces"
	movegen "chess/movegen/segregated"
	. "chess/searches/common"
	"fmt"
)

var _ ifaces.ExtendedSearch = BestMove
var _ = fmt.Sprintf(":)")

func BestMove(g *game.GameState, eval ifaces.Evaluator, qdepth, depth int) *game.Move {
	n := &Node{
		Move:  game.NullMove,
		Score: 314159,
	}
	newG := g.Copy()
	bestMove := alphabeta(newG, n, MinusInf, PlusInf, qdepth, depth, eval)

	//fmt.Sprintln(n.NextMoves(g.BlackTurn))
	//fmt.Sprintln("Best Move: ", bestMove.Move)
	//fmt.Sprintln("Best Score: ", bestMove.Score)

	return bestMove.Move
}

func alphabeta(g *game.GameState, n *Node, alpha, beta, qdepth, depth int, eval ifaces.Evaluator) *Node {
	if g.IsOver {
		n.Score = eval(g)
		return n
	}
	if depth == 0 {
		quiescence(g, n, alpha, beta, qdepth, eval)
		return n
	}
	if g.BlackTurn {
		return minimizingPlayer(g, n, alpha, beta, qdepth, depth, eval)
	}
	return maximizingPlayer(g, n, alpha, beta, qdepth, depth, eval)
}

func maximizingPlayer(g *game.GameState, n *Node, alpha, beta, qdepth, depth int, eval ifaces.Evaluator) *Node {
	mg := movegen.NewMoveGenerator(g)
	var bestMove *Node
	mv := mg.Next()
	for mv != nil {
		leaf := &Node{Move: mv}
		alphabeta(g, leaf, alpha, beta, qdepth, depth-1, eval)
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

func minimizingPlayer(g *game.GameState, n *Node, alpha, beta, qdepth, depth int, eval ifaces.Evaluator) *Node {
	mg := movegen.NewMoveGenerator(g)
	var bestMove *Node
	mv := mg.Next()
	for mv != nil {
		leaf := &Node{Move: mv}
		alphabeta(g, leaf, alpha, beta, qdepth, depth-1, eval)
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

func quiescence(g *game.GameState, n *Node, alpha, beta, qdepth int, eval ifaces.Evaluator) *Node {
	if qdepth == 0 {
		n.Score = eval(g)
		return n
	}
	if g.BlackTurn {
		return quiesc_minimize(g, n, alpha, beta, qdepth, eval)
	}
	return quiesc_maximize(g, n, alpha, beta, qdepth, eval)
}

func quiesc_minimize(g *game.GameState, n *Node, alpha, beta, qdepth int, eval ifaces.Evaluator) *Node {
	mg := movegen.NewMoveGenerator(g)
	mv := mg.NextCapture()
	if mv == nil {
		n.Score = eval(g)
		return n
	}
	var bestMove *Node
	for mv != nil {
		leaf := &Node{Move: mv}
		quiescence(g, leaf, alpha, beta, qdepth-1, eval)
		g.UnMove()

		bestMove = Min(bestMove, leaf)
		if leaf.Score < beta {
			beta = bestMove.Score
		}
		if beta+1 < alpha {
			break
		}
		mv = mg.NextCapture()
	}
	n.Score = reduce(bestMove.Score)
	return bestMove
}

func quiesc_maximize(g *game.GameState, n *Node, alpha, beta, qdepth int, eval ifaces.Evaluator) *Node {
	mg := movegen.NewMoveGenerator(g)
	mv := mg.NextCapture()
	if mv == nil {
		n.Score = eval(g)
		return n
	}
	var bestMove *Node
	for mv != nil {
		leaf := &Node{Move: mv}
		quiescence(g, leaf, alpha, beta, qdepth-1, eval)
		g.UnMove()

		bestMove = Max(bestMove, leaf)
		if leaf.Score > alpha {
			alpha = bestMove.Score
		}
		if beta+1 < alpha {
			break
		}
		mv = mg.NextCapture()
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
