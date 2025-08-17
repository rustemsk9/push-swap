package main

func four_Five(A, B *stack) {
	tmp := len(*A) - 3
	for move := tmp; move > 0; {
		pivot, _ := A.getPivot()
		if (*A)[0] < pivot {
			pb(A, B)
			move--
		} else {
			ra(A, B) // check if rra is better?
		}
	}
	three(A, B)
	for move := tmp; move > 0; move-- {
		pa(A, B)
	}
}

func two(A, B *stack) {
	if A.isSorted() {
		return
	}
	sa(A, B)
}

func three(A, B *stack) {
	if A.isSorted() {
		return
	} else if (*A)[0] < (*A)[1] && (*A)[1] > (*A)[2] {
		rra(A, B) //1 3 2 -> 2 1 3
	} else if (*A)[0] > (*A)[1] && (*A)[1] > (*A)[2] {
		sa(A, B) // 3 2 1 -> 2 3 1
	}
	if (*A)[0] > (*A)[1] && (*A)[0] < (*A)[2] {
		if len(*B) == 2 && (*B)[0] < (*B)[1] {
			ss(A, B)
		} else {
			sa(A, B) // 2 1 3 -> 1 2 3	//check stackB? for ss?
		}
	} else if (*A)[0] < (*A)[1] {
		if len(*B) == 2 && (*B)[0] < (*B)[1] {
			rrr(A, B)
		} else {
			rra(A, B) // 2 1 3 -> 1 2 3	//check stackB? for ss?
		}
	} else {
		if len(*B) == 2 && (*B)[0] < (*B)[1] {
			rr(A, B)
		} else {
			ra(A, B) //Same
		}
	}
}
