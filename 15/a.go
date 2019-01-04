package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Race byte

type Point struct {
	x,y int
}
const (
	Goblin Race = 0
	Elf    Race = 1
)

type Board [][]byte

type Unit struct {
	x, y int
	hp   int
	race Race
}

type Units struct {
	// goblins []*Unit
	list []*Unit
	pos  [][]*Unit
}

var board Board
var units Units

func main() {
	load("input.txt")

	board.print()

	r := 1
	for ; units.endOfCombat() && r < 1000000; r++ { // rounds
		// sort units by y.x
		
		round()
		units.sort()
		board.print()
		fmt.Printf("==> Round %d finished\n", r)
		units.removeDead()
	}

	sum := 0
	for _,u:=range units.list{
		sum += u.hp
	}
	fmt.Printf("Outcome: %d * %d = %d\n", r-1, sum, (r-1)*sum)
	fmt.Printf("Outcome: %d * %d = %d\n", r-2, sum, (r-2)*sum)
}

func load(path string) {
	board = Board{}
	units = Units{}

	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	lines := [][]byte{}
	for s.Scan() {
		line := []byte(s.Text())
		lines = append(lines, line)
	}
	height := len(lines)
	width := len(lines[0])

	board = make([][]byte, width)
	for x := range board {
		board[x] = make([]byte, height)
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			board[x][y] = lines[y][x]
		}
	}

	units.pos = make([][]*Unit, len(board))
	for x, _ := range board {
		units.pos[x] = make([]*Unit, len(board[x]))
		for y, _ := range board[x] {
			if board[x][y] == 'G' {
				board[x][y] = '.'
				units.add(Unit{x, y, 200, Goblin})
			} else if board[x][y] == 'E' {
				board[x][y] = '.'
				units.add(Unit{x, y, 200, Elf})
			}
		}
	}

}

func (units *Units) add(unit Unit) {
	units.pos[unit.x][unit.y] = &unit
	units.list = append(units.list, &unit)
}

func (units *Units) at(x int, y int) *Unit {
	return units.pos[x][y]
}

func (board Board) print() {
	for y := 0; y < len(board[0]); y++ {
		for x := 0; x < len(board); x++ {
			u := units.at(x, y)
			if u == nil {
				fmt.Print(string(board[x][y]))
			} else if u.race == Elf {
				fmt.Print("E")
			} else if u.race == Goblin {
				fmt.Print("G")
			}
		}
		fmt.Println()
	}
	for _,u := range units.list {
		fmt.Printf("%+v\n", u)
	}
}

// create distance map to each unit of specified race
func distmap(fromRace Race) [][]int {
	width := len(board)
	height := len(board[0])
	d := make([][]int, height)
	for i := range d {
		d[i] = make([]int, width)
	}

	for _, u := range units.list {
		if u.race == fromRace {
			distmapFill(1, u.x, u.y, d)
		}
	}
	return d
}

func distmapFill(dist int, x int, y int, d [][]int) {
	width := len(board)
	height := len(board[0])
	if board[x][y] == '#' {
		return
	}
	if x < 0 || x >= width || y < 0 || y >= height {
		return
	}
	if d[x][y] > 0 && d[x][y] <= dist {
		return
	}
	// space occuppied by a unit, cannot pass through here
	if units.at(x, y) != nil && dist > 1 {
		d[x][y] = 100000
		return
	} else {
		d[x][y] = dist
	}

	distmapFill(dist+1, x+1, y, d)
	distmapFill(dist+1, x-1, y, d)
	distmapFill(dist+1, x, y-1, d)
	distmapFill(dist+1, x, y+1, d)
}

func printDist(d [][]int) {
	for y := 0; y < len(d[0]); y++ {
		for x := range d {
			if d[x][y] == 100000 {
				fmt.Print(".")
			} else {
				fmt.Print(d[x][y])
			}
		}
		fmt.Println()
	}
}

func (race Race) opposite() Race {
	if race == Elf {
		return Goblin
	}
	return Elf
}

func (units *Units) sort() {
	sort.Slice(units.list, func(i, j int) bool {
		a := units.list[i]
		b := units.list[j]
		if a.y != b.y {
			return a.y < b.y
		}
		return a.x < b.x
	})
}

func round() {
	units.removeDead()
	units.sort()

	for _, unit := range units.list {
		if unit.hp <= 0 {
			continue
		}

		// fmt.Printf("processing unit %d,%d\n", unit.x, unit.y)

		// d contains distances to the nearest enemy (plus 1)
		d := distmap(unit.race.opposite())

		// printDist(d)

		to := findMoveTarget(d, unit.x, unit.y)

		
		if to != nil {
			units.move(unit, to.x-unit.x, to.y-unit.y)
		}

		target := findAttackTarget(unit)
		if target != nil {
				unit.attack(target)
		}
	}
}

func (u *Unit) attack(target *Unit) {
	fmt.Printf("%d,%d attacking %+v\n", u.x, u.y, target)
	target.hp -= 3
}

func (units *Units) move(unit *Unit, dx int, dy int) {
	// fmt.Printf("%d,%d moving by %d,%d\n", unit.x, unit.y, dx, dy)
	fmt.Printf("%d,%d moving to %d,%d\n", unit.x, unit.y, unit.x+dx, unit.y+dy)
	units.pos[unit.x][unit.y] = nil
	unit.x += dx
	unit.y += dy
	units.pos[unit.x][unit.y] = unit
}

func (units *Units) removeDead() {
	newlist := []*Unit{}
	for _,u := range units.list {
		if u.hp > 0 {
			newlist = append(newlist, u)
		}else{
			fmt.Printf("Removing dead unit %+v\n", u)
			units.pos[u.x][u.y] = nil
		}

	}
	units.list = newlist
}

func (units *Units) endOfCombat() bool {
	counts := []int{0,0}

	for _,u := range units.list {
		counts[u.race]++
	}
	return counts[0] > 0 && counts[1] > 0
}

func findMoveTarget(d [][]int, x int, y int) *Point {
	adjacent := []Point{
			{x, y-1},
			{x -1 , y},
			{x + 1, y},
			{x , y+1}}

	min := 10000
	var to *Point

	for i,_ := range adjacent {
		a := adjacent[i]
		val := d[a.x][a.y]
		// fmt.Printf(" (%d,%d) -> %2d %2d val=%2d\n",x,y,a.x, a.y, val)

		// adjacent to enemy, don't move
		if val == 1 {
			return nil
		}
		if val > 0 && val < min {
			to = &a
			min = val
			// fmt.Printf("   new min = %d, to = %+v\n", min, to)
		}
	}
	// fmt.Printf("   returning %+v\n", to)
	return to
}

func findAttackTarget(u *Unit) *Unit {
	adjacent := []Point{
			{u.x, u.y-1},
			{u.x -1 , u.y},
			{u.x + 1, u.y},
			{u.x , u.y+1}}

	var targets []*Unit

	for _,a := range adjacent {
		t := units.at(a.x, a.y)
		if t!= nil && t.race != u.race && t.hp > 0 {
			targets = append(targets, t)
		}
	}

	if len(targets) == 0 {
		return nil
	}

	sort.Slice(targets, func (a,b int) bool {
			return targets[a].hp < targets[b].hp
	})
	return targets[0]
}