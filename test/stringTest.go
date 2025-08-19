package test

import "fmt"

func StringTest() {
	a := "saaaa"
	b := []byte(a)
	for _, v := range a {
		if v == 'a' {
			fmt.Println(string(v))
		}
	}
	for _, v := range b {
		if v == 'a' {
			fmt.Println(v)
		}
	}
}
