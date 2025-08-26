package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = stdout
	buf.ReadFrom(r)
	return strings.TrimSpace(buf.String())
}

func Test1(t *testing.T) {
	input := "2 1 3 6 5 8"
	expectedOutput := "pb\npb\nra\nsa\nrrr\npa\npa"
	stackA := stackFromArgNums([]string{input})
	stackB := stack{}

	output := captureOutput(func() {
		solve(&stackA, &stackB)
	})

	if output != expectedOutput {
		t.Errorf("Test1 failed: expected:\n%s\n\tgot:\n\t%s", expectedOutput, output)
	}
}

func TestAlreadySorted(t *testing.T) {
	input := "1 2 3 4 5"
	expectedOutput := ""
	stackA := stackFromArgNums([]string{input})
	stackB := stack{}

	output := captureOutput(func() {
		solve(&stackA, &stackB)
	})

	if output != expectedOutput {
		t.Errorf("TestAlreadySorted failed: expected:\n%s\ngot:\n%s", expectedOutput, output)
	}
}

func TestReverseOrder(t *testing.T) {
	input := "5 4 3 2 1"
	expectedOutput := "pb\npb\nsa\npa\npa\nra\nra\nra\nra"
	stackA := stackFromArgNums([]string{input})
	stackB := stack{}

	output := captureOutput(func() {
		solve(&stackA, &stackB)
	})

	if output != expectedOutput {
		t.Errorf("TestReverseOrder failed: expected:\n%s\ngot:\n%s", expectedOutput, output)
	}
}

func TestSingleElement(t *testing.T) {
	input := "42"
	expectedOutput := ""
	stackA := stackFromArgNums([]string{input})
	stackB := stack{}

	output := captureOutput(func() {
		solve(&stackA, &stackB)
	})

	if output != expectedOutput {
		t.Errorf("TestSingleElement failed: expected:\n%s\ngot:\n%s", expectedOutput, output)
	}
}

// Test implementation goes here
// }
