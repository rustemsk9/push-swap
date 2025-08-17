package main

import (
	"fmt"
	"os"
	"sort"
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

// ...existing code...

func main() {
	// functions := map[string]func(A, B *stack){
	// 	"sa": sa, "sb": sb, "ss": ss, "pa": pa, "pb": pb, "ra": ra, "rb": rb, "rr": rr, "rra": rra, "rrb": rrb, "rrr": rrr}

	if len(os.Args) == 1 {
		os.Exit(2)
	}
	stackA := stackFromArgNums(os.Args[1:])
	stackB := stack{}

	// 1) Многофазное разбиение A на B по двум пивотам, пока в A не останется <= 3
	// 1) Многофазное разбиение A на B по двум пивотам, пока в A не останется <= 3
	for len(stackA) > 3 {
		p1, p2 := getPivots(stackA)
		limit := len(stackA) // обрабатываем ровно текущий размер A
		for processed := 0; processed < limit; processed++ {
			x := stackA[0]
			switch {
			case x < p1:
				// Малые — отправляем в B.
				pb(&stackA, &stackB)
				// Агрегация вращений: если следующий на вершине A — "большой" (>= p2),
				// то мы бы сделали rb сейчас и ra на следующем шаге. Вместо этого делаем rr.
				if len(stackA) > 0 && stackA[0] >= p2 {
					rr(&stackA, &stackB)
				} else {
					rb(&stackA, &stackB)
				}
			case x < p2:
				// Средние — просто в B без дополнительных вращений.
				pb(&stackA, &stackB)
			default:
				// Большие — крутим A.
				ra(&stackA, &stackB)
			}
		}
	}

	// 2) Досортировать малый A
	switch len(stackA) {
	case 3:
		three(&stackA, &stackB)
	case 2:
		two(&stackA, &stackB)
	}

	// 3) Вернуть из B в A: каждый раз подтягиваем максимум в B на вершину и pa
	for len(stackB) > 0 {
		idx := indexOfMax(stackB)
		rotateBToTop(&stackA, &stackB, idx)
		pa(&stackA, &stackB)
	}

	// Готово: A по возрастанию, B пуст.

	for _, x := range stackA {
		fmt.Print(int(x), "A ")
	}
	fmt.Println()
	for _, x := range stackB {
		fmt.Print(int(x), "B ")
	}
	return
	// echo -e "pb\npb\nsb\nrb\npb\nsa\nrb\npb\nrb\npb\nrb\npb" | ./cheker "2 1 3 6 5 8"
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


// // findInsertIdxBDesc возвращает индекс в B, который нужно поднять на верх,
// // чтобы после pb элемент из A оказался на вершине B и B оставался по убыванию.
// // Линейная проверка по «круговой» паре prev->curr; если подходящей щели нет,
// // вставляем после максимума.
// func findInsertIdxBDesc(b stack, x float64) int {
//     n := len(b)
//     if n == 0 {
//         return 0
//     }
//     maxIdx := 0
//     for i := 1; i < n; i++ {
//         if b[i] > b[maxIdx] {
//             maxIdx = i
//         }
//     }
//     // Ищем позицию i, куда x «вписывается» между prev и curr: prev >= x >= curr
//     for i := 0; i < n; i++ {
//         prev := b[(i-1+n)%n]
//         curr := b[i]
//         if prev >= x && x >= curr {
//             return i
//         }
//     }
//     // x либо больше max, либо меньше min — вставляем после максимума
//     return (maxIdx + 1) % n
// }

// func main() {
//     // functions := map[string]func(A, B *stack){
//     // 	"sa": sa, "sb": sb, "ss": ss, "pa": pa, "pb": pb, "ra": ra, "rb": rb, "rr": rr, "rra": rra, "rrb": rrb, "rrr": rrr}

//     if len(os.Args) == 1 {
//         os.Exit(2)
//     }
//     stackA := stackFromArgNums(os.Args[1:])
//     stackB := stack{}

//     // Переносим из A в B, поддерживая B по убыванию.
//     for len(stackA) > 0 {
//         bestCost := -1
//         bestIdxA := 0
//         bestPlan := RotPlan{}
//         for idxA := 0; idxA < len(stackA); idxA++ {
//             x := stackA[idxA]
//             idxB := findInsertIdxBDesc(stackB, x)
//             p := BestPlan(len(stackA), idxA, len(stackB), idxB)
//             cost := planCost(p)
//             if bestCost == -1 || cost < bestCost {
//                 bestCost = cost
//                 bestIdxA = idxA
//                 bestPlan = p
//             }
//         }
//         // На случай изменения стека: пересчитываем план под выбранные индексы
//         x := stackA[bestIdxA]
//         idxB := findInsertIdxBDesc(stackB, x)
//         plan := BestPlan(len(stackA), bestIdxA, len(stackB), idxB)
//         ExecutePlan(&stackA, &stackB, plan)
//         pb(&stackA, &stackB)
//         MaybeSS(&stackA, &stackB)
//     }

//     // Возвращаем из B в A — теперь A по возрастанию.
//     for len(stackB) > 0 {
//         pa(&stackA, &stackB)
//     }
