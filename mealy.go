package hangulmealy

import (
	"fmt"
)

type PiFunc func(*HangulText, rune)

type Mealy struct {
	Q      []string
	Sigma  []rune
	Pi     map[string]PiFunc
	D      [][]int
	Lambda [][]string
	q0     int

	q    int
	text *HangulText
}

func NewMealy(Q []string, Sigma []rune, Pi map[string]PiFunc, D [][]int, Lambda [][]string, q0 int) *Mealy {
	return &Mealy{Q, Sigma, Pi, D, Lambda, q0, q0, &HangulText{}}
}

func (mealy *Mealy) String() string {
	return fmt.Sprintf("Q: %v\nΣ: %v\nΠ: %v\nδ: %v\nλ: %v\nq0: %v",
		mealy.Q, mealy.Sigma, mealy.Pi, mealy.D, mealy.Lambda, mealy.q0)
}

func (mealy *Mealy) RunByRune(r rune) error {
	if r == 127 || r == 8 { // backspace/delete 처리
		if len(mealy.text.String()) > 0 {
			mealy.text.RemoveLast()
			ins := mealy.text.ins

			mealy.q = mealy.q0
			mealy.text.ClearBuf()
			for _, r := range ins {
				mealy.RunByRune(r)
			}
		}
		return nil
	}

	findPi := func(q int, r rune) string {
		ir := idxOfRune(mealy.Sigma, r)
		if ir == -1 {
			ir = idxOfRune(mealy.Sigma, '*')
		}

		pi := mealy.Lambda[q][ir]
		if pi == "" { // * - Else Lambda
			pi = mealy.Lambda[q][idxOfRune(mealy.Sigma, '*')]
		}
		return pi
	}

	ir := idxOfRune(mealy.Sigma, r)
	if ir == -1 || mealy.D[mealy.q][ir] == -1 {
		if mealy.q == mealy.q0 {
			ir = idxOfRune(mealy.Sigma, '*')
		} else {
			mealy.Pi[findPi(mealy.q, r)](mealy.text, r)

			if *verboseFlag {
				fmt.Println("c:", string(r), "pi:", findPi(mealy.q, r), "q:", mealy.Q[mealy.q0], "ins:", mealy.text.ins)
			}

			mealy.q = mealy.q0
			return mealy.RunByRune(r)
		}
	}

	mealy.Pi[findPi(mealy.q, r)](mealy.text, r)

	if *verboseFlag {
		fmt.Println("c:", string(r), "pi:", findPi(mealy.q, r), "q:", mealy.Q[mealy.q], "ins:", mealy.text.ins)
	}

	mealy.q = mealy.D[mealy.q][ir]

	return nil
}

func (mealy *Mealy) Run(s string) error {
	for _, c := range s {
		e := mealy.RunByRune(c)
		if e != nil {
			return e
		}
	}

	return nil
}

func (mealy *Mealy) HangulString() string {
	return mealy.text.String()
}
