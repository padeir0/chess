package main

import (
	xcmd "chess/command"
	ck "chess/command/commandkind"
	game "chess/game"
	ifaces "chess/interfaces"

	"chess/engines"

	gentest "chess/movegen"
	seggen "chess/movegen/segregated"

	rs "chess/game/result"

	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime/pprof"
	"sort"
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
		Curr:            game.InitialGame(game.ShuffledBoard()),
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
			cli.Curr = game.InitialGame(game.ShuffledBoard())
			if !gentest.CompareGens(cli.Curr, 5) {
				fmt.Println("movement generators don't match")
			}
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
	case ck.Compare:
		evalCompare(cli, cmd)
	case ck.Championship:
		evalChampionship()
	case ck.StopProfile:
		pprof.StopCPUProfile()
	case ck.Show:
		evalShow(cli, cmd)
	}
}

func evalMove(cli *cliState, cmd *xcmd.Command) bool {
	from := *cmd.Operands[0].Position
	to := *cmd.Operands[1].Position
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
		mgen := seggen.NewMoveGenerator(newG)
		mvs := seggen.ConsumeAllQuiet(mgen)
		hls := game.MoveToHighlight(mvs)
		fmt.Println(cli.Curr.Board.Show(hls))
	case "attacks":
		newG := cli.Curr.Copy()
		mgen := seggen.NewMoveGenerator(newG)
		mvs := seggen.ConsumeAllCaptures(mgen)
		hls := game.MoveToHighlight(mvs)
		fmt.Println(cli.Curr.Board.Show(hls))
	default:
		warn("unimplemented")
	}
}

func enginePlay(cli *cliState) {
	engines.QuiescenceIII.Play(cli.Curr)
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
			fmt.Println("BLACK -------------")
			start := time.Now()
			engines.MinimaxII.Play(cli.Curr)
			fmt.Printf("BLACK: %v\n", time.Since(start))
		} else {
			fmt.Println("WHITE --------------")
			start := time.Now()
			engines.MinimaxII.Play(cli.Curr)
			fmt.Printf("WHITE: %v\n", time.Since(start))
		}
		fmt.Println(cli.Curr.Board.String())
		fmt.Println("--------------------------")
	}
}

type engineScore struct {
	eng     ifaces.Engine
	score   float64
	times   []time.Duration
	average time.Duration
}

func makeBoards(number int) []*game.Board {
	boards := make([]*game.Board, number)
	for i := 0; i < number; i++ {
		boards[i] = game.ShuffledBoard()
	}
	return boards
}

func play100x(A, B ifaces.Engine) fight {
	white := &engineScore{
		eng:   A,
		score: 0,
		times: []time.Duration{},
	}
	black := &engineScore{
		eng:   B,
		score: 0,
		times: []time.Duration{},
	}
	init := time.Now()
	var numOfBoards = 50
	boards := makeBoards(numOfBoards)
	playoneside(boards, white, black)
	playoneside(boards, black, white)

	black.average = average(black.times)
	black.times = nil
	white.average = average(white.times)
	white.times = nil
	fmt.Printf("Comparison took: %v, Avg. %v: %v, Avg. %v: %v\n", time.Since(init), white.eng.String(), white.average, black.eng.String(), black.average)
	fmt.Printf("Result: %v %0.1f x %0.1f %v\n", white.eng.String(), white.score, black.score, black.eng.String())
	return fight{white, black}
}

func playoneside(boards []*game.Board, white, black *engineScore) {
	for i := 0; i < len(boards); i++ {
		g := game.InitialGame(boards[i])
		for !g.IsOver {
			if g.BlackTurn {
				start := time.Now()
				black.eng.Play(g)
				black.times = append(black.times, time.Since(start))
			} else {
				start := time.Now()
				white.eng.Play(g)
				white.times = append(white.times, time.Since(start))
			}
		}
		switch g.Result {
		case rs.Draw:
			white.score += 0.5
			black.score += 0.5
		case rs.WhiteWins:
			white.score += 1
		case rs.BlackWins:
			black.score += 1
		}
	}
}

func evalCompare(cli *cliState, cmd *xcmd.Command) {
	eng0Name := *cmd.Operands[0].Label
	eng1Name := *cmd.Operands[1].Label

	eng0, ok := engines.AllEngines[eng0Name]
	if !ok {
		warn("engine not found", eng0Name)
		return
	}
	eng1, ok := engines.AllEngines[eng1Name]
	if !ok {
		warn("engine not found: ", eng1Name)
		return
	}
	play100x(eng0, eng1)
}

func average(times []time.Duration) time.Duration {
	var sum time.Duration
	for _, t := range times {
		sum += t
	}
	return sum / time.Duration(len(times))
}

type fight struct {
	A *engineScore
	B *engineScore
}

func evalChampionship() {
	allFights := []fight{}
	for _, eng := range engines.AllEngines {
		for _, other := range engines.AllEngines {
			f := play100x(eng, other)
			allFights = append(allFights, f)
		}
	}
	sort.Slice(allFights, func(i, j int) bool {
		ires := math.Abs(allFights[i].A.score - allFights[i].B.score)
		jres := math.Abs(allFights[j].A.score - allFights[j].B.score)
		return ires < jres
	})
	fmt.Println("----------------FINAL-RESULT-----------------")
	for _, fight := range allFights {
		fmt.Printf("Result: %v %0.1f x %0.1f %v\n", fight.A.eng.String(), fight.A.score, fight.B.score, fight.B.eng.String())
	}
}
