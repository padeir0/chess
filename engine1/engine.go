package engine

import (
	. "chess/common"
	"chess/engine1/eval"
	movegen "chess/engine1/movegen"
	"chess/game"
	"fmt"
	"sync"
)

func BestMove(g *game.GameState) *game.Move {
	return concAlphaBeta(g, 6)
}

func best(n *node, isMax bool) (*game.Move, int) {
	var output *game.Move
	if isMax {
		var bestScore int = minusInf
		for _, leaf := range n.Leaves {
			if leaf.Score > bestScore {
				output = leaf.Move
				bestScore = leaf.Score
			}
		}
		n.Score = bestScore
		return output, bestScore
	}
	var bestScore int = plusInf
	for _, leaf := range n.Leaves {
		if leaf.Score < bestScore {
			output = leaf.Move
			bestScore = leaf.Score
		}
	}
	n.Score = bestScore
	return output, bestScore
}

func close(a, b, threshold int) bool {
	return Abs(a-b) <= threshold
}

var totNodes = 0

func metrics(ctx *context, start *node, best int) {
	avg := _m(start)
	for _, leaf := range start.Leaves {
		if close(leaf.Score, best, 50) {
			fmt.Printf("%v%v: %v\n", leaf.Move.From, leaf.Move.To, leaf.Score)
		}
	}
	fmt.Println("Best Score: ", best)
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

func concAlphaBeta(original *game.GameState, depth int) *game.Move {
	start := &node{
		Score:  0,
		Leaves: []*node{},
	}
	c := &context{
		TranspositionTable: map[board]int{},
		MaxSize:            5000,
	}

	g := original.Copy()
	mg := movegen.NewMoveGenerator(g)

	mv := mg.Next()
	if mv == nil {
		return nil
	}
	leaf := &node{Move: mv}
	start.AddLeaf(leaf)

	// we set the alpha, then paralelize
	alpha := alphaBeta(c, g, leaf, depth-1, minusInf, plusInf, g.BlackTurn)
	g.UnmakeMove(mv)

	var wg sync.WaitGroup
	for {
		mv = mg.Next()
		if mv == nil {
			break
		}
		wg.Add(1)
		newG := g.Copy()
		g.UnmakeMove(mv)
		leaf := &node{Move: mv}
		start.AddLeaf(leaf)
		go func() {
			defer wg.Done()
			alphaBeta(c, newG, leaf, depth-1, alpha, plusInf, g.BlackTurn)
		}()
	}

	wg.Wait()

	bestMove, bestScore := best(start, !g.BlackTurn)
	metrics(c, start, bestScore)

	return bestMove
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
	sync.Mutex
}

var window = []int{0, 5, 10, 20, 30, 40, 50}

func alphaBeta(c *context, g *game.GameState, n *node, depth int, alpha, beta int, IsMax bool) int {
	if depth == 0 {
		n.Score = eval.Evaluate(g)
		return n.Score
	}
	b := board{g.Board, g.BlackTurn}
	c.Mutex.Lock()
	v, ok := c.TranspositionTable[b]
	c.Mutex.Unlock()
	if ok {
		return v
	}
	mg := movegen.NewMoveGenerator(g)
	if IsMax {
		maxEval := minusInf
		for {
			mv := mg.Next()
			if mv == nil {
				break
			}
			leaf := &node{Move: mv}
			score := alphaBeta(c, g, leaf, depth-1, alpha, beta, false)
			g.UnmakeMove(mv)

			if score > maxEval {
				n.AddLeaf(leaf)
				maxEval = score
			}
			if score > alpha {
				alpha = score
			}
			if beta+window[depth] < alpha {
				break
			}
		}
		n.Score = maxEval
		c.Mutex.Lock()
		if c.CurrSize < c.MaxSize {
			c.TranspositionTable[b] = maxEval
			c.CurrSize++
		}
		c.Mutex.Unlock()
		return maxEval
	}
	minEval := plusInf
	for {
		mv := mg.Next()
		if mv == nil {
			break
		}
		leaf := &node{Move: mv}
		score := alphaBeta(c, g, leaf, depth-1, alpha, beta, true)
		g.UnmakeMove(mv)

		if score < minEval {
			n.AddLeaf(leaf)
			minEval = score
		}
		if score < beta {
			beta = score
		}
		if beta+window[depth] < alpha {
			break
		}
	}
	n.Score = minEval
	c.Mutex.Lock()
	if c.CurrSize < c.MaxSize {
		c.TranspositionTable[b] = minEval
		c.CurrSize++
	}
	c.Mutex.Unlock()
	return minEval
}
