package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	state, rules := read()

	state = "......" + state + "...."
	offset := -6

	i := 0
	for ; i <= 140; i++ {
		fmt.Printf("%4d: %s %d off: %d\n", i, state, sumPlants(&state, offset), offset)
		state = next(state, rules, &offset)
	}

	fmt.Println((50000000000-1000)*78+80467)
}

func sumPlants(state *string, offset int) int {
	var sum int
	for i:= 0; i < len(*state);i++{
		if (*state)[i] == '#' {
			sum += i + offset
		}
	}
	return sum
}

func next(init string, rules map[string]string, offset *int) string {
	init = adjStart(init, offset)
	init = adjEnd(init)
	out := []byte(init)
	*offset *= 1

	for i := 2; i < len(init)-2; i++ {
		chunk := init[i-2 : i+3]
		if rules[chunk] == "#" {
			out[i] = '#'
		} else {
			out[i] = '.'
		}
	}
	return string(out)
}

func adjStart(init string, offset *int) string {
	for strings.HasPrefix(init, "......") {
		init = init[1:]
		*offset = *offset + 1
	}
	return init
}

func adjEnd(init string) string {
	for !strings.HasSuffix(init, "...") {
		init = init + "."
	}
	return init
}

func read() (string, map[string]string) {
	f, _ := os.Open("input.txt")
	defer f.Close()

	var init string
	fmt.Fscanf(f, "initial state: %s\n", &init)
	fmt.Fscanf(f, "\n")

	rules := map[string]string{}
	var rule, result string
	for {
		_, err := fmt.Fscanf(f, "%s => %s\n", &rule, &result)
		if err != nil {
			break
		}
		rules[rule] = result
	}
	return init, rules
}
