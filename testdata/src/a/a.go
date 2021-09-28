package a

import "reflect"


func main() { // want main:"isPanic"
	f()
}

func f() { // want f:"isPanic"
	v := reflect.ValueOf(10)
	v.SetInt(20)
}