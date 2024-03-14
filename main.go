package main

import (
	"fmt"
	"math/rand"
	"os"
)

const (
	doors = 6
	show  = false
	file  = true
	debug = false
)

var (
	rounds           = 100
	reveals          = 1
	f       *os.File = nil
)

// revealed: 4 chosen: 2 winning: 1
const (
	winning option = 1 << iota
	chosen
	revealed
)

type option int

type game [doors]option

func (g game) isStat(n int, stat option) bool {
	if n < 0 || n >= doors {
		return false
	}
	return g[n]&stat == stat
}

func (g game) String() string {
	var str string
	for x := range g {
		str += func(x int) string {
			var str string
			if g.isWinning(x) {
				str += "w"
			} else {
				str += "-"
			}
			if g.isChosen(x) {
				str += "c"
			} else {
				str += "-"
			}
			if g.isRevealed(x) {
				str += "r"
			} else {
				str += "-"
			}
			return str
		}(x)
		if x != doors - 1 {
			str += " "
		}
	}
	return str
}

func (g game) isWinning(n int) bool {
	return g.isStat(n, winning)
}

func (g game) isChosen(n int) bool {
	return g.isStat(n, chosen)
}

func (g game) isRevealed(n int) bool {
	return g.isStat(n, revealed)
}

func (g game) isCorrect() bool {
	for k := range g {
		if g.isWinning(k) && g.isChosen(k) {
			return true
		}
	}
	return false
}

func (g *game) modifyField(n int, c option, remove ...bool) bool {
	if n < 0 || n >= doors {
		return false
	}

	if len(remove) >= 1 && remove[0] {
		if g.isStat(n, c) {
			(*g)[n] -= c
			return true
		} else {
			return false
		}
	} else {
		if !g.isStat(n, c) {
			(*g)[n] += c
			return true
		} else {
			return false
		}
	}
}

func (g *game) switchChoice() {
	newChoices := make([]int, 0)
	for x := range *g {
		if !g.isChosen(x) && !g.isRevealed(x) {
			newChoices = append(newChoices, x)
		}
	}

	for x := range *g {
		if g.isChosen(x) {
			g.modifyField(x, chosen, true)
		}
	}

	newChoice := newChoices[rand.Intn(len(newChoices))]
	g.modifyField(newChoice, chosen)
}

func (g *game) choose(n int) {
	g.modifyField(n, chosen)
}

func (g *game) reveal() {
	ok := make([]int, 0)
	for k := range *g {
		if !g.isChosen(k) && !g.isWinning(k) {
			ok = append(ok, k)
		}
	}
	reveal := ok[rand.Intn(len(ok))]
	g.modifyField(reveal, revealed)
}

func newGame() game {
	g := game{}
	g[rand.Intn(len(g))] = winning
	return g
}

func newRound(strategy bool, scores *[]bool) {
	g := newGame()
	nc := rand.Intn(len(g))
	g.choose(nc)
	for range reveals {
		g.reveal()
	}
	if strategy {
		g.switchChoice()
	}
	*scores = append(*scores, g.isCorrect())
	if show {
		fmt.Println(g)
	}
	if file && debug {
		_, err := f.WriteString(g.String() + "\n")
		if err != nil {
			panic(err)
		}
	}
}

func runGantlet(strategy bool) []bool {
	scores := make([]bool, 0)
	for range rounds {
		newRound(strategy, &scores)
	}
	return scores
}

func countWins(str string, scores []bool) float64 {
	trues := make([]bool, 0)
	for x := range scores {
		if scores[x] {
			trues = append(trues, scores[x])
		}
	}

	if show {
		fmt.Printf(str, float64(len(trues))/float64(len(scores)))
	}
	return float64(len(trues)) / float64(len(scores))
}

func main() {
	if file {
		var err error
		f, err = os.Create("stats.log")
		if err != nil {
			panic(err)
		}
	}
	var results [][2]float64
	for range doors - 2 {
		noSwitchWins := countWins("no switch wins: %f\n", runGantlet(false))
		switchWins := countWins("switch wins: %f\n", runGantlet(true))
		results = append(results, [2]float64{noSwitchWins, switchWins})
		if file {
			_, err := f.WriteString(fmt.Sprintf("no switch: %.2f%%\n", noSwitchWins * 100))
			if err != nil {
				panic(err)
			}

			_, err = f.WriteString(fmt.Sprintf("switch: %.2f%%\n", switchWins * 100))
			if err != nil {
				panic(err)
			}

			_, err = f.WriteString(fmt.Sprintf("reveals: %d, doors: %d\n", reveals, doors))
			if err != nil {
				panic(err)
			}

			_, err = f.WriteString(fmt.Sprintf("rounds: %d\n", rounds))
			if err != nil {
				panic(err)
			}

			_, err = f.WriteString(fmt.Sprintf("is it 50%% better like in original: %t\n", switchWins > noSwitchWins * 1.5))
			if err != nil {
				panic(err)
			}

			_, err = f.WriteString(fmt.Sprintf("is it better to switch: %t\n", switchWins > noSwitchWins))
			if err != nil {
				panic(err)
			}

			if reveals != doors - 1 {
				_, err = f.WriteString("\n")
				if err != nil {
					panic(err)
				}
			}

		}
		reveals += 1
	}
	if file {
		avgNoSwitch, avgSwitch := func() (float64, float64) {
			sumNo := 0.0
			sumSw := 0.0
			for _, y := range results {
				sumNo += y[0]
				sumSw += y[1]
			}
			return sumNo / float64(len(results)), sumSw / float64(len(results))
		}()
		_, err := f.WriteString(fmt.Sprintf("\nAverages:\nno switching: %.2f\nswitching: %.2f\n", avgNoSwitch, avgSwitch))
		if err != nil {
			panic(err)
		}
		if avgSwitch > avgNoSwitch {
			_, err := f.WriteString("\nit's better to switch the doors")
			if err != nil {
				panic(err)
			}
		} else {
			_, err := f.WriteString("\nit's better to not switch the doors")
			if err != nil {
				panic(err)
			}
		}
	}
}
