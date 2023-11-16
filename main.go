package main

import "fmt"

func main() {
	//fmt.Println("hello world")

	m := make(map[string]string)

	Add(m, "Map", "Maps are Go’s built-in associative data type (sometimes called hashes or dicts in other languages).")
	Add(m, "Values", "Go has various value types including strings, integers, floats, booleans, etc. Here are a few basic examples.")
	List(m)
	Remove(m, "Map")
	List(m)
	Add(m, "Variables", "In Go, variables are explicitly declared and used by the compiler to e.g. check type-correctness of function calls.")
	Get(m, "Variables")
	Reset(m)
	List(m)
}

func Add(m map[string]string, key string, value string) {
	m[key] = value
}

func Get(m map[string]string, key string) {
	fmt.Println("_______________________________")
	if value, ok := m[key]; ok {
		fmt.Println(key+":", value)
	} else {
		fmt.Println("Key not found:", key)
	}
}

func List(m map[string]string) {
	fmt.Println("_______________________________")
	if len(m) != 0 {
		for key, value := range m {
			fmt.Printf("%s: %s\n", key, value)
		}
	} else {
		fmt.Println("List Empty")
	}
}

func Remove(m map[string]string, value string) {
	delete(m, value)
}

func Reset(m map[string]string) {
	clear(m)
}
