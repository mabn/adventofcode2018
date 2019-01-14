package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"regexp"
	"strconv"
	"sort"
)

type Operation func(int, int, int, []int)
type Case struct {
	before []int
	opcode []int
	after []int	
}
type Mapping struct {
	opName string
	opId int
}
const(
	YES int = 1
	NO int = 2
)
var okMapping map[Mapping]int


func main() {
	okMapping = map[Mapping]int{}
	_,program := load()
	
	//part1(cases)

	// solveOpcodes(cases)

	mapping := []string{
		"gtir",
		"setr",
		"bori",
		"gtrr",
		"gtri",
		"eqir",
		"seti",
		"eqri",
		"eqrr",
		"borr",
		"addr",
		"mulr",
		"bani",
		"muli",
		"banr",
		"addi",
	}

	run(program, mapping)
}

func run(program []string, mapping []string) {
	ops := allOps()
	r := []int{0,0,0,0}
	for i := range program {
		line := program[i]

		
		
		nums := findNumbers(line)
		opId := nums[0]
		a := nums[1]
		b := nums[2]
		c := nums[3]
		fmt.Println("nums:", opId, a, b, c)
		op := ops[mapping[opId]]

		op(a,b,c,r)

		fmt.Printf("r=%+v\n", r)
	}
}

func solveOpcodes(cases []Case) {
	ops := allOps()

	// mapping[opcode] == function name
	var mapping []string
	for k := range ops {
		mapping = append(mapping, k)
	}
	sort.Strings(mapping)
	

	for name := range ops {
		for i:=0;i<=16;i++{
			checkMapping(name, i, cases, ops)
		}
	}

	for i := range okMapping {
		if okMapping[i] == YES {
			fmt.Println("ok: ", i)
		}
	}

	permutations(mapping, 0, func(pivot int) bool {
	 	return check(mapping, pivot, cases, ops)
	})

	fmt.Println("final mapping:", mapping)


}

func permutations(mapping []string, pivot int, f func(int) bool) bool {
	// short circuit recursive calls
	if !f(pivot) {
		return false
	}
	if pivot == len(mapping) {
		fmt.Println("SOLUTION!!!")
		return true
	}
	for i := pivot; i <= len(mapping) - 1; i++ {
		swap(pivot, i, mapping)
		ok := permutations(mapping, pivot + 1, f)
		if ok {
			return true
		}
		swap(pivot, i, mapping)
	}
	return false
}

func check(mapping []string, pivot int, cases[] Case, ops map[string]Operation) bool {
	if pivot == 0 {
		return true
	}

	fmt.Println("checking mapping: ", mapping[:pivot], "pivot=", pivot)
	opId := pivot-1
	opName := mapping[opId]
	m := Mapping{opName, opId}
	if okMapping[m] == 0 {
		panic("okMapping[m] == 0")
	}
	if okMapping[m] == YES {
		return true
	}
	fmt.Println(" does not work: ", m)
	return false
}

func checkMapping(opName string, opId int, cases[] Case, ops map[string]Operation) bool {
	m := Mapping{opName, opId}
	if okMapping[m] != 0 {
		//fmt.Printf(" short cutting %s=%d -> %d\n", opName, opId, okMapping[m])
		return okMapping[m] == YES
	}
	tested := 0
	for _,example := range cases {
		if opId != example.opcode[0] {
			continue
		}
		
		op := ops[opName]
		a := example.opcode[1]
		b := example.opcode[2]
		c := example.opcode[3]
		r := dup(example.before)
		op(a,b,c,r)
		if !same(example.after, r) {
			//fmt.Printf("  does not work for %+v, operation %s\n", example, opName)
			okMapping[m] = NO
			return false
		}
		tested++
	}
	fmt.Printf("%s=%d ok for %d cases\n", opName, opId, tested)
	if tested > 0 {
		okMapping[m] = YES
		return true
	}
	okMapping[m] = NO
	return false
}

func swap(a,b int, s[] string) {
	tmp := s[a]
	s[a] = s[b]
	s[b] = tmp
}

func part1(cases []Case) {
	threePlusCnt := 0
	for _,example := range cases {
		matches := 0
		a := example.opcode[1]
		b := example.opcode[2]
		c := example.opcode[3]

		lastMatched := ""
		ops := allOps()
		for i := range ops {
			op := ops[i]
			r := dup(example.before)
			//fmt.Printf("%s in:  %v", opName, r)
			op(a, b, c, r)
			//fmt.Printf(" out: %v\n", r)
			if same(example.after, r) {
				matches++
				lastMatched = i
				fmt.Println("op matches:", i)
			}
		}
		fmt.Println(">>> matches:", matches)
		if matches >= 3 {
			threePlusCnt++
		}
		if matches == 1 {
			fmt.Printf("found %s -> %d\n", lastMatched, example.opcode[0])
		}
		
	}

	
	fmt.Println(threePlusCnt)
}

func allOps() map[string]Operation {
	return map[string]Operation{
		"addr": addr, // 10
		"addi": addi, // 15
		"mulr": mulr, // 11
		"muli": muli, // 13
		"banr": banr, // 14
		"bani": bani, // 12
		"borr": borr, // 9
		"bori": bori, // 2
		"setr": setr, // 1
		"seti": seti, // 6
		"gtir": gtir, // 0
		"gtri": gtri, // 4
		"gtrr": gtrr, // 3
		"eqir": eqir, // 5
		"eqri": eqri, // 7
		"eqrr": eqrr, // 8
	}
}
func load() ([]Case, []string) {
	f,err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	readingProgram := false
	var cases []Case
	var program []string
	for s.Scan() {
		if !strings.HasPrefix(s.Text(), "Before") && !readingProgram {
			s.Scan()
			readingProgram = true
			continue
		}
		if !readingProgram {
			c := Case{}
			c.before = findNumbers(s.Text())
			s.Scan()
			c.opcode = findNumbers(s.Text())
			s.Scan()
			c.after = findNumbers(s.Text())
			s.Scan()
			if s.Text() != "" {
				panic("expected empty line")
			}
			//fmt.Printf("%+v\n", c)
			cases = append(cases, c)
		}else{
			program = append(program, s.Text())
		}
	}
	//fmt.Println("loaded", len(cases), "cases")
	fmt.Println("program len:", len(program))
	return cases, program
}

func findNumbers(s string) []int {
	re := regexp.MustCompile("[0-9]+")
	numbers := re.FindAllString(s, -1)
	var ret []int
	for _,x := range numbers {
		i,err := strconv.Atoi(x)
		if err != nil {
			panic(err)
		}
		ret = append(ret, i )
	}
	return ret
}

func addr(a, b, c int, r []int) {
	r[c] = r[a] + r[b]
}

func addi(a, b, c int, r []int) {
	r[c] = r[a] + b
}

func mulr(a, b, c int, r []int) {
	r[c] = r[a] * r[b]
}

func muli(a, b, c int, r []int) {
	r[c] = r[a] * b
}

func banr(a, b, c int, r []int) {
	r[c] = r[a] & r[b]
}

func bani(a, b, c int, r []int) {
	r[c] = r[a] & b
}

func borr(a, b, c int, r []int) {
	r[c] = r[a] | r[b]
}

func bori(a, b, c int, r []int) {
	r[c] = r[a] | b
}

func setr(a, b, c int, r []int) {
	r[c] = r[a]
}
func seti(a, b, c int, r []int) {
	r[c] = a
}

func gtir(a, b, c int, r []int) {
	if a > r[b] {
		r[c] = 1
	} else {
		r[c] = 0
	}
}

func gtri(a, b, c int, r []int) {
	if r[a] > b {
		r[c] = 1
	} else {
		r[c] = 0
	}
}

func gtrr(a, b, c int, r []int) {
	if r[a] > r[b] {
		r[c] = 1
	} else {
		r[c] = 0
	}
}

func eqir(a, b, c int, r []int) {
	if a == r[b] {
		r[c] = 1
	} else {
		r[c] = 0
	}
}
func eqri(a, b, c int, r []int) {
	if r[a] == b {
		r[c] = 1
	} else {
		r[c] = 0
	}
}
func eqrr(a, b, c int, r []int) {
	if r[a] == r[b] {
		r[c] = 1
	} else {
		r[c] = 0
	}
}

func dup(r []int) []int {
	ret := make([]int, len(r))
	copy(ret, r)
	return ret
}

func same(x, y []int) bool {
	if len(x) != len(y) {
		return false
	}

	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}
