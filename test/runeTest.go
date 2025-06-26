package test

import "fmt"

func RuneTest() {
	var s = "golang你好"
	charPrint(s)
	bytePrint(s)
}

func charPrint(s string) {
	a := []rune(s)
	for i := 0; i < len(a); i++ {
		fmt.Printf("%c ", a[i])
	}
}

func bytePrint(s string) {
	for i := 0; i < len(s); i++ {
		fmt.Printf("%c ", s[i])
	}
}
