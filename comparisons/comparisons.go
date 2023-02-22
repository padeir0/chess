package comparisons

import (
	"chess/game"
	rs "chess/game/result"
	ifaces "chess/interfaces"

	colors "chess/asciicolors"

	"fmt"
	"runtime"
	"sync"
	"time"
)

func Compare(a, b ifaces.Engine, amount int) FightResult {
	if amount%2 != 0 {
		panic("comparison number must be even")
	}
	dwl := newDuelWorkList(a, b, amount)
	results := dwl.Start(runtime.NumCPU())
	output := FightResult{
		White: &EngineScore{
			Eng:     a,
			Score:   0,
			Average: 0,
		},
		Black: &EngineScore{
			Eng:     b,
			Score:   0,
			Average: 0,
		},
	}
	whiteTimes := []time.Duration{}
	blackTimes := []time.Duration{}
	for _, res := range results {
		if res.White.Eng == output.White.Eng {
			output.White.Score += res.White.Score
			output.Black.Score += res.Black.Score
			whiteTimes = append(whiteTimes, res.White.Average)
			blackTimes = append(blackTimes, res.Black.Average)
		} else if res.White.Eng == output.Black.Eng {
			output.White.Score += res.Black.Score
			output.Black.Score += res.White.Score
			whiteTimes = append(whiteTimes, res.Black.Average)
			blackTimes = append(blackTimes, res.White.Average)
		}
	}
	output.White.Average = average(whiteTimes)
	output.Black.Average = average(blackTimes)
	return output
}

func newDuelWorkList(a, b ifaces.Engine, amount int) *duelWorkList {
	duels := makeDuels(a, b, amount)
	return &duelWorkList{
		queue: duels,
		top:   len(duels) - 1,
		out:   make(chan FightResult),
		Mutex: sync.Mutex{},
	}
}

type duelWorkList struct {
	queue []*Duel
	top   int
	out   chan FightResult
	sync.Mutex
}

func (this *duelWorkList) Pop() *Duel {
	this.Mutex.Lock()
	if this.top < 0 {
		return nil
	}
	out := this.queue[this.top]
	this.top--
	this.Mutex.Unlock()
	return out
}

func (this *duelWorkList) Out(fr FightResult) {
	this.out <- fr
}

func (this *duelWorkList) GetResults() []FightResult {
	output := make([]FightResult, len(this.queue))
	for i := range this.queue {
		output[i] = <-this.out
	}
	return output
}

func (this *duelWorkList) Start(procs int) []FightResult {
	for i := 0; i < procs; i++ {
		go work(this)
	}
	ticker := time.NewTicker(200 * time.Millisecond)
	go this.progressBarUwU(ticker)
	results := this.GetResults()
	ticker.Stop()
	fmt.Print("\033[1A\033[K")
	return results
}

func (this *duelWorkList) progressBarUwU(ticker *time.Ticker) {
	fmt.Println()
	for {
		_, ok := <-ticker.C
		if !ok {
			break
		}
		processed := len(this.queue) - (this.top + 1)
		total := len(this.queue)
		bar := makebar(processed, total)
		fmt.Printf("\033[1A\033[K%v %v / %v                       \n", bar, processed, total)
	}
}

func makebar(processed, total int) string {
	bars := (processed * 20 / total)
	output := "|" + colors.BackgroundGreen
	for i := 0; i < bars; i++ {
		output += " "
	}
	output += colors.Reset
	for i := bars; i < 20; i++ {
		output += " "
	}
	return output + "|"
}

func work(workList *duelWorkList) {
	job := workList.Pop()
	for job != nil {
		result := job.run()
		workList.Out(result)
		job = workList.Pop()
	}
}

type Duel struct {
	White ifaces.Engine
	Black ifaces.Engine
	Board game.Board
}

func (this *Duel) run() FightResult {
	white := &EngineScore{
		Eng:   this.White,
		Score: 0,
	}
	black := &EngineScore{
		Eng:   this.Black,
		Score: 0,
	}
	whiteTimes := []time.Duration{}
	blackTimes := []time.Duration{}
	g := game.InitialGame(&this.Board)
	for !g.IsOver {
		if g.BlackTurn {
			start := time.Now()
			black.Eng.Play(g)
			blackTimes = append(blackTimes, time.Since(start))
		} else {
			start := time.Now()
			white.Eng.Play(g)
			whiteTimes = append(whiteTimes, time.Since(start))
		}
	}
	switch g.Result {
	case rs.Draw:
		white.Score += 0.5
		black.Score += 0.5
	case rs.WhiteWins:
		white.Score += 1
	case rs.BlackWins:
		black.Score += 1
	}
	white.Average = average(whiteTimes)
	black.Average = average(blackTimes)
	return FightResult{white, black}
}

type FightResult struct {
	White *EngineScore
	Black *EngineScore
}

func (this FightResult) String() string {
	return fmt.Sprintf("%v %v %0.1f x %0.1f %v %v",
		this.White.Average, this.White.Eng.String(), this.White.Score,
		this.Black.Score, this.Black.Eng.String(), this.Black.Average)
}

type EngineScore struct {
	Eng     ifaces.Engine
	Score   float64
	Average time.Duration
}

func makeDuels(A, B ifaces.Engine, number int) []*Duel {
	duels := make([]*Duel, number)
	for i := 0; i < number; i += 2 {
		board := game.ShuffledBoard()
		duels[i] = &Duel{
			White: A,
			Black: B,
			Board: *board,
		}
		duels[i+1] = &Duel{
			White: B,
			Black: A,
			Board: *board,
		}
	}
	return duels
}

func average(times []time.Duration) time.Duration {
	if len(times) == 0 {
		return 0
	}
	var sum time.Duration
	for _, t := range times {
		sum += t
	}
	return sum / time.Duration(len(times))
}
