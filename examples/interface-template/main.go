package main

import "fmt"

func main() {
	var p PrinterStringInterface
	p = &PrinterString{}

	fmt.Printf("%s\n", p.Print("Hello world"))
}
