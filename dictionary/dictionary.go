package dictionary

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Add(filePath string, key string, value string) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f.Close()

	_, err = f.WriteString(key + ":" + value + "\n")
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func Get(filePath string, key string) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, key+":") {
			fmt.Println(line)
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Key not found:", key)
	}
}

func List(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	empty := true
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		empty = false
	}

	if empty {
		fmt.Println("List Empty")
	}
}

func Remove(filePath string, key string) {
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	content := string(contentBytes)

	lines := strings.Split(content, "\n")

	var updatedLines []string
	for _, line := range lines {
		if !strings.HasPrefix(line, key+":") {
			updatedLines = append(updatedLines, line)
		}
	}

	updatedContent := strings.Join(updatedLines, "\n")

	err = os.WriteFile(filePath, []byte(updatedContent), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func Reset(filePath string) {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	fmt.Println("File content reset successfully.")
}
