package main

import (
	xcmd "chess/command"
	ck "chess/command/commandkind"
	engine0 "chess/engine0"
	engine1 "chess/engine1"
	game "chess/game"
	"chess/movegen"

	rs "chess/game/result"

	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime/pprof"
	"time"
)

var asBlack = flag.Bool("black", false, "play as black")

func main() {
	flag.Parse()
	cli := newCliState()
	if !cli.ComputerIsBlack {
		enginePlay(cli)
	}
	for {
		fmt.Print(">")
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

type cliState struct {
	Saved map[string]game.GameState
	Curr  *game.GameState

	ComputerIsBlack bool
}

func newCliState() *cliState {
	return &cliState{
		Saved:           map[string]game.GameState{},
		Curr:            game.InitialGame(),
		ComputerIsBlack: !*asBlack,
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
	case ck.Move:
		if isOver(cli) {
			return
		}
		if evalMove(cli, cmd) {
			if isOver(cli) {
				return
			}
			start := time.Now()
			enginePlay(cli)
			fmt.Printf("%v\n", time.Since(start))
		}
		if isOver(cli) {
			return
		}
	case ck.Profile:
		file := *cmd.Operands[0].Label
		f, err := os.Create(file)
		if err != nil {
			warn(err)
			os.Exit(0)
		}
		pprof.StartCPUProfile(f)
	case ck.SelfPlay:
		doSelfPlay(cli)
	case ck.StopProfile:
		pprof.StopCPUProfile()
	case ck.Show:
		evalShow(cli, cmd)
	}
}

func evalMove(cli *cliState, cmd *xcmd.Command) bool {
	from := *cmd.Operands[0].Position
	to := *cmd.Operands[1].Position
	if len(cmd.Operands) == 3 {
		ok, _ := cli.Curr.Move(from, to)
		if !ok {
			warn("invalid move")
			return false
		}
		return true
	}
	ok, _ := cli.Curr.Move(from, to)
	if !ok {
		warn("invalid move")
		return false
	}
	return true
}

func evalShow(cli *cliState, cmd *xcmd.Command) {
	if len(cmd.Operands) == 0 {
		fmt.Println(cli.Curr.Board.String())
		return
	}
	whatToShow := *cmd.Operands[0].Label
	switch whatToShow {
	case "moves":
		newG := cli.Curr.Copy()
		mgen := movegen.NewMoveGenerator(newG)
		mvs := movegen.ConsumeAll(mgen)
		hls := game.MoveToHighlight(mvs)
		fmt.Println(cli.Curr.Board.Show(hls))
	default:
		warn("unimplemented")
	}
}

func enginePlay(cli *cliState) {
	mv := engine0.BestMove(cli.Curr)
	if mv == nil {
		warn("engine made no moves!!")
		return
	}
	ok, _ := cli.Curr.Move(mv.From, mv.To)
	if !ok {
		warn("engine made an illegal move!!", mv)
		return
	}
}

func isOver(cli *cliState) bool {
	if cli.Curr.IsOver {
		switch cli.Curr.Result {
		case rs.Draw:
			fmt.Println("Draw: ", cli.Curr.Reason)
		case rs.WhiteWins:
			fmt.Println("White Wins: ", cli.Curr.Reason)
		case rs.BlackWins:
			fmt.Println("Black Wins: ", cli.Curr.Reason)
		}
		return true
	}
	return false
}

func doSelfPlay(cli *cliState) {
	for !isOver(cli) {
		if cli.Curr.BlackTurn {
			fmt.Println("engine ZERO (BLACK) -------------")
			start := time.Now()
			mv := engine0.BestMove(cli.Curr)
			if mv == nil {
				warn("ZERO made no moves!!")
				return
			}
			ok, _ := cli.Curr.Move(mv.From, mv.To)
			if !ok {
				warn("ZERO made an illegal move!!")
				return
			}
			fmt.Printf("ZERO: %v\n", time.Since(start))
		} else {
			fmt.Println("engine ONE (WHITE) --------------")
			start := time.Now()
			mv := engine1.BestMove(cli.Curr)
			if mv == nil {
				warn("ONE made no moves!!")
				return
			}
			ok, _ := cli.Curr.Move(mv.From, mv.To)
			if !ok {
				warn("ONE made an illegal move!!", mv.From, mv.To)
				return
			}
			fmt.Printf("ONE: %v\n", time.Since(start))
		}
		fmt.Println(cli.Curr.Board.String())
		fmt.Println("--------------------------")
	}
}
