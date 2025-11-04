package tasks

import (
	"context"
	"fmt"
	"time"
)

/*
Задача: реализовать аналог unix tee в виде функции, используя базовые конструкции языка
*/

func RunUnixTee() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	in := gen(ctx, 10)

	out1, out2 := tee(ctx, in)

	println(<-out1) // 1
	println(<-out2) // 1

	println(<-out2) // 2
	println(<-out1) // 2

	println(<-out2) // 3*
	println(<-out2) // 0
}

// gen возвращает канал, из которого поступают значения от 1 до n.
//
// Канал закрывается, если все значения будут вычитаны или отменится ctx.
func gen(ctx context.Context, n int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for i := range n {
			select {
			case <-ctx.Done():
				fmt.Println("Context canceled gen")
				return
			case out <- i + 1:
			}
		}
	}()

	return out
}

// tee разделяет канал in на каналы out1 и out2 и возвращает их.

// Когда из канала in поступает новое значение, tee дублирует это значение
// и отправляет его в каналы out1 и out2. Если это значение не вычитано
// из каналов out1 и out2, то tee не отправляет следующее значение из канала in
// в каналы out1 и out2 до тех пор, пока предыдущее не будет вычитано.

// Порядок чтения каналов out1 и out2 не определен.

// Каналы out1 и out2 закрываются, если канал in будет закрыт или отменится ctx.
// teeWithPointers uses the insight that we can spawn a goroutine per consumer
// and use channels themselves for synchronization (no package imports needed).
//
// The trick: send on channels is blocking until the receiver is ready.
// We exploit this by spawning independent goroutines for each output channel.
func tee(ctx context.Context, in <-chan int) (<-chan int, <-chan int) {
	out1 := make(chan int)
	out2 := make(chan int)
	done1 := make(chan struct{})
	done2 := make(chan struct{})

	go func() {
		defer func() {
			close(out1)
			close(out2)
			close(done1)
			close(done2)
		}()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Context canceled tee")
				return
			case v, ok := <-in:
				if !ok {
					return
				}

				// Launch two goroutines to send to each output
				// They will block until the consumer reads
				go func() {
					select {
					case <-ctx.Done():
					case out1 <- v:
					}
					done1 <- struct{}{}
				}()

				go func() {
					select {
					case <-ctx.Done():
					case out2 <- v:
					}
					done2 <- struct{}{}
				}()

				// Wait for both goroutines to complete before next iteration
				<-done1
				<-done2
			}
		}
	}()

	return out1, out2
}
