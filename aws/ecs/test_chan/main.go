package main

import (
	"fmt"
)

func main() {
	c := make(chan int, 1)
	//c <- 1
	close(c)
	i, ok := <-c
	fmt.Println(ok, i)
	i, ok = <-c
	fmt.Println(ok, i)
}
