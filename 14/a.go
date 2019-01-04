package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	a := []int{3, 7}
	x := 0
	y := 1

	print(a, x, y)

	//input := 765071
	//part2: 
	input := 50000000
	for i := 0; len(a) < input + 10; i++ {
		a = next(a, x, y)
		x = moveElf(a, x)
		y = moveElf(a, y)
		
	}
	
	s := ""
	for _,v:=range a[input:input+10] {
		s = s + strconv.Itoa(v)
	}


	fmt.Println(s)

	var out []byte
	for i,_:=range a {
		out = append(out, byte(a[i] + '0'))
	}

	fmt.Println(string(out[input:input+10]))

	outs := string(out)
	fmt.Println(strings.Index(outs, "51589"))
	fmt.Println(strings.Index(outs, "01245"))
	fmt.Println(strings.Index(outs, "92510"))
	fmt.Println(strings.Index(outs, "59414"))
	fmt.Println(strings.Index(outs, "765071"))
}
func next(a []int, x int, y int) []int {
	n := a[x] + a[y]
	var digits []int
	for n > 0 {
		digits = append(digits, n % 10)
		n = n / 10
	}

	var reversed []int
	for i, _ := range digits {
		reversed = append(reversed, digits[len(digits)-i-1])
	}
	if len(reversed) == 0 {
		reversed = append(reversed, 0)
	}
	return append(a, reversed...)
}

func moveElf(a []int, pos int) (int) {
	r :=  (pos + a[pos] + 1) % len(a)
	return r
}

func print(a []int, x int, y int) {
	for i, _ := range a {
		fmt.Print(a[i])
	}
	fmt.Println()
	for i, _ := range a {
		if x == i {
			fmt.Print("x")
		} else if y == i {
			fmt.Print("y")
		} else {
			fmt.Print(" ")	
		}
	}
	fmt.Println()
}
