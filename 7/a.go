package main

import (
	"fmt"
	"os"
)

const (
	BaseCost int = 60
	NumWorkers int = 5
)

type Rule struct {
	from int
	to   int
}

type Step struct {
	id   int
	done bool
	assigned bool
	reqs []*Step
}

type Worker struct {
	id int
	step    *Step
	elapsed int
}

func main() {
	steps := load("input.txt")
	var workers []*Worker
	for i := 0; i < NumWorkers; i++ {
		workers = append(workers, &Worker{i, nil, 0})
	}

	time :=0
	for ; !allDone(steps) && time < 30000; time++ {
		fmt.Println("time=", time)
		for _, w := range workers {
			w.maybeFinish()
		}

		for _, w := range workers {
			if !w.working() {
				w.assignWork(steps)
			}
		}

		for _, w := range workers {
			w.elapsed++
		}
	}
	fmt.Println("elapsed time: ", time-1)
}

func (w *Worker) maybeFinish() {
	ready := w.step != nil && w.elapsed >= BaseCost+w.step.id+1
	if !ready {
		return
	}
	fmt.Printf("%d finished working %s\n", w.id, string(w.step.id + 'A'))
	w.step.done = true
	w.step.assigned = false
	w.step = nil
	w.elapsed = 0
}

func (w *Worker) working() bool {
	return w.step != nil
}

func (w *Worker) assignWork(steps map[int]*Step) {
	for id := 0; id < len(steps); id++ {
		s := steps[id]
		if s.done || s.assigned{
			continue
		}
		if s.ready() {
			w.step = s
			w.elapsed = 0
			s.assigned = true
			fmt.Printf("assigned %s to w%d\n", string(s.id + 'A'), w.id)
			break
		}
	}
}

func part1() {
	steps := load("input.txt")

	for !allDone(steps) {
		for id := 0; id < len(steps); id++ {
			s := steps[id]
			if s.done {
				continue
			}
			if s.ready() {
				s.done = true
				fmt.Print(string(s.id + 'A'))
				break
			}
		}
	}
	fmt.Println()
}

func load(path string) map[int]*Step {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	rules := []Rule{}
	steps := map[int]*Step{}
	for {
		var from, to rune
		_, err := fmt.Fscanf(f, "Step %c must be finished before step %c can begin.\n", &from, &to)
		if err != nil {
			break
		}
		r := Rule{int(from - 'A'), int(to - 'A')}
		rules = append(rules, r)
		if steps[r.from] == nil {
			steps[r.from] = &Step{r.from, false, false, nil}
		}
		if steps[r.to] == nil {
			steps[r.to] = &Step{r.to, false, false, nil}
		}

	}
	fmt.Println("rules", rules, "steps count", len(steps))

	for _, r := range rules {
		to := steps[r.to]
		to.reqs = append(to.reqs, steps[r.from])
	}

	// for i:=0; i < len(steps); i++ {
	// 	s:= steps[i]
	// 	fmt.Printf("%d %s requirements: \n", s.id, string(s.id+'A'))
	// 	for _,req := range s.reqs {
	// 		fmt.Printf("  %d %s\n", req.id, string(req.id+'A'))
	// 	}
	// }
	return steps
}

func (s *Step) ready() bool {
	for _, req := range s.reqs {
		if !req.done {
			return false
		}

	}
	return true
}

func allDone(steps map[int]*Step) bool {
	for _, s := range steps {
		if !s.done {
			return false
		}
	}
	return true
}
