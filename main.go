package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

const (
	doors   = 3
	reveals = 1
)

var rounds = 1000

// revealed: 4 chosen: 2 winning: 1
const (
	winning option = 1 << iota
	chosen
	revealed
)

type option int

type game [doors]option

func (r game) isStat(n int, stat option) bool {
	if n < 0 || n >= 3 {
		return false
	}
	return r[n]&stat == stat
}

func (r game) String() string {
	var str string
	for x := range r {
		str += fmt.Sprint(strconv.Itoa(x) + ": ")
		str += fmt.Sprint(r.isWinning(x), r.isChosen(x), r.isRevealed(x))
		str += "\n"
	}
	return str
}

func (r game) isWinning(n int) bool {
	return r.isStat(n, winning)
}

func (r game) isChosen(n int) bool {
	return r.isStat(n, chosen)
}

func (r game) isRevealed(n int) bool {
	return r.isStat(n, revealed)
}

func (r game) isCorrect() bool {
	for k := range r {
		if r.isWinning(k) && r.isChosen(k) {
			return true
		}
	}
	return false
}

func (r *game) modifyField(n int, c option, remove ...bool) bool {
	if n < 0 || n >= 3 {
		return false
	}

	if len(remove) >= 1 && remove[0] {
		if r.isStat(n, c) {
			(*r)[n] -= c
			return true
		} else {
			return false
		}
	} else {
		if !r.isStat(n, c) {
			(*r)[n] += c
			return true
		} else {
			return false
		}
	}
}

func (r *game) switchChoice() {
	newChoices := make([]int, 0)
	for x := range *r {
		if !r.isChosen(x) && !r.isRevealed(x) {
			newChoices = append(newChoices, x)
		}
	}

	for x := range *r {
		if r.isChosen(x) {
			r.modifyField(x, chosen, true)
		}
	}

	newChoice := newChoices[rand.Intn(len(newChoices))]
	r.modifyField(newChoice, chosen)
}

func (r *game) choose(n int) {
	r.modifyField(n, chosen)
}

func (r *game) reveal() {
	for range reveals {
		ok := make([]int, 0)
		for k := range *r {
			if !r.isChosen(k) && !r.isWinning(k) {
				ok = append(ok, k)
			}
		}
		reveal := ok[rand.Intn(len(ok))]
		r.modifyField(reveal, revealed)
	}
}

func newGame() game {
	r := game{}
	r[rand.Intn(len(r))] = winning
	return r
}

func newRound(strategy bool, scores *[]bool) {
	r := newGame()
	nc := rand.Intn(len(r))
	r.choose(nc)
	r.reveal()
	if strategy {
		r.switchChoice()
	}
	*scores = append(*scores, r.isCorrect())
}

func runGantlet(strategy bool) []bool {
	scores := make([]bool, 0)
	for range rounds {
		newRound(strategy, &scores)
	}
	return scores
}

func countWins(str string, scores []bool) {
	trues := make([]bool, 0)
	for x := range scores {
		if scores[x] {
			trues = append(trues, scores[x])
		}
	}

	fmt.Printf(str, float64(len(trues))/float64(len(scores)))
}

func main() {
	countWins("no switch wins: %f\n", runGantlet(false))
	countWins("switch wins: %f\n", runGantlet(true))
}
