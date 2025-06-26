package test

import "fmt"

type Test struct {
	A int64 `json:"a,omitempty"`
	B int64 `json:"b"`
}

func AnyThingTest() {
	/*numGoroutines := runtime.NumGoroutine()
	fmt.Printf("Current number of Goroutines: %d\n", numGoroutines)

	// 查看当前的逻辑处理器 P 的数量
	numP := runtime.NumCPU()
	fmt.Printf("Current number of logical CPUs (P): %d\n", numP)*/

	//var arr [5]int = [5]int{1, 2, 3, 4, 5}
	a := make([]int, 5, 6)
	a[0] = 1
	a[1] = 2
	a[2] = 3
	a[3] = 4
	a[4] = 5
	b := a[:3]
	fmt.Println(a)
	fmt.Println(b)
	b = append(b, 1, 2, 3)
	fmt.Println(b)
	fmt.Println(a)
}
