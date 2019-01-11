package main

import (
	"fmt"
	"container/ring"
)

func main() {
	players := 468
	lastMarble := 7184300

	n := ring.New(1)
	n.Value = 0
	first := n
	print(first)

	scores := make([]int, players)

	currentPlayer := 1
	for id := 1; id <= lastMarble; id++ {
		currentPlayer = (currentPlayer + 1) % players
		if id % 23 != 0 {
			n = n.Move(1)
			marble := ring.New(1)
			marble.Value = id
			n.Link(marble)
			n = marble
		} else {
			scores[currentPlayer] += id
			n = n.Move(-8)
			scores[currentPlayer] += n.Next().Value.(int)
			n.Unlink(1)
			n = n.Next()
			// print(first)
		}
	}
	max := 0
	for _,s := range scores {
		if s > max {
			max = s
		}
	}
	fmt.Println("max score", max)
}



func print(n *ring.Ring) {
	n.Do(func(x interface{}){
		fmt.Print(x," ")
	})

	fmt.Println()
}