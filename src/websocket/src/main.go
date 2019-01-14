// auth: kunlun
// date: 2019-01-11
// description:
package main

import (
	"fmt"
	"time"
)

func main() {
	msg := "starting main"
	fmt.Println(msg)
	bus := make(chan int)
	msg = "starting a gofunc"
	go counting(bus)
	for count := range bus {
		fmt.Println("count:", count)
	}
}

func counting(c chan<- int) {
	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)
		c <- i
	}
	close(c)
}
