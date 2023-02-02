package main

import (
	"bufio"
	xcmd "chess/command"
	ck "chess/command/commandkind"
	engine "chess/engine"
	game "chess/game"
	pr "chess/game/promotion"
	"fmt"
	"os"
	"os/exec"
)

type MoveHist struct {
	Cmd     *xcmd.Command
	Capture *game.Slot
}

type cliState struct {
	Moves []MoveHist
	Saved map[string]game.GameState
	Curr  *game.GameState
}

func newCliState() *cliState {
	return &cliState{
		Moves: []MoveHist{},
		Saved: map[string]game.GameState{},
		Curr:  game.InitialGame(),
	}
}

func main() {
	cli := newCliState()
	for {
		reader := bufio.NewReader(os.Stdin)
		cmdstr, err := reader.ReadString('\n')
		if err != nil {
			fatal(err)
		}
		cmd, err2 := xcmd.Parse(cmdstr)
		if err2 != nil {
			warn(err2)
			continue
		}
		eval(cli, cmd)
	}
}

func warn(stuff ...any) {
	fmt.Print("\u001b[31m")
	fmt.Println(stuff...)
	fmt.Print("\u001b[0m")
}

func fatal(anything ...any) {
	fmt.Println(anything...)
	os.Exit(0)
}

func eval(cli *cliState, cmd *xcmd.Command) {
	switch cmd.Kind {
	case ck.Clear:
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
	case ck.Quit:
		os.Exit(0)
	case ck.NO:
		fmt.Println("i'm sorry :(")
		os.Exit(0)
	case ck.Save:
		txt := *cmd.Operands[0].Label
		cli.Saved[txt] = *cli.Curr
	case ck.Restore:
		if len(cmd.Operands) == 0 {
			cli.Curr = game.InitialGame()
			return
		}
		txt := *cmd.Operands[0].Label
		saved, ok := cli.Saved[txt]
		if !ok {
			warn(txt + " doesn't exist")
			return
		}
		cli.Curr = &saved
	case ck.Next:
		warn("unimplemented")
	case ck.Move:
		ok, piece := evalMove(cli, cmd)
		if ok {
			cli.Moves = append(cli.Moves, MoveHist{cmd, piece})
		}
		PrintMoves(cli.Curr)
	case ck.Undo: // needs to know captured pieces
		warn("unimplemented")
		//evalUndo(cli, cmd)
	case ck.Show:
		if len(cmd.Operands) == 0 {
			fmt.Println(cli.Curr.Board.String())
			return
		}
		warn("unimplemented")
	}
}

func evalUndo(cli *cliState, cmd *xcmd.Command) {
	if len(cmd.Operands) == 1 {
		// multiple undos
	}
}

func undoOne(cl *cliState) {
	// all evaluated moves are possible
	for _, move := range cl.Moves {
		// promotion
		if len(move.Cmd.Operands) == 3 {
		}
	}
}

func evalMove(cli *cliState, cmd *xcmd.Command) (bool, *game.Slot) {
	from := *cmd.Operands[0].Position
	to := *cmd.Operands[1].Position
	if len(cmd.Operands) == 3 {
		promTxt := *cmd.Operands[2].Label
		prom := convertLabelToProm(promTxt)
		if prom == pr.InvalidPromotion {
			warn("invalid promotion")
			return false, nil
		}
		ok, capture := cli.Curr.Move(from, to, prom)
		if !ok {
			warn("invalid move")
			return false, nil
		}
		return true, capture
	}
	ok, capture := cli.Curr.Move(from, to, pr.InvalidPromotion)
	if !ok {
		warn("invalid move")
		return false, nil
	}
	return true, capture
}

func convertLabelToProm(label string) pr.Promotion {
	switch label {
	case "W", "w":
		return pr.Queen
	case "H", "h":
		return pr.Horsie
	case "B", "b":
		return pr.Bishop
	case "R", "r":
		return pr.Rook
	}
	return pr.InvalidPromotion
}

func PrintMoves(g *game.GameState) {
	moves := engine.GenerateMoves(g)
	output := []string{}
	for _, move := range moves {
		str := move.From.String() + move.To.String()
		output = append(output, str)
	}
	fmt.Println("BlackTurn: ", g.BlackTurn)
	fmt.Println(len(output), "moves:")
	fmt.Println(output)
}
