package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Race byte

type Point struct {
	x, y int
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
	load(os.Args[2])

	board.print(0)

	r := 1
	for ; units.endOfCombat() && r < 1000000; r++ { // rounds
		// sort units by y.x

		round()
		units.sort()
		board.print(r)
		fmt.Printf("==> Rnd %d finished\n", r)
		units.removeDead()
	}

	sum := 0
	for _, u := range units.list {
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

func (units *Units) getByRace(x int, y int, race Race) *Unit {

	u := units.pos[x][y]
	if u != nil && u.race == race {
		return u
	}
	return nil
}

func (board Board) print(round int) {
	fmt.Printf("Round: %d\n", round)
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
	for _, u := range units.list {
		fmt.Printf("%+v\n", u)
	}
}

// create distance map to each unit of specified race
func distmap(x int, y int) [][]int {
	width := len(board)
	height := len(board[0])
	d := make([][]int, height)
	for i := range d {
		d[i] = make([]int, width)
	}

	distmapFill(1, x, y, d)

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
				fmt.Print(" .")
			} else {
				fmt.Printf("%2d", d[x][y])
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

	cpy := make([]*Unit, len(units.list))
	copy(cpy, units.list)
	for _, unit := range cpy {
		if unit.hp <= 0 {
			continue
		}

		// fmt.Printf("processing unit %d,%d\n", unit.x, unit.y)

		
		

		destination := findInRangeDestination(unit)
		to := findMoveTarget(unit, destination)

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
	if target.hp <= 0 {
		units.remove(target)
	}
}

func (units *Units) move(unit *Unit, dx int, dy int) {
	// fmt.Printf("%d,%d moving by %d,%d\n", unit.x, unit.y, dx, dy)
	fmt.Printf("%d,%d moving to %d,%d\n", unit.x, unit.y, unit.x+dx, unit.y+dy)
	units.pos[unit.x][unit.y] = nil
	unit.x += dx
	unit.y += dy
	units.pos[unit.x][unit.y] = unit
}

func (units *Units) remove(u *Unit) {
	units.pos[u.x][u.y] = nil

	for i,_ := range units.list {
		if units.list[i] == u {
			units.list = append(units.list[:i], units.list[i+1:]...)
			return
		}
	}
}

func (units *Units) removeDead() {
	newlist := []*Unit{}
	for _, u := range units.list {
		if u.hp > 0 {
			newlist = append(newlist, u)
		} else {
			fmt.Printf("Removing dead unit %+v\n", u)
			units.pos[u.x][u.y] = nil
		}

	}
	units.list = newlist
}

func (units *Units) endOfCombat() bool {
	counts := []int{0, 0}

	for _, u := range units.list {
		counts[u.race]++
	}
	return counts[0] > 0 && counts[1] > 0
}

/* Finds in-range point - adjacent to an enemy, reachable by unit
 * and closest to the unit
 */
func findInRangeDestination(unit *Unit) *Point {
	d := distmap(unit.x, unit.y)
	enemyRace := unit.race.opposite()
	//fmt.Printf(" distmap from %d,%d\n", unit.x, unit.y)
	//printDist(d)
	if unit.x == 5 && unit.y == 15 {
		printDist(d)
	}

	// find closest point in d adjacent to enemy unit
	minDist := 10000
	var closest *Point
	for j := 1; j < len(d[0])-1; j++ {
		for i := 1; i < len(d)-1; i++ {
			for _, adj := range (Point{i, j}.adjacent()) {
				t := units.getByRace(adj.x, adj.y, enemyRace)
				if t != nil && d[i][j] < minDist && d[i][j] > 0 {
					minDist = d[i][j]
					closest = &Point{i, j}
				}
			}
		}
	}
	if closest != nil && (unit.x != closest.x || unit.y != closest.y) {
		fmt.Printf("%d,%d will move towards %+v, mindist: %d\n", unit.x, unit.y, closest, minDist)
		return closest
	}
	
	return nil
}

func findMoveTarget(unit *Unit, dest *Point) *Point {
	if dest == nil {
		return nil
	}
	d := distmap(dest.x, dest.y)

	minDist := 100000
	var target Point
	for _, a := range (Point{unit.x, unit.y}.adjacent()) {
		if d[a.x][a.y] > 0 && d[a.x][a.y] < minDist {
			minDist = d[a.x][a.y]
			target = Point{a.x, a.y}
		}
	}
	
	// fmt.Printf("   returning %+v\n", to)
	return &target
}

func findAttackTarget(u *Unit) *Unit {
	adjacent := []Point{
		{u.x, u.y - 1},
		{u.x - 1, u.y},
		{u.x + 1, u.y},
		{u.x, u.y + 1}}

	var targets []*Unit

	for _, a := range adjacent {
		t := units.at(a.x, a.y)
		if t != nil && t.race != u.race && t.hp > 0 {
			targets = append(targets, t)
		}
	}

	if len(targets) == 0 {
		return nil
	}

	sort.Slice(targets, func(a, b int) bool {
		return targets[a].hp < targets[b].hp
	})
	return targets[0]
}

func (p Point) adjacent() []Point {
	return []Point{
		{p.x, p.y - 1},
		{p.x - 1, p.y},
		{p.x + 1, p.y},
		{p.x, p.y + 1}}
}

// 183519 too low. 186432 is wrong ;(
// correnct answer is : 195774