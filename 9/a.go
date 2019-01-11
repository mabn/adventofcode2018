package main

import (
	"fmt"
)

type Node struct {
	id int
	next *Node
	prev *Node
}

func main() {
	players := 468
	lastMarble := 7184300


	n := &Node{0, nil, nil}
	n.next = n
	n.prev = n
	first := n
	first.print()

	scores := make([]int, players)

	currentPlayer := 1
	for id := 1; id <= lastMarble; id++ {
		currentPlayer = (currentPlayer + 1) % players
		if id % 23 != 0 {
			n = n.right(1)
			n = n.addMarble(id)
		} else {
			scores[currentPlayer] += id
			n = n.left(7)
			scores[currentPlayer] += n.id
			n = n.remove()
		}
		// fmt.Println("CURRENT:",id)
		// fmt.Println(first)
		// fmt.Println(first.next)
		// fmt.Println(first.next.next)
		// fmt.Println()
		// fmt.Println(first.prev)
		// fmt.Println(first.prev.prev)
		// fmt.Println(first.prev.prev.prev)
		// first.print()
	}
	max := 0
	for _,s := range scores {
		if s > max {
			max = s
		}
	}
	fmt.Println("max score", max)
}

func (n *Node) right(shift int) *Node {
	for i:=0;i<shift;i++{
		n = n.next
	}
	return n
}

func (n *Node) left(shift int) *Node {
	for i:=0;i<shift;i++{
		n = n.prev
	}
	return n
}

func (n *Node) addMarble(id int) *Node {
	right := n.next
	m := &Node{id, right, n}
	n.next = m
	right.prev = m
	return m
}

// returns node do the right
func (n *Node) remove() *Node {
	right := n.next
	left := n.prev
	left.next = right
	right.prev = left
	return right
}

func (n *Node) print() {
	fmt.Print(n.id," ") 
	for i := n.next; i.id != n.id; i = i.next {
		fmt.Print(i.id," ") 
	}
	fmt.Println()
}