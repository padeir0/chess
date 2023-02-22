package typeB

import (
	"chess/game"
	ifaces "chess/interfaces"
	movegen "chess/movegen/segregated"
	. "chess/searches/common"
	"fmt"
	"sort"
)

var _ ifaces.TypeBSearch = BestMove
var _ = fmt.Sprintf(":)")

var breadth = 5

func BestMove(g *game.GameState, eval ifaces.Evaluator, depth int, breadth []int) game.Move {
	n := &Node{
		Move:  *game.NullMove,
		Score: 314159,
	}
	newG := g.Copy()
	typeB(newG, n, depth, breadth, eval)

	if len(n.Leaves) == 0 {
		return *game.NullMove
	}

	bestMove := n.Leaves[0]

	// fmt.Println(n.NextMoves(g.BlackTurn))
	// fmt.Println("Best Move: ", bestMove.Move)
	// fmt.Println("Best Score: ", bestMove.Score)

	return bestMove.Move
}

/*
typeB (currentPosition, maximizing, depth) :
	if gameOver or depth == 0 :
        currentPosition.score = evaluate(currentPosition)
        return currentPosition.score
    expand moves
    eval positions
    if maximizing:
        sort descending
    else:
        sort ascending
    top5 := pick first 5
	for each position in top5:
		best = typeB(position, not maximizing)
	sort top5
	best = pick first of top5
	currentPosition.score = best.score
	return best.score
*/
//
func typeB(g *game.GameState, n *Node, depth int, breadth []int, eval ifaces.Evaluator) *Node {
	if depth == 0 || g.IsOver {
		n.Score = eval(g, depth)
		return n
	}
	gen := movegen.NewMoveGenerator(g)
	mv, ok := gen.Next()
	for ok {
		leaf := &Node{Move: mv}
		leaf.Score = eval(g, depth)
		n.AddLeaf(leaf)
		g.UnMove()
		mv, ok = gen.Next()
	}
	minMaxSort(g, n)
	top(n, breadth[depth])
	for _, leaf := range n.Leaves {
		ok, _ := g.Move(leaf.Move.From, leaf.Move.To)
		if !ok {
			panic("invalid move!!")
		}
		typeB(g, leaf, depth-1, breadth, eval)
		g.UnMove()
	}
	minMaxSort(g, n)
	if len(n.Leaves) == 0 {
		panic("no moves!!")
	}
	best := n.Leaves[0]
	n.Score = best.Score
	return n
}

func top(n *Node, amount int) {
	if len(n.Leaves) <= amount {
		return
	}
	n.Leaves = n.Leaves[:amount]
}

func minMaxSort(g *game.GameState, n *Node) {
	if g.BlackTurn {
		sortAscending(n)
	} else {
		sortDescending(n)
	}
}

func sortAscending(n *Node) {
	sort.SliceStable(n.Leaves, func(i, j int) bool {
		return n.Leaves[i].Score < n.Leaves[j].Score
	})
}

func sortDescending(n *Node) {
	sort.SliceStable(n.Leaves, func(i, j int) bool {
		return n.Leaves[i].Score > n.Leaves[j].Score
	})
}
