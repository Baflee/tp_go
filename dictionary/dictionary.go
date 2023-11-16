package dictionary

import "fmt"

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
