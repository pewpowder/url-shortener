package tasks

import (
	"context"
	"time"
)

func UnixTee() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	in := gen(ctx, 10)
	out1, out2 := tee(ctx, in)
	
	println(<-out1) // 1
	println(<-out2) // 1

	println(<-out2) // 2
	println(<-out1) // 2

	println(<-out2) // 3
	println(<-out2) // 0
	// Program exited.
}

// gen возвращает канал, из которого поступает значения от 1 до n.
//
// Канал закрывается, если все значения будут вычитаны или отменится ctx.
func gen(ctx context.Context, n int) <-chan int {
	out := make(chan int)
	
	go func() {
		defer close(out)
		for val := range n {
			select {
			case out <- (val + 1):
			case <-ctx.Done():
				return;
			}
		}
	}()
	
	return out
}

// tee разделяет канал in на каналы out1 и out2 и возвращает их.
//
// Когда из канала in поступает новое значение, оно дублируется
// и отправляется в каналы out1 и out2. Если это значение не вычитано
// из каналов out1 и out2, то tee не отправляет следующее значение из канала in
// в каналы out1 и out2 до тех пор, пока предыдущее не будет вычитано.
//
// Порядок чтения каналов out1 и out2 неопределен.
//
// Каналы out1 и out2 закрываются, если канал in будет закрыт или отменится ctx.
func tee(ctx context.Context, in <-chan int) (out1, out2 <-chan int) {
	ch1, ch2 := make(chan int), make(chan int)
	
	go func() {
		defer close(ch1)
		defer close(ch2)
		
		for {
			select {
			case <-ctx.Done():
				return
			case val := <-in:
				select {
				case <-ctx.Done():
					return
				case ch1 <- val:
					select {
					case <-ctx.Done():
					case ch2 <- val:
					}
				case ch2 <- val:
					select {
					case <-ctx.Done():
					case ch1 <- val:
					}
				}
			}
		}
	}()
	
	out1, out2 = ch1, ch2
	
	return
}