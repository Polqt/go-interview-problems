package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

func WalkOptimal(t *tree.Tree, ch chan int) {
	defer close(ch)
	walkRecursive(t, ch)
}

func walkRecursive(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}

	walkRecursive(t.Left, ch)
	ch <- t.Value
	walkRecursive(t.Right, ch)
}

func SameOptimal(t1, t2 *tree.Tree) bool {
	const bufferSize = 10
	ch1 := make(chan int, bufferSize)
	ch2 := make(chan int, bufferSize)

	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2

		if ok1 != ok2 {
			return false
		}

		if !ok1 {
			return true
		}

		if v1 != v2 {
			return false
		}
	}
}

func main() {
	ch := make(chan int)
	go WalkOptimal(tree.New(1), ch)
	fmt.Print("Values from tree.New(1): ")
	for i := 0; i < 10; i++ {
		fmt.Printf("%d ", <-ch)
	}
	fmt.Println() // Expected output: 1 2 3 4 5 6 7 8 9 10

	fmt.Println("SameOptimal(tree.New(1), tree.New(1)):", SameOptimal(tree.New(1), tree.New(1))) // Expected output: true
	fmt.Println("SameOptimal(tree.New(1), tree.New(2)):", SameOptimal(tree.New(1), tree.New(2))) // Expected output: false
}
