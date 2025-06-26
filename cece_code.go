package main

import "strconv"

/*
1：编写一个函数，使用两个 goroutine 交替打印数字和字母，如：1A2B3C...26Z
*/
func Code() {
	chanNum := make(chan struct{}, 1)
	chanStr := make(chan struct{}, 1)
	chanStr <- struct{}{}
	endChanNum := make(chan struct{})
	go func() {
		for i := 0; i < 26; i++ {
			<-chanStr
			print(i + 1)
			chanNum <- struct{}{}
		}
	}()
	go func() {
		for i := 0; i < 26; i++ {
			<-chanNum
			print(string('A' + i))
			chanStr <- struct{}{}
			if i == 25 {
				endChanNum <- struct{}{}
			}
		}
	}()
	<-endChanNum
	close(chanNum)
	close(chanStr)
}

/*
2:给你一个整数x，若x 是一个回文整数，返回 true ；否则，返回 false 。
eg: 121是回文数，123不是回文。
*/

func Code1(x int) bool {
	var str = strconv.Itoa(x)
	if len(str) < 2 {
		return true
	}
	i := 0
	j := len(str) - 1
	for i < j {
		if str[i] != str[j] {
			return false
		}
		i++
		j--
	}
	return true
}
