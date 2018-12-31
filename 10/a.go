package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type Point struct {
	x  int
	y  int
	dx int
	dy int
	id int
}

type Solution struct {
	mindist int64
	iteration int
}

func main() {
	points := read_input()
	solution := &Solution{mindist: 213125005000000, iteration: 0}
	for i := 0; i < 10400; i++ {
		next(points)
		if i > 10237 && i < 10241 {
			print(points, solution, i)
			fmt.Println("iteration ", i, ", first:", points[0].x, points[0].id)
		}
	}
}

func read_input() []Point {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	
	var points []Point
	id := 1
	for i:=0; i < 1000; i++{
		var p Point
		p.id = id
		id++
		_, err = fmt.Fscanf(file, "position=<%d,%d> velocity=<%d,%d>\n", &p.x, &p.y, &p.dx, &p.dy)
		if err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			panic(err)
		}

		points = append(points, p)
	}
	fmt.Printf("loaded %d points", len(points))
	return points
}

func print(points []Point, solution *Solution, iteration int) {
	sort.Slice(points, func(i, j int) bool {
		a := points[i]
		b := points[j]
		if a.x != b.x {
			return a.x < b.x
		}
		return a.y < b.y
	})


	sumx:=0
	sumy:=0
	for _,p := range points {
		sumx += p.x
		sumy += p.y
	}
	avgx:=sumx/len(points)
	avgy:=sumy/len(points)
	
	dist:=int64(0)
	for _,p := range points {
		dist += int64((avgx-p.x)*(avgx-p.x) + (avgy-p.y)*(avgy-p.y))
	}
	if(dist < solution.mindist) {
		solution.mindist = dist
		solution.iteration = iteration
	}
	fmt.Printf("\n --- avg xy=(%d,%d), distance=%d (mindist=%d for iteration %d)\n",avgx, avgy, dist, solution.mindist, solution.iteration)


	size := 500

	if points[0].x < 0 || points[0].x >= size - 1 {
		return
	}

	var out [500][500]byte
	for _,p := range points {
		if p.x > 0 && p.x < size && p.y > 0 && p.y < size {
			out[p.x][p.y] = 1
		}
	}

	var msg strings.Builder
	for j := 100; j < 150; j++ {
		for i := 170; i < 280; i++ {
			if out[i][j] == 1 {
				msg.WriteString("X")
			} else {
				msg.WriteString(".")
			}
		}
		msg.WriteString("\n")
	}
	msg.WriteString("\n")
	fmt.Println(msg.String())
}

func printp(p Point) {
	fmt.Printf("position=(%d,%d), velocity=(%d,%d) id=%d\n", p.x, p.y, p.dx, p.dy, p.id)
}

func next(points []Point) {
	for i, p := range points {
		points[i].x = p.x + p.dx
		points[i].y = p.y + p.dy
	}
}
