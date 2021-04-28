package main

import "fmt"

var (
	Version = "dev"
	Hash    = "unknown"
)

func main() {
	fmt.Println("hello world")
	fmt.Println("version:", Version)
	fmt.Println("hash:", Hash)
}
