package whatitdisplays

import (
	"fmt"
	"math/rand"
)

/*
We will see deadloack exception
1. First solution is to add waitgroup and close channel when all gouroutins stop counting
2. Second solution is to add timeout context, so it will be canceled after delay or if all gouroutins stop counting (not valid)
3. The commented code is that I added
*/
func RunChAndGouroutins1() {
	c := make(chan int, 1000)
	// var wg sync.WaitGroup

	// wg.Add(10)
	for i := 0; i < 100; i++ {
		go func() {
			foo(c)
			// wg.Done()
		}()
	}

	// go func() {
	// 	wg.Wait()
	// 	close(c)
	// }()

	sum := 0
	for r := range c {
		sum += r
	}

	fmt.Println(sum)
}

func foo(c chan int) {
	r := rand.Int()
	for i := 0; i < r; i++ {
		c <- r
	}
}
