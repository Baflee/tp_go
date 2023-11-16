package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tp_go/dictionary"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	const filePath = "dictionary.txt"

	for {
		fmt.Println("__________________________________________________________")
		fmt.Println("Choisissez une action: add, define, remove, list, ou quit.")
		action, _ := reader.ReadString('\n')
		action = strings.TrimSpace(action)

		switch action {
		case "add":
			actionAdd(filePath, reader)
		case "define":
			actionDefine(filePath, reader)
		case "remove":
			actionRemove(filePath, reader)
		case "list":
			actionList(filePath)
		case "quit":
			return
		default:
			fmt.Println("Action non reconnue")
		}
	}
}

func actionAdd(filePath string, reader *bufio.Reader) {
	fmt.Print("Entrez le mot : ")
	key, _ := reader.ReadString('\n')
	key = strings.TrimSpace(key)

	fmt.Print("Entrez la définition : ")
	value, _ := reader.ReadString('\n')
	value = strings.TrimSpace(value)

	dictionary.Add(filePath, key, value)
	fmt.Println("Mot ajouté avec succès.")
}

func actionDefine(filePath string, reader *bufio.Reader) {
	fmt.Print("Entrez le mot à définir : ")
	key, _ := reader.ReadString('\n')
	key = strings.TrimSpace(key)

	result, err := dictionary.Get(filePath, key)
	if err != nil {
		fmt.Println("Erreur :", err)
	} else {
		fmt.Println("Définition :", result)
	}
}

func actionRemove(filePath string, reader *bufio.Reader) {
	fmt.Print("Entrez le mot à supprimer : ")
	key, _ := reader.ReadString('\n')
	key = strings.TrimSpace(key)

	dictionary.Remove(filePath, key)
	fmt.Println("Mot supprimé avec succès")
}

func actionList(filePath string) {
	result, err := dictionary.List(filePath)
	if err != nil {
		fmt.Println("Erreur :", err)
	} else {
		fmt.Println("Liste : \n", result)
	}
}
