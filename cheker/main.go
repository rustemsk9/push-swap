package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type command struct {
	name string
	next *command
	prev *command
}

var reverseFunctions = map[string]string{
	"sa": "sa", "sb": "sb", "ss": "ss",
	"pa": "pb", "pb": "pa",
	"ra": "rra", "rb": "rrb", "rr": "rrr",
	"rra": "ra", "rrb": "rb", "rrr": "rr",
}

func addLink(head **command, new *command) {
	if *head == nil {
		*head = new
	} else {
		tmp := *head
		for tmp.next != nil {
			tmp = tmp.next
		}
		tmp.next = new
		new.prev = tmp
	}
}

func scanner(functions map[string]func(A, B *stack)) *command {
	var head *command
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if _, ok := functions[input]; !ok {
			fmt.Println("Error")
			os.Exit(1)
		}
		addLink(&head, &command{input, nil, nil})
	}
	return head
}

func stackGUI(stack *stack, color ui.Color) *widgets.Sparkline {
	tmp := widgets.NewSparkline()
	tmp.Data = *stack
	tmp.LineColor = color
	tmp.TitleStyle.Fg = ui.ColorWhite
	return tmp
}

func visualizer(curent *command, A, B *stack, functions map[string]func(A, B *stack)) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to init termui: %v", err)
	}
	defer ui.Close()

	sA := stackGUI(A, ui.ColorRed)
	sB := stackGUI(B, ui.ColorBlue)

	sAg := widgets.NewSparklineGroup(sA)
	sAg.Title = "Stack A"
	sBg := widgets.NewSparklineGroup(sB)
	sBg.Title = "Stack B"

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0/2, ui.NewCol(1.0, sAg)),
		ui.NewRow(1.0/2, ui.NewCol(1.0, sBg)))
	ui.Render(grid)
	uiEvents := ui.PollEvents()
	// ticker := time.NewTicker(time.Millisecond).C

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "g":
				if curent != nil {
					functions[curent.name](A, B)
					curent = curent.next
					sAg.Sparklines[0].Data = *A
					sBg.Sparklines[0].Data = *B
					ui.Render(grid)
				}
			case "f":
				if curent != nil && curent.prev != nil {
					curent = curent.prev
					// Apply reverse operation
					reverse := reverseFunctions[curent.name]
					functions[reverse](A, B)
					sAg.Sparklines[0].Data = *A
					sBg.Sparklines[0].Data = *B
					ui.Render(grid)
				}
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			}
			// case <-ticker:
			// if curent != nil {
			// 	functions[curent.name](A, B)
			// 	curent = curent.next
			// 	sAg.Sparklines[0].Data = *A
			// 	sBg.Sparklines[0].Data = *B
			// 	ui.Render(grid)
			// }
		}
	}
}

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

func updateIfNeg(smallest int, stackA *stack) {
	smallest *= -1
	for i, x := range *stackA {
		(*stackA)[i] = x + float64(smallest)
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
	if smallest < 0 {
		updateIfNeg(smallest, &stackA)
	}
	return stackA
}

func main() {
	functions := map[string]func(A, B *stack){
		"sa": sa, "sb": sb, "ss": ss, "pa": pa, "pb": pb, "ra": ra, "rb": rb, "rr": rr, "rra": rra, "rrb": rrb, "rrr": rrr,
	}

	if len(os.Args) == 1 {
		os.Exit(2)
	}
	stackA := stackFromArgNums(os.Args[1:])
	stackB := stack{}
	commandList := scanner(functions)
	visualizer(commandList, &stackA, &stackB, functions)
	fmt.Println(stackA)
	check(&stackA, &stackB)
}

/*
	Todo:
		* check if numbers are larger than int? custom Atoi? Atof?
*/
