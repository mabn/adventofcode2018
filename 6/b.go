package main

import (
	"fmt"
	"os"
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

	board := [400][400]int{}

	result := 0
	for i := 0; i <= 349; i++ {
		for j := 0; j <= 349; j++ {
			for k, _ := range points {
				p := &points[k]
				board[i][j] += abs(p.x-i) + abs(p.y-j)
			}
			if board[i][j] < 10000 {
				result++
			}
		}
	}
	fmt.Println(result)
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
