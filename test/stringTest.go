package test

import "fmt"

func StringTest() {
	a := "saaaa"
	b := a
	fmt.Println(a)
	fmt.Println(b)
	b = "aaaab"
	fmt.Println(a)
	fmt.Println(b)
}
