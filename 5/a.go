// 12:45 - 13:25
package main

import (
	"fmt"
	"container/list"
	"unicode"
	"io/ioutil"
	"os"
)
type Item struct {
	downcase rune
	upper bool
	str string
}

func main() {
	//input := "dabAcCaCBAcCcaDA"
	input := read()
	l := list.New()
	for i, _ := range input {
		char := rune(input[i])
		l.PushBack(Item{downcase: unicode.ToLower(char), upper: unicode.IsUpper(char), str: string(char)})
	}

	
	optimize(l)
	fmt.Println(l.Len())

	for unit := 'a'; unit <= 'z'; unit++ {
		c := copy(l)
		removeUnit(unit, c)
		optimize(c)
		fmt.Printf("removing %s -> %d\n", string(unit), c.Len())
	}


	//print(l)

	
}

func print(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Print(e.Value.(Item).str)
	}
	fmt.Println()
}

func react(l *list.List) {
	for e := l.Front(); e != nil && e.Next() != nil; e = e.Next() {
		curr := e.Value.(Item)
		next := e.Next().Value.(Item)
		if (curr.downcase == next.downcase) && (curr.upper != next.upper) {
			// fmt.Printf("merging %s %s\n", curr.str, next.str)

			l.Remove(e.Next())
			l.Remove(e)
		}
	}
}

func optimize(l *list.List) {
	len := l.Len()
	for {
		react(l)
		newLen := l.Len()
		if len == newLen {
			break
		}
		len = newLen
	}
}

func read() string {
	f,err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bytes,err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func copy(l *list.List) *list.List {
	c := list.New()
	c.PushBackList(l)
	return c
}

func removeUnit( unit rune, l *list.List) {
	e := l.Front()
	for e != nil  {
		//fmt.Printf("element %s unit %s equal %s\n", string(e.Value.(Item).downcase), string(unit), e.Value.(Item).downcase == unit)
		if e.Value.(Item).downcase == unit {
			next := e.Next()
			l.Remove(e)
			e = next
		}else{
			e = e.Next()
		}
	}
}
