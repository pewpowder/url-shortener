package tasks

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// реализовать функцию processParallel
// прокинуть контекст

func processData(v int) int {
	// select {
	// case <-ctx.Done():
	// 	// cancel computation
	// case <-time.After(time.Duration(rand.Intn(10)) * time.Second):
	// 	// return v * 2
	// }

	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	return v * 2
}

func ProcessParallel() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	in := make(chan int)
	out := make(chan int)

	go func() {
		defer close(in)

	breakLoop:
		for i := range 10 {
			select {
			case in <- i:
			case <-ctx.Done():
				fmt.Println("Context canceled MAIN")
				break breakLoop
			}
		}
	}()

	start := time.Now()
	processParallel(ctx, in, out, 5)

	for v := range out {
		fmt.Println("v =", v)
	}

	fmt.Println("main duration:", time.Since(start))
}

func processParallel(ctx context.Context, in, out chan int, numWorkers int) {
	// https://www.youtube.com/watch?v=wZCLVt_5-4c
	for range numWorkers {
		go func() {
			defer close(out)

		breakLoop:
			for {
				select {
				case val, ok := <-in:
					if !ok {
						return
					}
					out <- processData(val)
				case <-ctx.Done():
					fmt.Println("Context canceled PARALLEL_PROCESS")
					break breakLoop
				}
			}
		}()
	}

	// if I can cahnge processData (commented code)
	// size := make(chan struct{}, numWorkers)
	// for {
	// 	select {
	// 	case size <- struct{}{}:
	// 		go func() {
	// 			select {
	// 			case val := <-in:
	// 				processData(val)
	// 			case <-ctx.Done():
	// 				fmt.Println("Context canceled PARALLEL_PROCESS")
	// 				break breakLoop
	// 			}
	// 		}()
	// 	}
	// }
}
