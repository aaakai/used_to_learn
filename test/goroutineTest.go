package test

import (
	"fmt"
	"runtime"
	"time"
)

func cal() {
	for i := 0; i < 10; i++ {
		runtime.Gosched()
	}
}
func GoroutineTest() {

	runtime.GOMAXPROCS(1)
	currentTime := time.Now().UnixNano()
	go cal()
	cal()
	currentTime = time.Now().UnixNano() - currentTime
	fmt.Println("===", currentTime)
}
