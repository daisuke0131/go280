package a

import (
	"reflect"
)

func main() { // want main:"isPanic"
	f() // want "panic"
	g()
}

func f() { // want f:"isPanic"
	v := reflect.ValueOf(10)
	v.SetInt(20) // want "panic"
}

func g() {
	defer func() {
		if p := recover(); p != nil {
			println("recover")
		}
	}()
	v := reflect.ValueOf(10)
	v.SetInt(20)
}
