package main

import (
	"tp_go/dictionary"
)

func main() {
	//fmt.Println("hello world")

	m := make(map[string]string)

	dictionary.Add(m, "Map", "Maps are Go’s built-in associative data type (sometimes called hashes or dicts in other languages).")
	dictionary.Add(m, "Values", "Go has various value types including strings, integers, floats, booleans, etc. Here are a few basic examples.")
	dictionary.List(m)
	dictionary.Remove(m, "Map")
	dictionary.List(m)
	dictionary.Add(m, "Variables", "In Go, variables are explicitly declared and used by the compiler to e.g. check type-correctness of function calls.")
	dictionary.Get(m, "Variables")
	dictionary.Reset(m)
	dictionary.List(m)
}
