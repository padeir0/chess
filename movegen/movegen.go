package movegen

import (
	"chess/game"
	pc "chess/game/piece"
	basicgen "chess/movegen/basic"
	seggen "chess/movegen/segregated"
	"fmt"
)

// tests if the move/unmove behaviour is working properly
// by taking multiple positions, moving until the game is over
// then unmoving back to the original position and comparing
// if it is unaltered
// breadth is the number of positions (first moves) to test,
// similar to the number of samples in a monte carlo search
func TestMoveUnmove(g *game.GameState, breadth int) string {
	startGen := seggen.NewMoveGenerator(g)
	_, ok := startGen.Next()
	for b := 0; b < breadth && ok; b++ {
		i := 0
		startPos := g.Copy()
		for !g.IsOver {
			mvgen := seggen.NewMoveGenerator(g)
			mvgen.Next()
			err := g.CheckInvalid()
			if err != "" {
				return err
			}
			i += 1
		}
		for ; i > 0; i-- {
			g.UnMove()
		}
		err := g.CheckInvalid()
		if err != "" {
			return err
		}
		err = startPos.CheckInvalid()
		if err != "" {
			return err
		}
		err = CheckGameEquals(startPos, g)
		if err != "" {
			return err
		}
		g.UnMove()
		_, ok = startGen.Next()
	}
	return ""
}

func CheckGameEquals(this *game.GameState, other *game.GameState) string {
	if this.Board != other.Board {
		return "boards don't match\n" +
			this.Board.String() + "\n" +
			other.Board.String() + "\n"
	}
	if this.BlackTurn != other.BlackTurn {
		return fmt.Sprintf("turns doesn't match: %v, %v",
			this.BlackTurn, other.BlackTurn)
	}
	if this.WhiteKingPosition != other.WhiteKingPosition {
		return fmt.Sprintf("white king position doesn't match: %v, %v",
			this.WhiteKingPosition, other.WhiteKingPosition)
	}
	if this.BlackKingPosition != other.BlackKingPosition {
		return fmt.Sprintf("black king position doesn't match: %v, %v",
			this.BlackKingPosition, other.BlackKingPosition)
	}
	if this.IsOver != other.IsOver {
		return fmt.Sprintf("isOver doesn't match: %v, %v",
			this.IsOver, other.IsOver)
	}
	if this.Result != other.Result {
		return fmt.Sprintf("result doesn't match: %v, %v",
			this.Result, other.Result)
	}
	if this.MovesSinceLastCapture != other.MovesSinceLastCapture {
		return fmt.Sprintf("MovesSinceLastCapture doesn't match: %v, %v",
			this.MovesSinceLastCapture, other.MovesSinceLastCapture)
	}
	if this.TotalValuablePieces != other.TotalValuablePieces {
		return fmt.Sprintf("TotalValuablePieces doesn't match: %v, %v",
			this.TotalValuablePieces, other.TotalValuablePieces)
	}
	// this may be wrong since pieces may be unordered
	// we may need to sort and remove invalid slots
	game.OrderSlots(this.WhitePieces)
	game.OrderSlots(other.WhitePieces)
	if len(this.WhitePieces) != len(other.WhitePieces) {
		return fmt.Sprintf("white pieces length doesn't match: %v, %v",
			len(this.WhitePieces), len(other.WhitePieces))
	}
	for i := range this.WhitePieces {
		if this.WhitePieces[i] != other.WhitePieces[i] {
			return "white piece " +
				this.WhitePieces[i].String() +
				" doesn't match with " +
				other.WhitePieces[i].String() +
				fmt.Sprintf("\n%v\n%v\n", this.WhitePieces, other.WhitePieces)
		}
	}
	game.OrderSlots(this.BlackPieces)
	game.OrderSlots(other.BlackPieces)
	if len(this.BlackPieces) != len(other.BlackPieces) {
		return fmt.Sprintf("black pieces length doesn't match: %v, %v",
			len(this.BlackPieces), len(other.BlackPieces))
	}
	for i := range this.BlackPieces {
		if this.BlackPieces[i] != other.BlackPieces[i] {
			return "black piece " +
				this.BlackPieces[i].String() +
				" doesn't match with " +
				other.BlackPieces[i].String() +
				fmt.Sprintf("\n%v\n%v\n", this.BlackPieces, other.BlackPieces)
		}
	}
	// this should be enough
	thisMove, thisOk := this.Moves.Top()
	otherMove, otherOk := other.Moves.Top()
	if thisOk != otherOk || thisMove != otherMove {
		return fmt.Sprintf("move history doesn't match: (%v, %v), (%v, %v)",
			thisOk, thisMove, otherOk, otherMove)
	}
	return ""
}

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
		fmt.Println(newG.Board.String())
		fmt.Println(newG2.Board.String())
		fmt.Println("segregated")
		fmt.Println(mvs)
		fmt.Println()
		fmt.Println("basic")
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

func move2move(mv game.Move) move {
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
