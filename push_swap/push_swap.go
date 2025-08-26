package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type stack struct {
	val             int
	currentPosition int
	finalIndex      int
	pushPrice       int
	aboveMedian     bool
	cheapest        bool
	targetNode      *stack
	next            *stack
	prev            *stack
}

var checker bool

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

func stackFromArgNums(stackA **stack, argv []string) {
	isDup := dupChecker()
	var smallest int
	checker = false
	for _, str := range argv {
		if str == "--checker" {
			checker = true
			continue
		}
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

			appendNode(stackA, num)
		}
	}
	if smallest < 0 {
		updateIfNeg(smallest, stackA)
	}
}

func updateIfNeg(smallest int, stackA **stack) {
	offset := -smallest
	for node := *stackA; node != nil; node = node.next {
		node.val += offset
	}
}

func appendNode(s **stack, val int) {
	if s == nil {
		return
	}

	node := &stack{
		val:  val,
		next: nil,
		prev: nil,
	}
	if *s == nil {
		*s = node
	} else {
		last := findLastNode(*s)
		last.next = node
		node.prev = last
	}
	// Optionally print the appended value for debugging
	// fmt.Println("Appended", node.val)
}

func findLastNode(s *stack) *stack {
	if s == nil {
		return nil
	}
	for s.next != nil {
		s = s.next
	}
	return s
}

type valueStack struct {
	A, B *stack
}
type sizeStack struct {
	A, B *stack
}

func Init_nodes(a, b *stack) {
	set_curr_position(a)
	set_curr_position(b)
	set_target_node(a, b)
	set_price(a, b)
	set_cheapest(b)
}

func set_price(a, b *stack) {
	lenA := stack_len(a)
	lenB := stack_len(b)
	for nodeB := b; nodeB != nil; nodeB = nodeB.next {
		nodeB.pushPrice = nodeB.currentPosition
		if !nodeB.aboveMedian {
			nodeB.pushPrice = lenB - nodeB.currentPosition
		}
		if nodeB.targetNode != nil {
			if nodeB.targetNode.aboveMedian {
				nodeB.pushPrice += nodeB.targetNode.currentPosition
			} else {
				nodeB.pushPrice += lenA - nodeB.targetNode.currentPosition
			}
		}
	}
}

func set_cheapest(b *stack) {
	if b == nil {
		return
	}
	bestMatchValue := int(^uint(0) >> 1) // Max int
	var bestMatchNode *stack

	for node := b; node != nil; node = node.next {
		node.cheapest = false // reset
		if node.pushPrice < bestMatchValue {
			bestMatchValue = node.pushPrice
			bestMatchNode = node
		}
	}
	if bestMatchNode != nil {
		bestMatchNode.cheapest = true
	}
}

func set_curr_position(s *stack) {
	if s == nil {
		return
	}
	// Find length
	length := 0
	for node := s; node != nil; node = node.next {
		length++
	}
	centerline := length / 2

	i := 0
	for node := s; node != nil; node = node.next {
		node.currentPosition = i
		if i <= centerline {
			node.aboveMedian = true
		} else {
			node.aboveMedian = false
		}
		i++
	}
}

func (s *stack) isEmpty() bool {
	return s == nil
}

func (s *stack) isSorted() bool {
	if s == nil || s.next == nil {
		return true
	}
	for node := s; node.next != nil; node = node.next {
		if node.val > node.next.val {
			return false
		}
	}
	return true
}

func (s *stack) peekIndex() int {
	if s == nil {
		return -1
	}
	return s.currentPosition
}

func find_smallest(s *stack) *stack {
	if s == nil {
		return nil
	}
	smallest := s
	for node := s; node != nil; node = node.next {
		if node.val < smallest.val {
			smallest = node
		}
	}
	return smallest
}

func set_target_node(a, b *stack) {
	for nodeB := b; nodeB != nil; nodeB = nodeB.next {
		var targetNode *stack
		bestMatchVal := int(^uint(0) >> 1) // Max int

		for nodeA := a; nodeA != nil; nodeA = nodeA.next {
			if nodeA.val > nodeB.val && nodeA.val < bestMatchVal {
				bestMatchVal = nodeA.val
				targetNode = nodeA
			}
		}

		if targetNode == nil {
			nodeB.targetNode = find_smallest(a)
		} else {
			nodeB.targetNode = targetNode
		}
	}
}

func tiny_sort(a **stack) {
	highestNode := find_highest(*a)
	if *a == highestNode {
		ra(a, false)
	} else if (*a).next == highestNode {
		rra(a, false)
	}
	if (*a).val > (*a).next.val {
		sa(a, false)
	}
}

func find_highest(s *stack) *stack {
	if s == nil {
		return nil
	}
	highest := s
	for node := s; node != nil; node = node.next {
		if node.val > highest.val {
			highest = node
		}
	}
	return highest
}

func finish_rotation(stack **stack, topNode *stack, stackName rune) {
	for *stack != topNode {
		if stackName == 'a' {
			if topNode.aboveMedian {
				ra(stack, false)
			} else {
				rra(stack, false)
			}
		} else if stackName == 'b' {
			if topNode.aboveMedian {
				rb(stack, false)
			} else {
				rrb(stack, false)
			}
		}
	}
}

func handle_five(a, b **stack) {
	for stack_len(*a) > 3 {
		Init_nodes(*a, *b)
		finish_rotation(a, find_smallest(*a), 'a')
		pb(b, a, false)
	}
}

// Helper function to get the length of the stack
func stack_len(s *stack) int {
	length := 0
	for node := s; node != nil; node = node.next {
		length++
	}
	return length
}

func return_cheapest(s *stack) *stack {
	if s == nil {
		return nil
	}
	for node := s; node != nil; node = node.next {
		if node.cheapest {
			return node
		}
	}
	return nil
}

func rotate_both(a, b **stack, cheapest *stack) { //
	for *b != cheapest && *a != cheapest.targetNode {
		rr(a, b, false)
	}
	set_curr_position(*a)
	set_curr_position(*b)
}

func reverse_rotate_both(a, b **stack, cheapest *stack) {
	for *b != cheapest && *a != cheapest.targetNode {
		rrr(a, b, false)
	}
	set_curr_position(*a)
	set_curr_position(*b)
}

func Move_nodes(a, b **stack) {
	cheapest := return_cheapest(*b)
	if cheapest == nil {
		return
	}
	if cheapest.aboveMedian && cheapest.targetNode.aboveMedian {
		rotate_both(a, b, cheapest)
	} else if !cheapest.aboveMedian && !cheapest.targetNode.aboveMedian {
		reverse_rotate_both(a, b, cheapest)
	}
	finish_rotation(b, cheapest, 'b')
	finish_rotation(a, cheapest.targetNode, 'a')
	pa(a, b, false)
}

func push_sw(a, b **stack) {
	lenA := stack_len(*a)
	if lenA == 5 {
		handle_five(a, b)
	} else {
		for i := 0; i < lenA-3; i++ {
			pb(b, a, false)
		}
	}

	tiny_sort(a)

	for *b != nil {
		Init_nodes(*a, *b)
		Move_nodes(a, b)
	}

	set_curr_position(*a)
	smallest := find_smallest(*a)
	if smallest != nil && smallest.aboveMedian {
		for *a != smallest {
			ra(a, checker)
		}
	} else {
		for *a != smallest {
			rra(a, checker)
		}
	}
}

// SECTION move

func Push(dest **stack, src **stack) {
	if *src == nil {
		return
	}
	nodeToPush := *src
	*src = (*src).next
	if *src != nil {
		(*src).prev = nil
	}
	nodeToPush.prev = nil
	if *dest == nil {
		*dest = nodeToPush
		nodeToPush.next = nil
	} else {
		nodeToPush.next = *dest
		(*dest).prev = nodeToPush
		*dest = nodeToPush
	}
}

func pa(a **stack, b **stack, checker bool) {
	Push(a, b)
	if !checker {
		fmt.Print("pa\n")
	}
}

func pb(b **stack, a **stack, checker bool) {
	Push(b, a)
	if !checker {
		fmt.Print("pb\n")
	}
}

// SECTION reverse_rotate
// reverseRotate moves the last element of the stack to the front.
func reverseRotate(head **stack) {
	if head == nil || *head == nil || (*head).next == nil {
		return
	}
	last := findLastNode(*head)
	if last == nil || last.prev == nil {
		return
	}
	// Detach last node
	last.prev.next = nil
	last.prev = nil
	// Move last to front
	last.next = *head
	(*head).prev = last
	*head = last
}

// rra reverses stack a.
func rra(a **stack, checker bool) {
	reverseRotate(a)
	if !checker {
		fmt.Println("rra")
	}
}

// rrb reverses stack b.
func rrb(b **stack, checker bool) {
	reverseRotate(b)
	if !checker {
		fmt.Println("rrb")
	}
}

// rrr reverses both stacks a and b.
func rrr(a, b **stack, checker bool) {
	reverseRotate(a)
	reverseRotate(b)
	if !checker {
		fmt.Println("rrr")
	}
}

// SECTION rotate
// rotate moves the first element of the stack to the end.
func rotate(head **stack) {
	if head == nil || *head == nil || (*head).next == nil {
		return
	}
	first := *head
	last := findLastNode(*head)

	*head = first.next
	(*head).prev = nil

	last.next = first
	first.prev = last
	first.next = nil
}

// ra rotates stack a.
func ra(a **stack, checker bool) {
	rotate(a)
	if !checker {
		fmt.Println("ra")
	}
}

// rb rotates stack b.
func rb(b **stack, checker bool) {
	rotate(b)
	if !checker {
		fmt.Println("rb")
	}
}

// rr rotates both stacks a and b.
func rr(a, b **stack, checker bool) {
	rotate(a)
	rotate(b)
	if !checker {
		fmt.Println("rr")
	}
}

// SECTION swap
// swap swaps the first two elements of the stack.
func swap(head **stack) {
	if head == nil || *head == nil || (*head).next == nil {
		return
	}
	first := *head
	second := first.next

	// Adjust pointers to swap first and second nodes
	first.next = second.next
	if second.next != nil {
		second.next.prev = first
	}
	second.prev = nil
	second.next = first
	first.prev = second

	*head = second
}

// sa swaps the first two elements of stack a.
func sa(a **stack, checker bool) {
	swap(a)
	if !checker {
		fmt.Println("sa")
	}
}

// sb swaps the first two elements of stack b.
func sb(b **stack, checker bool) {
	swap(b)
	if !checker {
		fmt.Println("sb")
	}
}

// ss swaps the first two elements of both stacks a and b.
func ss(a, b **stack, checker bool) {
	swap(a)
	swap(b)
	if !checker {
		fmt.Println("ss")
	}
}

func printNode(s *stack) {
	for node := s; node != nil; node = node.next {
		fmt.Printf("val: %d, currPos: %d, finalIdx: %d, pushPrice: %d, aboveMedian: %v, cheapest: %v\n",
			node.val, node.currentPosition, node.finalIndex, node.pushPrice, node.aboveMedian, node.cheapest)
	}
}

func main() {
	if len(os.Args) == 1 {
		os.Exit(2)
	}
	// fmt.Println("Initializing stack A from arguments...")
	var stackA *stack
	stackFromArgNums(&stackA, os.Args[1:])
	var stackB *stack
	// printNode(stackA)
	if !stackA.isSorted() {
		if stack_len(stackA) == 2 {
			sa(&stackA, checker)
		} else if stack_len(stackA) == 3 {
			tiny_sort(&stackA)
		} else {
			push_sw(&stackA, &stackB)
			// push_swap(&stackA, &stackB) --- IGNORE ---
		}
	}
	// printNode(stackA)
}
