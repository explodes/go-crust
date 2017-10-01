package main

import "fmt"

func main() {

	foo := []string{"a", "b", "c"}

	fmt.Println(foo[:0])
	fmt.Println(foo[:1])
	fmt.Println(foo[:2])
	fmt.Println(foo[:3])
	fmt.Println(foo[:4])
}
