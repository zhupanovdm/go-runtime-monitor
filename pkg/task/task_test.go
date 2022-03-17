package task

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func ExampleTask_With() {
	ctx := context.TODO()

	middleware := func(task Task) Task {
		return func(ctx context.Context) {
			fmt.Println("before task")
			task(ctx)
			fmt.Println("after task")
		}
	}

	t := Task(func(ctx context.Context) {
		fmt.Println("task execution")
	}).With(middleware)

	// task will be executed within our middleware
	t(ctx)

	fmt.Println("completed")

	// Output:
	// before task
	// task execution
	// after task
	// completed
}

func ExampleCompletionWait() {
	ctx := context.TODO()

	var wg sync.WaitGroup
	defer wg.Wait()

	t1 := Task(func(ctx context.Context) {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("task 1")
	}).With(CompletionWait(&wg))

	t2 := Task(func(ctx context.Context) {
		time.Sleep(200 * time.Millisecond)
		fmt.Println("task 2")
	}).With(CompletionWait(&wg))

	// the program will wait for completion of both tasks through wg
	go t1(ctx)
	go t2(ctx)

	fmt.Println("completed")

	// Output:
	// completed
	// task 1
	// task 2
}

func ExamplePeriodicRun() {
	ctx, cancel := context.WithTimeout(context.TODO(), 500*time.Millisecond)
	defer cancel()

	counter := 0

	periodic := Task(func(ctx context.Context) { counter++ }).With(PeriodicRun(50 * time.Millisecond))

	// the task will be executed every 50ms until context is cancelled
	go periodic(ctx)

	time.Sleep(550 * time.Millisecond)

	fmt.Println(counter > 8)

	// Output:
	// true
}

func ExamplePeriodicRun_withWait() {
	ctx, cancel := context.WithTimeout(context.TODO(), 500*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup

	counter := 0

	periodic := Task(func(ctx context.Context) { counter++ }).
		With(PeriodicRun(50*time.Millisecond), CompletionWait(&wg))

	// the task will be executed every 50ms until context is cancelled
	go periodic(ctx)

	// will wait until context is cancelled
	wg.Wait()
	fmt.Println(counter > 8)

	// Output:
	// true
}
