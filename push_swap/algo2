package main

import (
	"fmt"
)

func algoTwo(A, B *stack) {
	pushHalfSorted(A, B, 2)
	values := valueStack{A, B}
	size2 := sizeStack{&stack{}, &stack{}}
	initSize2(values, size2)
	fmt.Println(*size2.A, *size2.B)
	bringBack(A, B)
}

func bringBack(A, B *stack) {
	for len(*B) > 0 {
		pa(A, B)
		if (*A)[0] > (*A)[1] {
			sa(A, B)
		}
	}
}

func pushHalfSorted(A, B *stack, half int) {
	// defer fmt.Println(*A, *B)
	if len(*A) <= half {
		return
	}
	pivotA, _ := A.getPivot()
	pivotB, _ := B.getPivot()
	for len(*A) > 0 && (*A)[0] != pivotA {
		if len(*B) > 1 && (*A)[0] > (*A)[1] && (*B)[0] < (*B)[1] {
			ss(A, B)
		} else if (*A)[0] > pivotA {
			if (*B)[0] < pivotB { //
				// fmt.Println(pivotB)
				rr(A, B)
			} else {
				ra(A, B)
			}
		} else {
			pb(A, B)
			pivotB, _ = B.getPivot()
		}
	}
	pb(A, B)
	pushHalfSorted(A, B, half)
}

// if (*A)[0] < pivotA && (*A)[0] > (*A)[1] && (*B)[0] < (*B)[1] {
// 	ss(A, B)
// 	//push twice?
// } else if (*A)[0] > pivotA && errB == nil && len(*B) > 0 && (*B)[0] < pivotB {
// 	rr(A, B)
// 	// 	//if top 2 of A are below pivot(large number to be ra'd) && top 2 of B are (small numberes that should be rr)
// 	// 	// A, top[0] > bot[1], B, top[0] > bot[1] swap
// 	// 	//then push pb twice
// 	// }
// 	// else if (*A)[0] > pivotA { //&& errB != nil && (*B)[0] > pivotB {
// 	// 	rr(A, B)
