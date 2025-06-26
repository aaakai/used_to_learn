package test

import (
	"fmt"
	"time"
)

func ChannelTest2() {
	oddChan := make(chan struct{})  // 奇数通道
	evenChan := make(chan struct{}) // 偶数通道
	done := make(chan struct{})     // 完成通道

	// 奇数 goroutine
	go func() {
		for i := 1; i <= 100; i += 2 {
			<-oddChan // 等待信号
			fmt.Println("奇数:", i)
			evenChan <- struct{}{} // 通知偶数 goroutine
		}
	}()

	// 偶数 goroutine
	go func() {
		for i := 2; i <= 100; i += 2 {
			<-evenChan // 等待信号
			fmt.Println("偶数:", i)
			if i < 100 { // 防止最后一次发送导致 deadlock
				oddChan <- struct{}{}
			}
		}
		close(done) // 通知主 goroutine 完成
	}()

	// 启动第一个 goroutine
	oddChan <- struct{}{}

	// 等待完成
	<-done
}

func ChannelTest1() {
	one := make(chan struct{})
	two := make(chan struct{})
	done := make(chan struct{})
	go func() {
		for i := 1; i <= 100; i += 2 {
			<-one
			fmt.Println("奇数：", i)
			two <- struct{}{}
		}
	}()
	go func() {
		for i := 2; i <= 100; i += 2 {
			<-two
			fmt.Println("偶数：", i)
			if i < 100 {
				one <- struct{}{}
			}
		}
		done <- struct{}{}
	}()
	one <- struct{}{}
	<-done
	return
}

func ChannelTest() {
	done := make(chan struct{}, 2)
	closeDone := make(chan struct{})
	ch := make(chan int, 1)
	ch <- 1
	go printMod(0, ch, done)
	go printMod(1, ch, done)
	go func() {
		time.Sleep(time.Second * 5)
		closeDone <- struct{}{}
	}()
	for {
		select {
		case <-closeDone:
			fmt.Println("close")
			return
		case <-done:
			fmt.Println("done")
		}
	}
}

func printMod(mod int, ch chan int, done chan<- struct{}) {
	for i := 0; i < 100; i++ {
		if i%2 == mod {
			fmt.Printf("%dchannel打印\n ", mod)
			fmt.Printf("%d 的结果\n", i)
		}
	}
	done <- struct{}{}
}

func ChannelTest3() {
	oneChan := make(chan struct{}) // 奇数通道
	twoChan := make(chan struct{}) // 偶数通道
	threeChan := make(chan struct{})
	done := make(chan struct{}) // 完成通道

	// 奇数 goroutine
	go func() {
		for i := 1; i <= 97; i += 3 {
			<-oneChan // 等待信号
			fmt.Println(i)
			twoChan <- struct{}{} // 通知偶数 goroutine
		}
	}()

	// 偶数 goroutine
	go func() {
		for i := 2; i <= 98; i += 3 {
			<-twoChan // 等待信号
			fmt.Println(i)
			threeChan <- struct{}{}
		}
	}()

	// 偶数 goroutine
	go func() {
		for i := 3; i <= 99; i += 3 {
			<-threeChan // 等待信号
			fmt.Println(i)
			if i < 99 {
				oneChan <- struct{}{}
			}
		}
		close(done)
	}()
	// 启动第一个 goroutine
	oneChan <- struct{}{}

	// 等待完成
	<-done
}
