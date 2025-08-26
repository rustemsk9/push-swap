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
	a := *A
	if len(a) != 3 {
		return
	}
	switch {
	case a[0] < a[1] && a[1] < a[2]:
		// Already sorted
		return
	case a[0] > a[1] && a[1] < a[2] && a[0] < a[2]:
		sa(A, B)
	case a[0] > a[1] && a[1] > a[2]:
		sa(A, B)
		rra(A, B)
	case a[0] > a[1] && a[1] < a[2] && a[0] > a[2]:
		ra(A, B)
	case a[0] < a[1] && a[1] > a[2] && a[0] < a[2]:
		sa(A, B)
		ra(A, B)
	case a[0] < a[1] && a[1] > a[2] && a[0] > a[2]:
		rra(A, B)
	}
}
