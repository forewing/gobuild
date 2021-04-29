package main

// int r69()
// {
//	    return 69;
// }
import "C"

import "fmt"

func main() {
	hello()
}

func hello() {
	fmt.Println(C.r69())
}
