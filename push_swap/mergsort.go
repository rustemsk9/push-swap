package main

import (
	"fmt"
	"os"
)

func initSize2(value valueStack, size sizeStack) { //maybe do drops based on size?
	var count float64
	count++
	for i, _ := range (*value.A)[1:] {
		if (*value.A)[i] < (*value.A)[i+1] {
			count++
		} else {
			*size.A = append(*size.A, count)
			count = 1
		}
	}
	*size.A = append(*size.A, count)
	count = 1
	for i, _ := range (*value.B)[1:] {
		if (*value.B)[i] < (*value.B)[i+1] {
			count++
		} else {
			*size.B = append(*size.B, count)
			count = 1
		}
	}
	*size.B = append(*size.B, count)
	fmt.Fprintln(os.Stderr, "+++++++++++++++++++++++++++++")
	fmt.Fprintln(os.Stderr, (*value.A), *size.A, len(*size.A))
	fmt.Fprintln(os.Stderr, (*value.B), *size.B, len(*size.B))
	fmt.Fprintln(os.Stderr, "+++++++++++++++++++++++++++++")

}

func insertion(A, B *stack, side bool, Asize, Bsize int) {
	var p_ func(A, B *stack)
	var r_ func(A, B *stack)
	var compare func(A, B *stack) bool
	var push int
	var shift int
	if side == left {
		p_ = pa
		r_ = ra
		push = Bsize
		shift = Asize
		compare = func(A, B *stack) bool {
			if (*A)[0] > (*B)[0] {
				return true
			}
			return false
		}
	} else {
		p_ = pb
		r_ = rb
		push = Asize
		shift = Bsize
		compare = func(A, B *stack) bool {
			if (*A)[0] < (*B)[0] {
				return true
			}
			return false
		}
	}
	for push > 0 || shift > 0 {
		if shift == 0 || push > 0 && len(*A) != 0 && len(*B) != 0 && compare(A, B) {
			p_(A, B)
			push--
			shift++
		} else {
			r_(A, B)
			shift--
		}
		if check(A, B) {
			fmt.Println("CHECK PASSED")
		}
	}
}

func updateStacks(size sizeStack, side bool) {
	numA, _ := size.A.pop()
	numB, _ := size.B.pop()

	var tmp *stack
	if side == left {
		tmp = size.A
	} else {
		tmp = size.B
	}
	*tmp = append(*tmp, numA+numB)
	fmt.Fprintln(os.Stderr, "A:", *size.A, "::", len(*size.A))
	fmt.Fprintln(os.Stderr, "B:", *size.B, "::", len(*size.B))
}

func mergeSort(values valueStack) {
	pushAlot(values.A, values.B, len(*values.A)/2)

	size := sizeStack{&stack{}, &stack{}}
	initSize2(values, size)
	stablelize(values, size)
	side := left
	fmt.Fprintln(os.Stderr, "A:", values.A, "\nB:", values.B)
	for len(*size.A) >= 1 && len(*size.B) >= 1 {
		insertion(values.A, values.B, side, int((*size.A)[0]), int((*size.B)[0]))
		fmt.Fprintln(os.Stderr, "A:", values.A, "\nB:", values.B)
		updateStacks(size, side)
		fmt.Fprintln(os.Stderr, "---------------")
		side = !side
	}
}

func stablelize(values valueStack, size sizeStack) { // If b is already sorted just push it all over
	// newAsize := stack{}
	// newBsize := stack{}
	side := left
	// sizeStack := sizeStack{&stack{}, &stack{}}

	fmt.Fprintln(os.Stderr, "STABLELIZING")
	fmt.Fprintln(os.Stderr, "A:", values.A, "\nB:", values.B)
	for len(*size.A) < len(*size.B) {
		insertion(values.A, values.B, left, int((*size.A)[0]), int((*size.B)[0]))
		fmt.Fprintln(os.Stderr, "-A:", values.A, "\nB:", values.B, len(*size.B))
		updateStacks(size, left)
		fmt.Fprintln(os.Stderr, "---------------")
		side = !side
	}
	for len(*size.A) >= len(*size.B)+1 {
		insertion(values.A, values.B, right, int((*size.A)[0]), int((*size.B)[0]))
		fmt.Fprintln(os.Stderr, "+A:", values.A, "::", len(*size.A), "\nB:", values.B, len(*size.B))
		updateStacks(size, right)
		fmt.Fprintln(os.Stderr, "---------------")
		side = !side
	}
	fmt.Fprintln(os.Stderr, "DONE")

}
