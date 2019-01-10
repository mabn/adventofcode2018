package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var nextId rune

func main() {
	nextId = 'A'
	// root := load("in1.txt")
	root := load("input.txt")
	fmt.Println("total", root.totalMeta)
	fmt.Println("value", root.value)
}

func load(path string) *Node {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		panic(err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanWords)

	return readNode(s)
}

type Node struct {
	id        string
	metadata  []int
	children  []*Node
	totalMeta int
	value     int
}

func readNode(s *bufio.Scanner) *Node {
	id := string(nextId)
	nextId++
	childNum := read(s)
	metaNum := read(s)

	var children []*Node
	totalMeta := 0
	for i := 0; i < childNum; i++ {
		child := readNode(s)
		children = append(children, child)
		totalMeta += child.totalMeta
	}
	var metadata []int
	var nodeValueChildren int
	var localMeta int
	for i := 0; i < metaNum; i++ {
		m := read(s)
		metadata = append(metadata, m)
		totalMeta += m
		localMeta += m

		if m == 0 {
			panic("m == 0")
		}

		if m <= len(children) {
			nodeValueChildren += children[m-1].value
		}

	}
	var nodeValue int
	if len(children) == 0 {
		nodeValue = localMeta
	} else {
		nodeValue = nodeValueChildren
	}
	fmt.Printf("node %s value: %d\n", id, nodeValue)
	return &Node{id, metadata, children, totalMeta, nodeValue}
}

func read(s *bufio.Scanner) int {
	if !s.Scan() {
		fmt.Println("END OF INPUT")
		return -1
	}
	num, err := strconv.Atoi(s.Text())
	if err != nil {
		panic(err)
	}
	// fmt.Println(num)
	return num
}
