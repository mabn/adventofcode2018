package main

import (
	"fmt"
)



func main() {
	
	//fill(&grid, 5177)
	fmt.Println("122,79 @57", grid(57)[122][79])
	fmt.Println("217,196 @39", grid(39)[217][196])

	g := grid(5177)
	
	max := -1
	maxx := -1
	maxy := -1
	maxsquare := 1

	for square := 1; square <= 300; square++ {
		for x := 1; x <= 300-square+1; x++ {
			for y := 1; y <= 300-square+1; y++ {
				s := sum(x,y,square,&g)
				if s > max {
					max = s
					maxx = x
					maxy = y
					maxsquare = square
				}
			}
		}
		fmt.Printf("square %d, current max: %d at %d,%d,%d\n", square, max, maxx, maxy, maxsquare)
	}

	fmt.Printf("max %d at (%d,%d) for square %d\n", max, maxx, maxy, maxsquare)
}

func grid(serial int) [301][301]int {
	g := [301][301]int{}
	for x := 1; x <= 300; x++ {
		for y := 1; y <= 300; y++ {
			level := ((x + 10) * y + serial) * (x + 10)
			level = (level / 100) % 10
			g[x][y] = level - 5
		}
	}
	return g
} 

func sum(x int, y int, square int, grid *[301][301]int) int {
	sum := 0
	for i := 0; i < square; i++ {
		for j := 0; j < square; j++ {
			sum += grid[x+i][y+j]
		}
	}
	return sum
}