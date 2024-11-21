package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	parent := context.Background()

	// Create two child contexts with their own cancel functions,
	// and create two further grandchildren from the first child
	child1, cancel1 := context.WithCancel(parent)
	grandchild1, cancelg1 := context.WithCancel(child1)
	grandchild2, cancelg2 := context.WithCancel(child1)
	child2, cancel2 := context.WithCancel(parent)

	defer cancelg1()
	defer cancelg2()
	defer cancel2()

	// Simulate canceling child1
	cancel1()

	// Check states
	fmt.Println("child1:", child1.Err())       // Outputs: context canceled
	fmt.Println("child1a:", grandchild1.Err()) // Outputs: context canceled
	fmt.Println("child1b:", grandchild2.Err()) // Outputs: context canceled
	fmt.Println("child2:", child2.Err())       // Outputs: <nil> (still active)
	fmt.Println("parent:", parent.Err())       // Outputs: <nil> (still active)

	time.Sleep(5 * time.Second)
}
