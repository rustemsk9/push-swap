package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

//	type stack struct {
//		curr []float32
//		targetStack *stack
//	}
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
		// Use plan-based selective pushes
		maxPushes := min(len(*A)/3, 2)
		pushed := 0

		for len(*A) > 3 && pushed < maxPushes {
			bestCost := -1
			bestIdxA := 0
			maxConsider := min(len(*A), 4)
			for idxA := 0; idxA < maxConsider; idxA++ {
				x := (*A)[idxA]
				idxB := findInsertIdxBDesc(*B, x)
				p := BestPlan(len(*A), idxA, len(*B), idxB)
				cost := planCost(p)
				if bestCost == -1 || cost < bestCost {
					bestCost = cost
					bestIdxA = idxA
				}
			}
			x := (*A)[bestIdxA]
			idxB := findInsertIdxBDesc(*B, x)
			plan := BestPlan(len(*A), bestIdxA, len(*B), idxB)
			ExecutePlan(A, B, plan)
			pb(A, B)
			pushed++
			MaybeSS(A, B)
		}

		if len(*A) >= 4 && !(*A).isSorted() {
			if len(*A) >= 2 && (*A)[0] > (*A)[1] {
				ra(A, B)
			}
		}

		// Solve remaining A
		solve(A, B)
		// Return all from B to A
		for len(*B) > 0 {
			pa(A, B)
		}
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

// getPivots возвращает два порога: ~33% и ~66% перцентили.
func getPivots(a stack) (p1, p2 float64) {
	n := len(a)
	if n == 0 {
		return 0, 0
	}
	cp := make([]float64, n)
	copy(cp, a)
	sort.Float64s(cp)
	p1 = cp[n/3]
	p2 = cp[(2*n)/3] // 33% и 66% перцентили
	if p1 > p2 {
		p1, p2 = p2, p1
	}
	return
}

// indexOfMax возвращает индекс максимального элемента стека.
func indexOfMax(b stack) int {
	if len(b) == 0 {
		return 0
	}
	maxIdx := 0
	for i := 1; i < len(b); i++ {
		if b[i] > b[maxIdx] {
			maxIdx = i
		}
	}
	return maxIdx
}

// rotateBToTop крутит только B, чтобы b[idx] оказался на вершине.
func rotateBToTop(A, B *stack, idx int) {
	n := len(*B)
	if n == 0 || idx < 0 || idx >= n {
		return
	}
	if idx <= n/2 {
		for i := 0; i < idx; i++ {
			rb(A, B)
		}
	} else {
		for i := 0; i < n-idx; i++ {
			rrb(A, B)
		}
	}
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

// // findInsertIdxBDesc возвращает индекс в B, который нужно поднять на верх,
// // чтобы после pb элемент из A оказался на вершине B и B оставался по убыванию.
// // Линейная проверка по «круговой» паре prev->curr; если подходящей щели нет,
// // вставляем после максимума.
func findInsertIdxBDesc(b stack, x float64) int {
	n := len(b)
	if n == 0 {
		return 0
	}
	maxIdx := 0
	for i := 1; i < n; i++ {
		if b[i] > b[maxIdx] {
			maxIdx = i
		}
	}
	// Ищем позицию i, куда x «вписывается» между prev и curr: prev >= x >= curr
	for i := 0; i < n; i++ {
		prev := b[(i-1+n)%n]
		curr := b[i]
		if prev >= x && x >= curr {
			return i
		}
	}
	// x либо больше max, либо меньше min — вставляем после максимума
	return (maxIdx + 1) % n
}

// canSolveSmall checks if a small stack (3 or fewer elements) can be sorted efficiently
func canSolveSmall(a stack) bool {
	return len(a) <= 3
}

func main() {
	if len(os.Args) == 1 {
		os.Exit(2)
	}
	stackA := stackFromArgNums(os.Args[1:])
	stackB := stack{}

	// Special handling for small stacks
	if len(stackA) <= 3 {
		solve(&stackA, &stackB)
		return
	}

	// For the optimal strategy: be more selective about pushes
	maxPushes := min(len(stackA)/3, 2) // Push at most 2 elements for most cases
	pushed := 0

	for len(stackA) > 3 && pushed < maxPushes {
		bestCost := -1
		bestIdxA := 0

		// Only consider the first few elements for pushing (more selective)
		maxConsider := min(len(stackA), 4)
		for idxA := 0; idxA < maxConsider; idxA++ {
			x := stackA[idxA]
			idxB := findInsertIdxBDesc(stackB, x)
			p := BestPlan(len(stackA), idxA, len(stackB), idxB)
			cost := planCost(p)
			if bestCost == -1 || cost < bestCost {
				bestCost = cost
				bestIdxA = idxA
			}
		}

		// Execute the best plan
		x := stackA[bestIdxA]
		idxB := findInsertIdxBDesc(stackB, x)
		plan := BestPlan(len(stackA), bestIdxA, len(stackB), idxB)
		ExecutePlan(&stackA, &stackB, plan)
		pb(&stackA, &stackB)
		pushed++

		// Try to optimize with swaps after each push
		MaybeSS(&stackA, &stackB)
	}

	if len(stackA) >= 4 && !stackA.isSorted() {
		// Check if ra would help improve the order
		if len(stackA) >= 2 && stackA[0] > stackA[1] {
			// Try ra to see if it improves things
			ra(&stackA, &stackB)
		}
	}
	// Solve remaining A (should be mostly sorted now)
	solve(&stackA, &stackB)

	// Return all elements from B to A
	for len(stackB) > 0 {
		pa(&stackA, &stackB)
	}

	fmt.Println(stackA)
}
