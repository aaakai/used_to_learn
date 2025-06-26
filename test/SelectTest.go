package test

import (
	"fmt"
	"time"
)

func SelectTest() {
	chan1 := make(chan int)
	chan2 := make(chan int)
	chan3 := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second * 1)
			chan1 <- 1
		}
	}()
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second * 1)
			chan2 <- 2
		}
	}()

	go func() {
		time.Sleep(time.Second * 11)
		chan3 <- 3
	}()

	for {
		select {
		case val := <-chan1:
			fmt.Println("chan1", val)
		case val := <-chan2:
			fmt.Println("chan2", val)
		case val := <-chan3:
			fmt.Println("chan3", val)
			return
			//default:
			//	fmt.Println("default")
			//	return
		}
	}
}
