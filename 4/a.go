package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type GuardSleep struct {
	guard int
	duration int
}

func main() {
	//guard_id -> [start,stop], [start,stop], [start,stop]

	guards := read_input()
	sleeps := compute_sleep_durations(guards)
	fmt.Println(sleeps)
	guard_id := find_most_sleepy(sleeps)
	fmt.Println("most sleepy", guard_id)

	minute,_:=find_most_freq_minute(guards[guard_id])
	fmt.Println("answer A", minute*guard_id)

	for guard,sleeps := range guards {
		minute, times := find_most_freq_minute(sleeps)

		fmt.Printf("%d %d %d\n", guard, minute, times)
	}
}

type Break struct {
	from int
	to int
}

func read_input() map[int][]Break {
	file, _ := os.Open("input.txt")
	defer file.Close()

	p, err := regexp.Compile(".* (\\d+):(\\d+)\\] (?:(falls asleep)|(wakes up)|(?:Guard #(\\d+) begins shift))")
	if err != nil {
		panic(err)
	}

	guards := map[int][]Break{}
	current_guard := -1

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := p.FindStringSubmatch(line)

		hour, _ := strconv.Atoi(matches[1])
		minute, _ := strconv.Atoi(matches[2])

		if len(matches[3]) > 0 {
			// handle falls asleep
			guards[current_guard] = append(guards[current_guard], Break{from: minute, to: 59})
		}
		if len(matches[4]) > 0 {
			// handle wakes up
			num_breaks := len(guards[current_guard])
			guards[current_guard][num_breaks-1].to = minute
		}
		if len(matches[5]) > 0 {
			// handle begins shift
			current_guard, _ = strconv.Atoi(matches[5])
		}

		// fmt.Println("matches:")
		// for i, m := range matches {
		// 	fmt.Println("  ", i, m)
		// }

		fmt.Printf("current_guard: %d time; %d:%d\n", current_guard, hour, minute)
		// fmt.Println(guards)
	}
	return guards
}

func compute_sleep_durations(guards map[int][]Break) []GuardSleep {
	sleeps := []GuardSleep{}
	for k,breaks := range guards {
		sum := 0
		for _,b := range breaks {
			sum += b.to - b.from
		}
		sleeps = append(sleeps, GuardSleep{guard: k, duration: sum})
	}
	return sleeps
}

func find_most_sleepy(sleeps []GuardSleep) int {
	max_sleep := 0
	guard := -1

	for _,v := range sleeps {
		if max_sleep < v.duration {
			max_sleep = v.duration
			guard = v.guard
		}
	}
	return guard
}

func find_most_freq_minute(sleeps []Break) (int,int) {
	// fmt.Println(sleeps)
	max_occurrences := -1
	most_occurring := -1
	for min:=0; min < 60; min++{
		occurrences := 0
		for _,sleep := range sleeps {
			if sleep.from <= min && sleep.to > min {
				occurrences++
			}
		}
		if occurrences > max_occurrences {
			max_occurrences = occurrences
			most_occurring = min
		}

		// fmt.Printf("%d -> %d\n", min, occurrences)
	}
	return most_occurring,max_occurrences
}