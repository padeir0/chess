package minimax

import (
	"chess/game"
	"chess/movegen"
	"fmt"
)

func BestMove(g *game.GameState, eval game.Evaluator) *game.Move {
	totNodes = 0
	n := &node{
		Move:  game.NullMove,
		Score: 314159,
	}
	newG := g.Copy()
	miniMax(newG, n, 3, eval)
	mv, _ := best(n, g.BlackTurn)
	//metrics(n, mv, best, score, g.BlackTurn)
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

var totNodes = 0

func metrics(start *node, bestMove *game.Move, bestScore, score int, isBlack bool) {
	avg := _m(start)
	for _, leaf := range start.Leaves {
		mv, score := best(leaf, !isBlack)
		fmt.Printf("%v -> %v: %v\n", leaf.Move, mv, score)
	}
	fmt.Println("Best Move: ", bestMove)
	fmt.Println("Best Score: ", bestScore)
	fmt.Println("AB Score: ", score)
	fmt.Printf("Average breadth: %v\n", avg)
	fmt.Printf("Total Nodes: %v\n", totNodes)
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

func miniMax(g *game.GameState, n *node, depth int, eval game.Evaluator) int {
	if depth == 0 || g.IsOver {
		n.Score = eval(g)
		return n.Score
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
			score := (miniMax(g, leaf, depth-1, eval) * 3) / 4
			g.UnMove()
			n.AddLeaf(leaf)

			if score < minEval {
				minEval = score
			}
		}
		n.Score = minEval
		return minEval
	}
	maxEval := minusInf
	for {
		mv := mg.Next()
		if mv == nil {
			break
		}
		leaf := &node{Move: mv}
		score := (miniMax(g, leaf, depth-1, eval) * 3) / 4
		g.UnMove()
		n.AddLeaf(leaf)

		if score > maxEval {
			maxEval = score
		}
	}
	n.Score = maxEval
	return maxEval
}
