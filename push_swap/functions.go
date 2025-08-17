package main

import (
	"fmt"
)

func sa(A, B *stack) {
	A.swap()
	fmt.Println("sa")
}

func sb(A, B *stack) {
	B.swap()
	fmt.Println("sb")
}

func ss(A, B *stack) {
	A.swap()
	B.swap()
	fmt.Println("ss")
}

func pa(A, B *stack) {
	num, _ := B.pop()
	*A = prepend(*A, num)
	fmt.Println("pa")
}

func pb(A, B *stack) {
	num, _ := A.pop()
	*B = prepend(*B, num)
	fmt.Println("pb")
}

func ra(A, B *stack) {
	A.rotateUp()
	fmt.Println("ra")
}

func rb(A, B *stack) {
	B.rotateUp()
	fmt.Println("rb")
}

func rr(A, B *stack) {
	A.rotateUp()
	B.rotateUp()
	fmt.Println("rr")
}

func rra(A, B *stack) {
	A.rotateDown()
	fmt.Println("rra")
}

func rrb(A, B *stack) {
	B.rotateDown()
	fmt.Println("rrb")
}

func rrr(A, B *stack) {
	A.rotateDown()
	B.rotateDown()
	fmt.Println("rrr")
}

func check(A, B *stack) bool {
	if A.isSorted() && B.isEmpty() {
		return true
	} else {
		return false
	}
}
