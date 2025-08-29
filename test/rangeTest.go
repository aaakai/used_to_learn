package test

import "fmt"

func RangeTest() {
	arr := []int{1, 2, 3}
	newArr := []*int{}
	//v是临时变量 每次迭代重新赋值 但是地址不变
	for _, v := range arr {
		newArr = append(newArr, &v)
		fmt.Println(newArr)
		for _, v := range newArr {
			fmt.Println(*v)
		}
	}
	newArr1 := []*int{}

	for i, _ := range arr {
		newArr1 = append(newArr1, &arr[i])
		fmt.Println(newArr1)
		for _, v := range newArr1 {
			fmt.Println(*v)
		}
	}
}
