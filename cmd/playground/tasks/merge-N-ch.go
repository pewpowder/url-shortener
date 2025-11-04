package tasks

import (
	"fmt"
	"sync"
)

// merge N channels
func MergeNCh() {
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 20)
	ch3 := make(chan int, 15)

fillLoop:
	for i := range 50 {
		select {
		case ch1 <- i:
		case ch2 <- i:
		case ch3 <- i:
		default:
			fmt.Println("All channels filled")
			break fillLoop
		}
	}

	close(ch1)
	close(ch2)
	close(ch3)

	mergedCh := merge(ch1, ch2, ch3)

	for val := range mergedCh {
		fmt.Println(val)
	}

	fmt.Println("Finished")
}

func merge[T any](chans ...chan T) chan T {
	out := make(chan T, 10) // buffered ch instead of unbuffered make merge/reading a little bit fater
	var wg sync.WaitGroup

	for _, ch := range chans {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
