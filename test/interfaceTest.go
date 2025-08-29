package test

import (
	"fmt"
	"reflect"
)

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

func Interface1() {
	v := reflect.ValueOf(1)
	i := v.Interface().(int)
	fmt.Println(reflect.TypeOf(i))
	fmt.Println(reflect.ValueOf(i))
}
