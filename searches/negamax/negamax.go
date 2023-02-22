package negamax

import (
	"chess/game"
	ifaces "chess/interfaces"
	movegen "chess/movegen/basic"
	. "chess/searches/common"
	"fmt"
)

func BestMove(g *game.GameState, eval ifaces.Evaluator, depth int) game.Move {
	n := &Node{
		Move:  *game.NullMove,
		Score: 314159,
	}
	newG := g.Copy()
	bestScore, bestNode := negaMax(newG, n, depth, eval)

	fmt.Sprintln(n.NextMoves(g.BlackTurn))
	fmt.Sprintln("Best Move: ", bestNode.Move)
	fmt.Sprintln("Best Score: ", bestScore, bestNode.Score)

	return bestNode.Move
}

func negaMax(g *game.GameState, n *Node, depth int, eval ifaces.Evaluator) (int, *Node) {
	if depth == 0 || g.IsOver {
		n.Score = eval(g, depth)
		return player(g) * n.Score, nil
	}
	mg := movegen.NewMoveGenerator(g)

	bestScore := MinusInf
	var bestNode *Node

	mv, ok := mg.Next()
	for ok {
		leaf := &Node{Move: mv}
		score, _ := negaMax(g, leaf, depth-1, eval)
		g.UnMove()
		// for debugging
		// n.AddLeaf(leaf)

		if -score > bestScore {
			bestScore = -score
			bestNode = leaf
		}
		mv, ok = mg.Next()
	}
	n.Score = bestNode.Score
	return bestScore, bestNode
}

func player(g *game.GameState) int {
	if g.BlackTurn {
		return -1
	}
	return 1
}
