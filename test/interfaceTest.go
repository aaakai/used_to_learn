package test

import "fmt"

type Duck interface {
	Walk()
}

type Cat struct{}

func (c Cat) Walk() {
	fmt.Println("cat.walk")
}

type CatPointer struct{}

func (c *CatPointer) Walk() {
}

func InterfaceT() {
	var duck Duck = &Cat{}
	switch duck.(type) {
	case *Cat:
		cat := duck.(*Cat)
		cat.Walk()
	}

	var intDuck interface{} = &Cat{}
	switch intDuck.(type) {
	case *Cat:
		cat := duck.(*Cat)
		cat.Walk()
	}
}
