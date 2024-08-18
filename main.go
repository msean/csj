package main

import (
	"fmt"
	"reflect"
)

func Inslice(val interface{}, slice interface{}) bool {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice && sliceValue.Kind() != reflect.Array {
		return val == slice
	}

	for i := 0; i < sliceValue.Len(); i++ {
		_val := sliceValue.Index(i).Interface()
		if val == _val {
			return true
		}
	}

	return false
}

func main() {
	v := 4
	vs := []int{1, 2, 3}
	fmt.Println(Inslice(v, vs))
}
