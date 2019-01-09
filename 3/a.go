package main

import (
	"fmt"
	"os"
)

const (
	N int = 1000
)

func main() {
	f, _ := os.Open("input.txt")
	defer f.Close()
	fabric := [N][N]int{}
	fabricclaims := [N][N]int{}
	claims := map[int]bool{}
	overlapping := map[int]bool{}
	var id, x, y, w, h int
	for {
		_, err := fmt.Fscanf(f, "#%d @ %d,%d: %dx%d\n", &id, &x, &y, &w, &h)
		if err != nil {
			break
		}
		fmt.Println(id, x, y, w, h)

		for j := y; j < y+h; j++ {
			for i := x; i < x+w; i++ {
				fabric[i][j] = fabric[i][j] + 1
				claims[id] =true
				if fabricclaims[i][j] != 0 {
					overlapping[fabricclaims[i][j]] = true
					overlapping[id] = true
				}
				fabricclaims[i][j] = id
			}
		}
	}

	cnt :=0
	for x:=0;x<N;x++{
		for y:=0; y<N; y++ {
			if fabric[x][y] > 1 {
				cnt++
			}
		}
	}

	for y:=0;y<10;y++{
		for x:=0; x<10; x++ {
			fmt.Printf("%d", fabric[x][y])
		}
		fmt.Println()
	}

	fmt.Println("ok", cnt)

	// part 2

	for k,_ := range claims {
		if !overlapping[k] {
			fmt.Println("part2", k)
			os.Exit(0)
		}
	}
	fmt.Println("part 2 answer not found")
}
