package engine

import (
	eval "chess/engine/evaluation"
	movegen "chess/engine/movegen"
	"chess/game"
	"fmt"
	"math"
	"sync"
)

func BestMove(g *game.GameState) *game.Move {
	return concAlphaBeta(g, 5, eval.Evaluate)
}

func best(n *node, isMax bool) (*game.Move, float64) {
	if isMax {
		var output *game.Move
		var bestScore float64 = minusInf
		for _, leaf := range n.Leaves {
			if leaf.Score > bestScore {
				output = leaf.Move
				bestScore = leaf.Score
			}
		}
		n.Score = bestScore
		return output, bestScore
	}
	var output *game.Move
	var bestScore float64 = plusInf
	for _, leaf := range n.Leaves {
		if leaf.Score < bestScore {
			output = leaf.Move
			bestScore = leaf.Score
		}
	}
	n.Score = bestScore
	return output, bestScore
}

func close(a, b, threshold float64) bool {
	return math.Abs(a-b) <= threshold
}

var totNodes = 0

func metrics(ctx *context, start *node, best float64) {
	avg := _m(start)
	//for _, leaf := range start.Leaves {
	//	if close(leaf.Score, best, 0.5) {
	//		fmt.Printf("%v%v: %.4f\n", leaf.Move.From, leaf.Move.To, leaf.Score)
	//	}
	//}
	fmt.Printf("Average breadth: %v\n", avg)
	fmt.Printf("Total Nodes: %v\n", totNodes)
	fmt.Printf("History Size: %v\n", len(ctx.History))
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
		if leafAvg != 0 {
			avg = append(avg, leafAvg)
		}
	}
	var sum float64 = float64(currNumLeafs)
	for _, a := range avg {
		sum += a
	}
	return sum / float64(len(avg)+1)
}

func concAlphaBeta(g *game.GameState, depth int, eval Evaluator) *game.Move {
	start := &node{
		Move: &game.Move{
			State: g,
		},
		Score:  0,
		Leaves: []*node{},
	}
	c := &context{
		History: map[board]float64{},
		MaxSize: 5000,
	}

	// we set the alpha, then paralelize
	mg := movegen.NewMoveGenerator(g)

	n := mg.Next()
	if n == nil {
		return nil
	}
	leaf := newNode(n)
	start.AddLeaf(leaf)
	alpha := alphaBeta(c, leaf, depth-1, minusInf, plusInf, g.BlackTurn, eval)

	n = mg.Next()
	var wg sync.WaitGroup
	for n != nil {
		wg.Add(1)
		leaf := newNode(n)
		start.AddLeaf(leaf)
		go func() {
			defer wg.Done()
			alphaBeta(c, leaf, depth-1, alpha, plusInf, g.BlackTurn, eval)
		}()
		n = mg.Next()
	}

	wg.Wait()

	bestMove, bestScore := best(start, !g.BlackTurn)
	metrics(c, start, bestScore)

	return bestMove
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
	History  map[board]float64
	CurrSize int
	MaxSize  int
	sync.Mutex
}

const window = 0.3

func alphaBeta(c *context, n *node, depth int, alpha, beta float64, IsMax bool, eval Evaluator) float64 {
	if depth == 0 {
		n.Score = eval(n.Move.State)
		return n.Score
	}
	b := board{n.Move.State.Board, n.Move.State.BlackTurn}
	c.Mutex.Lock()
	v, ok := c.History[b]
	c.Mutex.Unlock()
	if ok {
		return v
	}
	mg := movegen.NewMoveGenerator(n.Move.State)
	if IsMax {
		maxEval := minusInf
		mv := mg.Next()
		for mv != nil {
			leaf := newNode(mv)

			score := alphaBeta(c, leaf, depth-1, alpha, beta, false, eval)
			if score > maxEval {
				n.AddLeaf(leaf)
				maxEval = score
			}
			if score > alpha {
				alpha = score
			}
			if beta+window < alpha {
				break
			}
			mv = mg.Next()
		}
		n.Score = maxEval
		c.Mutex.Lock()
		if c.CurrSize < c.MaxSize {
			c.History[b] = maxEval
			c.CurrSize++
		}
		c.Mutex.Unlock()
		return maxEval
	}
	minEval := plusInf
	mv := mg.Next()
	for mv != nil {
		leaf := newNode(mv)

		score := alphaBeta(c, leaf, depth-1, alpha, beta, true, eval)
		if score < minEval {
			n.AddLeaf(leaf)
			minEval = score
		}
		if score < beta {
			beta = score
		}
		if beta+window < alpha {
			break
		}
		mv = mg.Next()
	}
	n.Score = minEval
	c.Mutex.Lock()
	if c.CurrSize < c.MaxSize {
		c.History[b] = minEval
		c.CurrSize++
	}
	c.Mutex.Unlock()
	return minEval
}

func newNode(mv *game.Move) *node {
	return &node{Move: mv, Leaves: []*node{}}
}
