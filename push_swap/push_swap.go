package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type stack []float64

func dupChecker() func(int) bool {
	duplicate := map[int]bool{}
	return func(num int) bool {
		if _, ok := duplicate[num]; !ok {
			duplicate[num] = true
			return false
		} else {
			return true
		}
	}
}

func stackFromArgNums(argv []string) stack {
	isDup := dupChecker()
	stackA := stack{}
	var smallest int
	for _, str := range argv {
		split := strings.Split(str, " ")
		for _, x := range split {
			num, err := strconv.Atoi(x)
			if num < smallest {
				smallest = num
			}
			if err != nil {
				fmt.Println("Bad input")
				os.Exit(2)
			}
			if isDup(num) {
				fmt.Println(num, "is duplicated, Error")
				os.Exit(2)
			}
			stackA = append(stackA, float64(num))
		}
	}
	return stackA
}

type valueStack struct {
	A, B *stack
}
type sizeStack struct {
	A, B *stack
}

func solve(A, B *stack) {
	switch len(*A) {
	case 4, 5:
		four_Five(A, B)
	case 3:
		three(A, B)
	case 2:
		two(A, B)
	default:
		values := valueStack{A, B}

		mergeSort(values)
	}
}

func main() {
	// functions := map[string]func(A, B *stack){
	// 	"sa": sa, "sb": sb, "ss": ss, "pa": pa, "pb": pb, "ra": ra, "rb": rb, "rr": rr, "rra": rra, "rrb": rrb, "rrr": rrr}

	if len(os.Args) == 1 {
		os.Exit(2)
	}
	stackA := stackFromArgNums(os.Args[1:])
	stackB := stack{}

	solve(&stackA, &stackB)

	for _, x := range stackA {
		fmt.Print(int(x), " ")
	}
	fmt.Println()
	for _, x := range stackB {
		fmt.Print(int(x), " ")
	}
	fmt.Println()
}

func pushAlot(A, B *stack, amount int) {
	for amount > 0 { // Sorted then the left over? random? // argest set then random?
		pb(A, B)
		amount--
	}
}

func pushAlotRev(A, B *stack, amount int) {
	for amount > 0 { // Sorted then the left over? random? // argest set then random?
		rra(A, B)
		pb(A, B)
		amount--
	}
}

const (
	right bool = true
	left  bool = false
)
