package engine

import (
	. "chess/common"
	"chess/engine0/eval"
	"chess/game"
	"chess/movegen"
	"fmt"
)

func BestMove(g *game.GameState) *game.Move {
	totNodes = 0
	ctx := &context{
		TranspositionTable: map[board]int{},
		CurrSize:           0,
		MaxSize:            5000,
	}
	n := &node{
		Move:  game.NullMove,
		Score: 314159,
	}
	newG := g.Copy()
	score := alphaBeta(ctx, newG, n, 4, minusInf, plusInf)
	mv, best := best(n, g.BlackTurn)
	metrics(ctx, n, mv, best, score)
	return mv
}

func best(n *node, isBlackturn bool) (*game.Move, int) {
	var output *game.Move
	if isBlackturn {
		var bestScore int = plusInf
		for _, leaf := range n.Leaves {
			if leaf.Score <= bestScore {
				output = leaf.Move
				bestScore = leaf.Score
			}
		}
		n.Score = bestScore
		return output, bestScore
	}
	var bestScore int = minusInf
	for _, leaf := range n.Leaves {
		if leaf.Score >= bestScore {
			output = leaf.Move
			bestScore = leaf.Score
		}
	}
	n.Score = bestScore
	return output, bestScore
}

func isClose(a, b, threshold int) bool {
	return Abs(a-b) <= threshold
}

var totNodes = 0

func metrics(ctx *context, start *node, bestMove *game.Move, best, score int) {
	avg := _m(start)
	for _, leaf := range start.Leaves {
		fmt.Printf("%v%v: %v\n", leaf.Move.From, leaf.Move.To, leaf.Score)
	}
	fmt.Println("Best Move: ", bestMove)
	fmt.Println("Best Score: ", best)
	fmt.Println("AB Score: ", score)
	fmt.Printf("Average breadth: %v\n", avg)
	fmt.Printf("Total Nodes: %v\n", totNodes)
	fmt.Printf("Transposition Table Size: %v\n", len(ctx.TranspositionTable))
}

func _m(start *node) int {
	if start == nil {
		panic(start)
	}
	totNodes++
	currNumLeafs := 0
	avg := []int{}
	for _, leaf := range start.Leaves {
		currNumLeafs++
		leafAvg := _m(leaf)
		if leafAvg != 0 {
			avg = append(avg, leafAvg)
		}
	}
	var sum int = int(currNumLeafs)
	for _, a := range avg {
		sum += a
	}
	return sum / int(len(avg)+1)
}

type node struct {
	Move *game.Move

	Score int

	Leaves []*node
}

func (this *node) AddLeaf(n *node) {
	if this.Leaves == nil {
		this.Leaves = make([]*node, 5)[:0]
	}
	this.Leaves = append(this.Leaves, n)
}

var minusInf int = -(1 << 16)
var plusInf int = (1 << 16)

type Evaluator func(g *game.GameState) int

type board struct {
	b       game.Board
	isBlack bool
}

type context struct {
	TranspositionTable map[board]int
	CurrSize           int
	MaxSize            int
}

func alphaBeta(c *context, g *game.GameState, n *node, depth int, alpha, beta int) int {
	if depth == 0 || g.IsOver {
		n.Score = eval.Evaluate(g)
		return n.Score
	}
	b := board{g.Board, g.BlackTurn}
	v, ok := c.TranspositionTable[b]
	if ok {
		n.Score = v
		return v
	}
	mg := movegen.NewMoveGenerator(g)
	if g.BlackTurn {
		minEval := plusInf
		for {
			mv := mg.Next()
			if mv == nil {
				break
			}
			leaf := &node{Move: mv}
			score := alphaBeta(c, g, leaf, depth-1, alpha, beta)
			g.UnMove()
			n.AddLeaf(leaf)

			if score < minEval {
				minEval = score
			}
			if score < alpha {
				break
			}
			if score < beta {
				beta = score
			}
		}
		n.Score = minEval
		if c.CurrSize < c.MaxSize {
			c.TranspositionTable[b] = minEval
			c.CurrSize++
		}
		return minEval
	}
	maxEval := minusInf
	for {
		mv := mg.Next()
		if mv == nil {
			break
		}
		leaf := &node{Move: mv}
		score := alphaBeta(c, g, leaf, depth-1, alpha, beta)
		g.UnMove()
		n.AddLeaf(leaf)

		if score > maxEval {
			maxEval = score
		}
		if score > beta {
			break
		}
		if score > alpha {
			alpha = score
		}
	}
	n.Score = maxEval
	if c.CurrSize < c.MaxSize {
		c.TranspositionTable[b] = maxEval
		c.CurrSize++
	}
	return maxEval
}
