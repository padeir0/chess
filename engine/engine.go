package engine

import (
	eval "chess/engine/evaluation"
	movegen "chess/engine/movegen"
	"chess/game"
	"fmt"
	"math"
)

func BestMove(g *game.GameState) *game.Move {
	start := &node{
		Move: &game.Move{
			State: g,
		},
		Score:  0,
		Leaves: []*node{},
	}
	c := &context{
		History: map[board]float64{},
	}
	res := alphaBeta(c, start, 5, minusInf, plusInf, !g.BlackTurn, eval.Evaluate)
	metrics(start, res)
	return best(start, !g.BlackTurn)
}

func best(n *node, isMax bool) *game.Move {
	if isMax {
		var output *game.Move
		var bestScore float64 = minusInf
		for _, leaf := range n.Leaves {
			if leaf.Score > bestScore {
				output = leaf.Move
				bestScore = leaf.Score
			}
		}
		return output
	}
	var output *game.Move
	var bestScore float64 = plusInf
	for _, leaf := range n.Leaves {
		if leaf.Score < bestScore {
			output = leaf.Move
			bestScore = leaf.Score
		}
	}
	return output
}

func close(a, b, threshold float64) bool {
	return math.Abs(a-b) <= threshold
}

var totNodes = 0

func metrics(start *node, best float64) {
	avg := _m(start)
	for _, leaf := range start.Leaves {
		if close(leaf.Score, best, 0.5) {
			fmt.Printf("%v%v: %.4f\n", leaf.Move.From, leaf.Move.To, leaf.Score)
		}
	}
	fmt.Printf("Average breadth: %v\n", avg)
	fmt.Printf("Total Nodes: %v\n", totNodes)
}

func _m(start *node) float64 {
	if start != nil {
		totNodes++
	}
	currNumLeafs := 0
	avg := []float64{}
	for _, leaf := range start.Leaves {
		currNumLeafs++
		leafAvg := _m(leaf)
		avg = append(avg, leafAvg)
	}
	var sum float64 = float64(currNumLeafs)
	for _, a := range avg {
		sum += a
	}
	return sum / float64(len(avg)+1)
}

type node struct {
	Move *game.Move

	Score float64

	Leaves []*node
}

func (this *node) AddLeaf(n *node) {
	this.Leaves = append(this.Leaves, n)
}

var minusInf float64 = -(1 << 16)
var plusInf float64 = (1 << 16)

type Evaluator func(g *game.GameState) float64

type board struct {
	b       game.Board
	isBlack bool
}

type context struct {
	History map[board]float64
}

func alphaBeta(c *context, n *node, depth int, alpha, beta float64, IsMax bool, eval Evaluator) float64 {
	if depth == 0 {
		n.Score = eval(n.Move.State)
		return n.Score
	}
	b := board{n.Move.State.Board, n.Move.State.BlackTurn}
	v, ok := c.History[b]
	if ok {
		return v
	}
	mg := movegen.NewMoveGenerator(n.Move.State)
	if IsMax {
		maxEval := minusInf
		mv := mg.Next()
		for mv != nil {
			leaf := newNode(mv)
			n.AddLeaf(leaf)

			score := alphaBeta(c, leaf, depth-1, alpha, beta, false, eval)
			if score > maxEval {
				maxEval = score
			}
			if score > alpha {
				alpha = score
			}
			if beta <= alpha {
				break
			}
			mv = mg.Next()
		}
		n.Score = maxEval
		c.History[b] = maxEval
		return maxEval
	}
	minEval := plusInf
	mv := mg.Next()
	for mv != nil {
		leaf := newNode(mv)
		n.AddLeaf(leaf)

		score := alphaBeta(c, leaf, depth-1, alpha, beta, true, eval)
		if score < minEval {
			minEval = score
		}
		if score < beta {
			beta = score
		}
		if beta <= alpha {
			break
		}
		mv = mg.Next()
	}
	n.Score = minEval
	c.History[b] = minEval
	return minEval
}

func newNode(mv *game.Move) *node {
	return &node{Move: mv, Leaves: []*node{}}
}
