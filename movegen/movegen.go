package movegen

import (
	"chess/game"
	pc "chess/game/piece"
	basicgen "chess/movegen/basic"
	seggen "chess/movegen/segregated"
	"fmt"
)

// makes sure Segregated == Basic
func CompareGens(g *game.GameState, depth int) bool {
	if depth == 0 {
		return true
	}
	newG := g.Copy()
	mgen := seggen.NewMoveGenerator(newG)
	mvs := seggen.ConsumeAll(mgen)

	newG2 := g.Copy()
	mgen2 := basicgen.NewMoveGenerator(newG2)
	mvs2 := basicgen.ConsumeAll(mgen2)

	if len(mvs) != len(mvs2) {
		fmt.Println(mvs)
		fmt.Println(mvs2)
		return false
	}
	moveSet := map[move]struct{}{}
	moveSet2 := map[move]struct{}{}
	for i := range mvs {
		moveSet[move2move(mvs[i])] = struct{}{}
		moveSet2[move2move(mvs2[i])] = struct{}{}
	}
	for i := range mvs {
		_, ok := moveSet[move2move(mvs2[i])]
		if !ok {
			fmt.Println("move not in set", mvs2[i])
			fmt.Println("segregated")
			fmt.Println(mvs)
			fmt.Println()
			fmt.Println("basic")
			fmt.Println(mvs2)
			return false
		}
		_, ok = moveSet2[move2move(mvs[i])]
		if !ok {
			fmt.Println("move not in set", mvs[i])
			fmt.Println("segregated")
			fmt.Println(mvs)
			fmt.Println()
			fmt.Println("basic")
			fmt.Println(mvs2)
			return false
		}
	}
	for _, mv := range mvs {
		g.Move(mv.From, mv.To)
		ok := CompareGens(g, depth-1)
		if !ok {
			return false
		}
		g.UnMove()
	}
	return true
}

// value based so we can use in maps
type move struct {
	from, to game.Point
	piece    pc.Piece
	captured pc.Piece
	mslc     int
}

func move2move(mv *game.Move) move {
	if mv.Capture != nil {
		return move{
			from:     mv.From,
			to:       mv.To,
			piece:    mv.Piece,
			captured: mv.Capture.Piece,
			mslc:     mv.MovesSinceLastCapture,
		}
	}
	return move{
		from:  mv.From,
		to:    mv.To,
		piece: mv.Piece,
		mslc:  mv.MovesSinceLastCapture,
	}
}
