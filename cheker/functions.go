package main

import (
	"fmt"
	"os"
)

func sa(A, B *stack) {
	A.swap()
}

func sb(A, B *stack) {
	B.swap()
}

func ss(A, B *stack) {
	A.swap()
	B.swap()
}

func pa(A, B *stack) {
	num, _ := B.pop()
	*A = prepend(*A, num)
}

func pb(A, B *stack) {
	num, _ := A.pop()
	*B = prepend(*B, num)
}

func ra(A, B *stack) {
	A.rotateUp()
}

func rb(A, B *stack) {
	B.rotateUp()
}

func rr(A, B *stack) {
	A.rotateUp()
	B.rotateUp()
}

func rra(A, B *stack) {
	A.rotateDown()
}

func rrb(A, B *stack) {
	B.rotateDown()
}

func rrr(A, B *stack) {
	A.rotateDown()
	B.rotateDown()
}

func check(A, B *stack) {
	if A.isSorted() && B.isEmpty() {
		fmt.Println("OK")
		os.Exit(0)
	} else {
		fmt.Println("KO")
		os.Exit(1)
	}
}
