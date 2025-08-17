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

// RotPlan captures how many times to call each rotation for the minimal-cost plan.
type RotPlan struct {
	rr, rrr int // overlapped rotations
	ra, rra int // residual rotations on A
	rb, rrb int // residual rotations on B
}

// rotationsToTop returns the number of up-rotations (ra/rb) and down-rotations (rra/rrb)
// needed to bring index idx to the top in a stack of length n.
func rotationsToTop(n, idx int) (up, down int) {
	if n <= 0 || idx < 0 || idx >= n {
		return 0, 0
	}
	up = idx
	down = n - idx
	return up, down
}

// planCost sums the total number of operations in a RotPlan.
func planCost(p RotPlan) int {
	return p.rr + p.rrr + p.ra + p.rra + p.rb + p.rrb
}

// BestPlan computes the minimal plan to bring A[idxA] and B[idxB] to the top,
// exploiting rr/rrr where beneficial.
func BestPlan(lenA, idxA, lenB, idxB int) RotPlan {
	aUp, aDown := rotationsToTop(lenA, idxA)
	bUp, bDown := rotationsToTop(lenB, idxB)

	// Plan 1: both up (rr)
	p1 := RotPlan{}
	p1.rr = min(aUp, bUp)
	p1.ra = aUp - p1.rr
	p1.rb = bUp - p1.rr

	// Plan 2: both down (rrr)
	p2 := RotPlan{}
	p2.rrr = min(aDown, bDown)
	p2.rra = aDown - p2.rrr
	p2.rrb = bDown - p2.rrr

	// Plan 3: A up, B down
	p3 := RotPlan{}
	p3.ra = aUp
	p3.rrb = bDown

	// Plan 4: A down, B up
	p4 := RotPlan{}
	p4.rra = aDown
	p4.rb = bUp

	// Choose the cheapest plan
	best := p1
	bestCost := planCost(p1)

	if c := planCost(p2); c < bestCost {
		best, bestCost = p2, c
	}
	if c := planCost(p3); c < bestCost {
		best, bestCost = p3, c
	}
	if c := planCost(p4); c < bestCost {
		best, bestCost = p4, c
	}
	return best
}

// ExecutePlan applies the plan using your primitive ops.
// After this, A[idxA] and B[idxB] will be at the top (if idxB was valid for lenB > 0).
// You typically follow with pb/pa as intended.
func ExecutePlan(A, B *stack, p RotPlan) {
	for i := 0; i < p.rr; i++ {
		rr(A, B)
	}
	for i := 0; i < p.rrr; i++ {
		rrr(A, B)
	}
	for i := 0; i < p.ra; i++ {
		ra(A, B)
	}
	for i := 0; i < p.rra; i++ {
		rra(A, B)
	}
	for i := 0; i < p.rb; i++ {
		rb(A, B)
	}
	for i := 0; i < p.rrb; i++ {
		rrb(A, B)
	}
}

// MaybeSS tries to use ss when both stacks benefit. Otherwise performs needed single swaps.
// - For A (ascending target): swap if A[0] > A[1].
// - For B (descending maintenance): swap if B[0] < B[1].
func MaybeSS(A, B *stack) {
	doA := len(*A) >= 2 && (*A)[0] > (*A)[1]
	doB := len(*B) >= 2 && (*B)[0] < (*B)[1]

	switch {
	case doA && doB:
		ss(A, B)
	case doA:
		sa(A, B)
	case doB:
		sb(A, B)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// functions := map[string]func(A, B *stack){
	// 	"sa": sa, "sb": sb, "ss": ss, "pa": pa, "pb": pb, "ra": ra, "rb": rb, "rr": rr, "rra": rra, "rrb": rrb, "rrr": rrr}

	if len(os.Args) == 1 {
		os.Exit(2)
	}
	stackA := stackFromArgNums(os.Args[1:])
	stackB := stack{}
	// Example: push all A to B using optimal rotations
	for len(stackA) > 0 {
		bestCost := -1
		bestIdxA := 0
		// bestPlan := RotPlan{}
		for idxA := 0; idxA < len(stackA); idxA++ {
			idxB := 0 // можно улучшить: искать лучшую позицию для вставки в B
			p := BestPlan(len(stackA), idxA, len(stackB), idxB)
			cost := planCost(p)
			if bestCost == -1 || cost < bestCost {
				bestCost = cost
				bestIdxA = idxA
				// bestPlan = p
			}
		}
		plan := BestPlan(len(stackA), bestIdxA, len(stackB), 0)
		ExecutePlan(&stackA, &stackB, plan)
		pb(&stackA, &stackB)
		MaybeSS(&stackA, &stackB)
	}
	for _, x := range stackA {
		fmt.Print(int(x), "B ")
	}
	fmt.Println()
	for _, x := range stackB {
		fmt.Print(int(x), "A ")
	}
	return
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
