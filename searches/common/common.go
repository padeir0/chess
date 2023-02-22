package common

import (
	"chess/game"
	"fmt"
)

var MinusInf int = -(1 << 16)
var PlusInf int = (1 << 16)

func Max(a, b *Node) *Node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.Score >= b.Score {
		return a
	}
	return b
}

func Min(a, b *Node) *Node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.Score <= b.Score {
		return a
	}
	return b
}

type Node struct {
	Move game.Move

	Score int

	Leaves []*Node
}

func (this *Node) AddLeaf(n *Node) {
	if this.Leaves == nil {
		this.Leaves = make([]*Node, 5)[:0]
	}
	this.Leaves = append(this.Leaves, n)
}

func (this *Node) NextMoves(isBlack bool) string {
	output := ""
	for _, leaf := range this.Leaves {
		mv, score := leaf.Best(!isBlack)
		output += fmt.Sprintf("%v -> %v: %v\n", leaf.Move, mv, score)
	}
	return output
}

type pseudomove struct {
	from, to game.Point
}

func (this *pseudomove) String() string {
	return this.from.String() + this.to.String()
}

func (this *Node) HasDuplicates() bool {
	uniques := map[pseudomove]struct{}{}
	for _, leaf := range this.Leaves {
		psmv := pseudomove{leaf.Move.From, leaf.Move.To}
		_, ok := uniques[psmv]
		if ok {
			return true
		}
		uniques[psmv] = struct{}{}
	}
	return false
}

func (this *Node) Best(isBlackturn bool) (game.Move, int) {
	var output game.Move
	if isBlackturn {
		var bestScore int = PlusInf
		for _, leaf := range this.Leaves {
			if leaf.Score <= bestScore {
				output = leaf.Move
				bestScore = leaf.Score
			}
		}
		this.Score = bestScore
		return output, bestScore
	}
	var bestScore int = MinusInf
	for _, leaf := range this.Leaves {
		if leaf.Score >= bestScore {
			output = leaf.Move
			bestScore = leaf.Score
		}
	}
	this.Score = bestScore
	return output, bestScore
}
