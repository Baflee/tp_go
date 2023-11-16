package main

import (
	"tp_go/dictionary"
)

func main() {
	//fmt.Println("hello world")

	const filePath = "dictionary.txt"
	dictionary.Reset(filePath)
	dictionary.Add(filePath, "Map", "Maps are Go’s built-in associative data type (sometimes called hashes or dicts in other languages).")
	dictionary.Add(filePath, "Values", "Go has various value types including strings, integers, floats, booleans, etc. Here are a few basic examples.")
	dictionary.List(filePath)
	dictionary.Remove(filePath, "Map")
	dictionary.List(filePath)
	dictionary.Add(filePath, "Variables", "In Go, variables are explicitly declared and used by the compiler to e.g. check type-correctness of function calls.")
	dictionary.Get(filePath, "Variables")
	dictionary.List(filePath)
}
