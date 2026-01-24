package main

import (
	"unsafe"

	"github.com/pewpowder/url-shortener/cmd/playground/tasks"
)

func IsLittleEndian() bool {
	var value int16 = 0b001
	ptr := (unsafe.Pointer(&value))
	oneBit := (*int8)(ptr)

	return *oneBit == 1
}

func main() {
	// ch1 := make(chan int)
	// ch2 := make(chan int)

	// go Producer(ch1)
	// go Producer(ch2)

	// for {
	// 	// select {
	// 	// case value := <- ch1:
	// 	// 	fmt.Println("ch1 1", value);
	// 	// 	return
	// 	// default:
	// 	// }

	// 	select {
	// 	case value := <- ch1:
	// 		time.Sleep(time.Second *2)
	// 		fmt.Println("ch1 2", value);
	// 		return
	// 	case value := <- ch2:
	// 		fmt.Println("ch2 2", value)
	// 	}
	// }

	// Tasks
	// tasks.RunInnocentsArray()
	// tasks.RunUnixTee()
	// tasks.RunPubSub()
	// tasks.MergeNCh()
	// tasks.ProcessParallel()
	tasks.HttpRequests()

	// What does it display?
	// whatitdisplays.RunChAndGouroutins1()
}