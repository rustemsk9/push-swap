package main

func four_Five(A, B *stack) {
	if len(*A) <= 3 {
		three(A, B)
		return
	}

	// For 4 elements, try strategic approach similar to optimal solution
	if len(*A) == 4 {
		// Check if we can benefit from ra + sa + rrr pattern
		if len(*B) >= 2 && (*B)[0] < (*B)[1] {
			// Try the optimal pattern: ra, sa, rrr
			ra(A, B)               // Rotate A
			if (*A)[0] > (*A)[1] { // If swap would help
				ss(A, B) // Swap both if B benefits too
			} else {
				sa(A, B) // Otherwise just swap A
			}
			if !A.isSorted() {
				rrr(A, B) // Reverse rotate both
			}
		} else {
			// Fallback: try rotations and swaps
			ra(A, B)
			if (*A)[0] > (*A)[1] {
				sa(A, B)
			}
			if !A.isSorted() {
				rra(A, B)
			}
		}
	} else {
		// For 5 elements, be conservative
		three(A, B)
	}
}

func two(A, B *stack) {
	if A.isSorted() {
		return
	}
	// For 2 elements, we just need to swap
	// Check if we can optimize with B stack
	if len(*B) >= 2 && (*B)[0] < (*B)[1] {
		ss(A, B) // Optimize both stacks
	} else {
		sa(A, B)
	}
}

func three(A, B *stack) {
	if A.isSorted() {
		return
	}

	a0, a1, a2 := (*A)[0], (*A)[1], (*A)[2]

	// Simple and reliable case analysis
	if a0 > a1 && a1 < a2 && a0 < a2 {
		// 2 1 3 -> 1 2 3
		if len(*B) >= 2 && (*B)[0] < (*B)[1] {
			ss(A, B)
		} else {
			sa(A, B)
		}
	} else if a0 < a1 && a1 > a2 && a0 < a2 {
		// 1 3 2 -> 1 2 3
		if len(*B) >= 2 && (*B)[0] < (*B)[1] {
			rrr(A, B)
		} else {
			rra(A, B)
		}
	} else if a0 < a1 && a1 > a2 && a0 > a2 {
		// 2 3 1 -> 1 2 3
		if len(*B) >= 2 && (*B)[0] < (*B)[1] {
			rr(A, B)
		} else {
			ra(A, B)
		}
	} else if a0 > a1 && a1 > a2 {
		// 3 2 1 -> 1 2 3 (need sa + rra)
		sa(A, B)
		if !A.isSorted() {
			if len(*B) >= 2 && (*B)[0] < (*B)[1] {
				rrr(A, B)
			} else {
				rra(A, B)
			}
		}
	} else if a0 > a1 && a1 < a2 && a0 > a2 {
		// 3 1 2 -> 1 2 3 (need ra + sa)
		if len(*B) >= 2 && (*B)[0] < (*B)[1] {
			rr(A, B)
		} else {
			ra(A, B)
		}
		if !A.isSorted() {
			if len(*B) >= 2 && (*B)[0] < (*B)[1] {
				ss(A, B)
			} else {
				sa(A, B)
			}
		}
	}
}
