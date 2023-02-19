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
	nodes = 0
	qnodes = 0
	n := &Node{
		Move:  game.NullMove,
		Score: 314159,
	}
	bestMove := alphabeta(g, n, MinusInf, PlusInf, qdepth, depth, eval)

	//fmt.Println("nodes: ", nodes, "qnodes: ", qnodes)
	//fmt.Println(n.NextMoves(g.BlackTurn))
	//fmt.Println("Best Move: ", bestMove.Move)
	//fmt.Println("Best Score: ", bestMove.Score)

	return bestMove.Move
}

var nodes = 0

func alphabeta(g *game.GameState, n *Node, alpha, beta, qdepth, depth int, eval ifaces.Evaluator) *Node {
	nodes++
	if g.IsOver {
		n.Score = eval(g, depth)
		return n
	}
	if depth == 0 {
		quiescence(g, n, alpha, beta, depth, qdepth, eval)
		return n
	}
	if g.BlackTurn {
		return minimizingPlayer(g, n, alpha, beta, qdepth, depth, eval)
	}
	return maximizingPlayer(g, n, alpha, beta, qdepth, depth, eval)
}

func maximizingPlayer(g *game.GameState, n *Node, alpha, beta, qdepth, depth int, eval ifaces.Evaluator) *Node {
	mg := movegen.NewMoveGenerator(g)
	mv := mg.Next()
	if mv == nil {
		panic("nil move!!")
	}
	var alphaMove *Node
	for mv != nil {
		leaf := &Node{Move: mv}
		alphabeta(g, leaf, alpha, beta, qdepth, depth-1, eval)
		n.AddLeaf(leaf)
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

func minimizingPlayer(g *game.GameState, n *Node, alpha, beta, qdepth, depth int, eval ifaces.Evaluator) *Node {
	mg := movegen.NewMoveGenerator(g)
	mv := mg.Next()
	if mv == nil {
		panic("nil move!!")
	}
	var betaMove *Node
	for mv != nil {
		leaf := &Node{Move: mv}
		alphabeta(g, leaf, alpha, beta, qdepth, depth-1, eval)
		n.AddLeaf(leaf)
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

var qnodes = 0

// sometimes not capturing (evading) is the best move
// even more so in checks (capturing with a king and losing the game
// is a massive blunder)
// so, to take this into account we use the standing pat,
// we make the search ignore these blunders
func quiescence(g *game.GameState, n *Node, alpha, beta, depth, qdepth int, eval ifaces.Evaluator) *Node {
	qnodes++
	if qdepth == 0 || g.IsOver {
		n.Score = eval(g, depth+qdepth)
		return n
	}
	if g.BlackTurn {
		return quiesc_minimize(g, n, alpha, beta, depth, qdepth, eval)
	}
	return quiesc_maximize(g, n, alpha, beta, depth, qdepth, eval)
}

func quiesc_minimize(g *game.GameState, n *Node, alpha, beta, depth, qdepth int, eval ifaces.Evaluator) *Node {
	standPat := eval(g, depth+qdepth)
	n.Score = standPat
	if standPat < beta {
		beta = standPat
	}

	mg := movegen.NewMoveGenerator(g)
	mv := mg.NextCapture()
	if mv == nil {
		n.Score = standPat
		return n
	}
	var betaMove *Node
	for mv != nil {
		leaf := &Node{Move: mv}
		quiescence(g, leaf, alpha, beta, depth, qdepth-1, eval)
		g.UnMove()

		if leaf.Score <= alpha {
			n.Score = alpha
			return leaf
		}
		if leaf.Score < beta {
			beta = leaf.Score
			betaMove = leaf
		}
		mv = mg.NextCapture()
	}
	n.Score = beta
	return betaMove
}

func quiesc_maximize(g *game.GameState, n *Node, alpha, beta, depth, qdepth int, eval ifaces.Evaluator) *Node {
	standPat := eval(g, depth+qdepth)
	n.Score = standPat
	if standPat > alpha {
		alpha = standPat
	}

	mg := movegen.NewMoveGenerator(g)
	mv := mg.NextCapture()
	if mv == nil {
		n.Score = standPat
		return n
	}
	var alphaMove *Node
	for mv != nil {
		leaf := &Node{Move: mv}
		quiescence(g, leaf, alpha, beta, depth, qdepth-1, eval)
		g.UnMove()

		if leaf.Score >= beta {
			n.Score = beta
			return leaf
		}
		if leaf.Score > alpha {
			alpha = leaf.Score
			alphaMove = leaf
		}
		mv = mg.NextCapture()
	}
	n.Score = alpha
	return alphaMove
}
