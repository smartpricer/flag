package main

import (
	"fmt"

	"github.com/smartpricer/flag"
)

func main() {
	var (
		config   string
		length   float64
		age      int
		name     string
		lastName string
		female   bool
	)

	flag.StringVar(&config, "config", "", "help message")
	flag.StringVar(&name, "name", "", "help message")
	flag.StringVar(&lastName, "last-name", "", "help message")

	flag.IntVar(&age, "age", 0, "help message")
	flag.Float64Var(&length, "length", 0, "help message")
	flag.BoolVar(&female, "female", false, "help message")

	flag.Parse()

	fmt.Println("length:", length)
	fmt.Println("age:", age)
	fmt.Println("name:", name)
	fmt.Println("lastName:", lastName)
	fmt.Println("female:", female)
}
