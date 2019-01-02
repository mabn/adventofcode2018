package main

import (
	"container/list"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x  int
	y  int
	id string
}

type Location struct {
	closest  *Point
	distance int
	tie      bool
}

func main() {
	points := read()

	board := [400][400]Location{}

	// fillFrom(&points[0], &board)
	// fillFrom(&points[27], &board)
	// print(&board)

	for i, _ := range points {
		fillFrom(&points[i], &board)
		print(&board)
	}

	count := map[string]int{}
	edgesstr := map[string]int{}

	for i := 0; i <= 349; i++ {
		for j := 0; j <= 349; j++ {
			if board[i][j].closest == nil {
				continue
			}
			if i == 0 || j == 0 || i == 349 || j == 349 {
				edgesstr[board[i][j].closest.id] = 1
			}
			if !board[i][j].tie {
				count[board[i][j].closest.id]++
			}

		}
	}

	for k, _ := range edgesstr {
		fmt.Println("edgestr: ", k)
	}

	for k, _ := range count {
		if edgesstr[k] != 1 {
			fmt.Printf("%d %s\n", count[k], k)
		}
	}
	fmt.Println(count)
}

func read() []Point {

	f, _ := os.Open("input.txt")
	defer f.Close()

	id := 'A'
	points := []Point{}
	for i := 0; i < 100; i++ {
		var x, y int
		_, err := fmt.Fscanf(f, "%d, %d\n", &x, &y)
		if err != nil {
			break
		}
		points = append(points, Point{x: x, y: y, id: string(id)})
		fmt.Printf("%d,%d %s\n", x, y, string(id))
		id = id + 1

	}
	return points
}

func print(board *[400][400]Location) {
	var buf strings.Builder
	for j := 0; j < 350; j++ {
		buf.WriteString(strconv.Itoa(j) + " ")
		for i := 0; i < 350; i++ {
			if board[i][j].closest != nil && !board[i][j].tie {
				buf.WriteString(board[i][j].closest.id)
			} else if board[i][j].tie {
				buf.WriteString(".")
			} else {
				buf.WriteString(" ")
			}

		}
		buf.WriteString("\n")
	}
	buf.WriteString("\n")
	fmt.Println(buf.String())
}

func fillFrom(p *Point, board *[400][400]Location) {
	stack := list.New()
	stack.PushBack(*p)
	for stack.Len() > 0 {
		fillFromHelper(p, board, stack)
	}

	//fillFromHelper(*p.x-1, *p.y, p, board)
}

func fillFromHelper(p *Point, board *[400][400]Location, stack *list.List) {
	e := stack.Back()
	coord := Point(e.Value.(Point))
	stack.Remove(e)
	x := coord.x
	y := coord.y

	if x < 0 || x > 350 || y < 0 || y > 350 {
		return
	}

	p_distance := abs(p.x-x) + abs(p.y-y)

	if x == 3 && y == 7 {
		fmt.Printf("board: %s, p_distance: %d workiing: %s\n", board[x][y], p_distance, p)
	}

	if board[x][y].closest == nil || board[x][y].distance > p_distance {
		board[x][y].closest = p
		board[x][y].distance = p_distance
		board[x][y].tie = false
	} else if board[x][y].closest == p {
		return // already visited
	} else if board[x][y].distance == p_distance {
		// tie!
		board[x][y].closest = p
		board[x][y].tie = true

		// probably won't work and needs some more traveling

	} else {
		return
	}

	stack.PushBack(Point{x: x - 1, y: y})
	stack.PushBack(Point{x: x + 1, y: y})
	stack.PushBack(Point{x: x, y: y - 1})
	stack.PushBack(Point{x: x, y: y + 1})
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
