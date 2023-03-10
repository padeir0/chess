package main

import (
	colors "chess/asciicolors"
	xcmd "chess/command"
	ck "chess/command/commandkind"
	comps "chess/comparisons"
	game "chess/game"
	ifaces "chess/interfaces"

	"chess/engines"

	movegenTest "chess/movegen"
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
		Curr:            game.InitialGame(game.InitialBoard()),
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
			cli.Curr = game.InitialGame(game.InitialBoard())
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
	case ck.Test:
		test()
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
	case "attacked":
		showAttacked(cli)
	case "defended":
		showDefended(cli)
	default:
		warn("unimplemented")
	}
}

func enginePlay(cli *cliState) {
	engines.TypeB_Psqt.Play(cli.Curr)
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
			engines.QuiescenceIII.Play(cli.Curr)
			fmt.Printf("BLACK: %v\n", time.Since(start))
		} else {
			fmt.Println("WHITE --------------")
			start := time.Now()
			engines.QuiescenceIII.Play(cli.Curr)
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
	start := time.Now()
	res := comps.Compare(eng0, eng1, 200)
	fmt.Println("final: ", res)
	fmt.Println("comparison took: ", time.Since(start))
}

type duel struct {
	A ifaces.Engine
	B ifaces.Engine
}

var duels = []duel{
	//{engines.TypeB_Mat, engines.Random},
	//{engines.TypeB_Mat, engines.RandCapt},
	//{engines.TypeB_Mat, engines.AlphaBetaII_Mat},
	//{engines.TypeB, engines.AlphaBetaII},
	//{engines.TypeB_Psqt, engines.AlphaBetaII_Psqt},

	//{engines.Random, engines.RandCapt},
	//{engines.Minimax, engines.Random},
	//{engines.Minimax, engines.Minimax_Mat},
	//{engines.Minimax, engines.Minimax_Psqt},
	//{engines.Minimax, engines.RandCapt},
	//{engines.Minimax, engines.MinimaxII},

	{engines.AlphaBeta_Mat, engines.RandCapt},
	{engines.AlphaBeta_Psqt, engines.RandCapt},
	{engines.AlphaBeta, engines.RandCapt},

	{engines.AlphaBeta_Mat, engines.AlphaBetaII_Mat},
	{engines.AlphaBeta_Mat, engines.AlphaBetaIII_Mat},
	{engines.AlphaBetaII_Mat, engines.AlphaBetaIII_Mat},
	{engines.AlphaBetaIII_Mat, engines.AlphaBetaIV_Mat},
	{engines.AlphaBetaIV_Mat, engines.AlphaBetaV_Mat},

	{engines.AlphaBeta_Psqt, engines.AlphaBetaII_Psqt},
	{engines.AlphaBeta_Psqt, engines.AlphaBetaIII_Psqt},
	{engines.AlphaBetaII_Psqt, engines.AlphaBetaIII_Psqt},
	{engines.AlphaBetaIII_Psqt, engines.AlphaBetaIV_Psqt},
	{engines.AlphaBetaIV_Psqt, engines.AlphaBetaV_Psqt},

	//{engines.Quiescence, engines.AlphaBeta},
	//{engines.Quiescence_Mat, engines.AlphaBeta_Mat},
	//{engines.Quiescence_Psqt, engines.AlphaBeta_Psqt},

	//{engines.QuiescenceII, engines.AlphaBetaII},
	//{engines.QuiescenceII_Mat, engines.AlphaBetaII_Mat},
	//{engines.QuiescenceII_Psqt, engines.AlphaBetaII_Psqt},

	//{engines.QuiescenceIII, engines.AlphaBetaIII},
	//{engines.QuiescenceIII_Mat, engines.AlphaBetaIII_Mat},
	//{engines.QuiescenceIII_Psqt, engines.AlphaBetaIII_Psqt},
}

func evalChampionship() {
	allFights := []comps.FightResult{}
	for _, duel := range duels {
		start := time.Now()
		res := comps.Compare(duel.A, duel.B, 200)
		allFights = append(allFights, res)
		fmt.Println(res, " : ", time.Since(start))
	}
	sort.Slice(allFights, func(i, j int) bool {
		ires := math.Abs(allFights[i].White.Score - allFights[i].Black.Score)
		jres := math.Abs(allFights[j].White.Score - allFights[j].Black.Score)
		return ires < jres
	})
	fmt.Println("----------------FINAL-RESULT-----------------")
	for _, fight := range allFights {
		fmt.Println(fight)
	}
}

func test() {
	g := game.InitialGame(game.InitialBoard())
	err := movegenTest.TestMoveUnmove(g, 5)
	if err != "" {
		fmt.Printf("TestMoveUnmove %vfailed%v: %v\n", colors.Red, colors.Reset, err)
	}
	for i := 0; i < 10; i++ {
		g = game.InitialGame(game.ShuffledBoard())
		err = movegenTest.TestMoveUnmove(g, 5)
		if err != "" {
			fmt.Printf("TestMoveUnmove %vfailed%v: %v\n", colors.Red, colors.Reset, err)
		}
	}

	g = game.InitialGame(game.InitialBoard())
	if !movegenTest.CompareGens(g, 5) {
		fmt.Println("CompareGens failed")
	}

	for i := 0; i < 10; i++ {
		g = game.InitialGame(game.ShuffledBoard())
		if !movegenTest.CompareGens(g, 5) {
			fmt.Println("CompareGens failed")
		}
	}
}

func showAttacked(cli *cliState) {
	pieces := cli.Curr.WhitePieces
	if cli.Curr.BlackTurn {
		pieces = cli.Curr.BlackPieces
	}
	hls := []game.Highlight{}
	for _, slot := range pieces {
		if slot.IsInvalid() {
			continue
		}
		if cli.Curr.IsAttacked(slot.Pos, cli.Curr.BlackTurn) {
			newhl := game.Highlight{
				Pos:   slot.Pos,
				Color: colors.BackgroundMagenta,
			}
			hls = append(hls, newhl)
		}
	}
	fmt.Println(cli.Curr.Board.Show(hls))
}

func showDefended(cli *cliState) {
	pieces := cli.Curr.WhitePieces
	if cli.Curr.BlackTurn {
		pieces = cli.Curr.BlackPieces
	}
	hls := []game.Highlight{}
	for _, slot := range pieces {
		if slot.IsInvalid() {
			continue
		}
		if cli.Curr.IsAttacked(slot.Pos, !cli.Curr.BlackTurn) {
			newhl := game.Highlight{
				Pos:   slot.Pos,
				Color: colors.BackgroundMagenta,
			}
			hls = append(hls, newhl)
		}
	}
	fmt.Println(cli.Curr.Board.Show(hls))
}
