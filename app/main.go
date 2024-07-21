package main

import "fmt"

func app(s []int, vars ...int) (s2 []int) {
	s2 = s
	s2 = append(s2, vars...)
	return
}

func main() {
	s := []int{1, 2, 3}
	s1 := app(s, 34)
	fmt.Println(len(s), cap(s))
	fmt.Println(len(s1), cap(s1))
}
