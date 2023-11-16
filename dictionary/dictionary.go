package dictionary

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func init() {
	go handleAddRequests()
	go handleRemoveRequests()
	go handleGetRequests()
	go handleListRequests()
	go handleResetRequests()
}

var (
	addChan      = make(chan AddRequest)
	removeChan   = make(chan GetRemoveRequest)
	getChan      = make(chan GetRemoveRequest)
	listChan     = make(chan ResetListRequest)
	resetChan    = make(chan ResetListRequest)
	responseChan = make(chan Response)
)

type AddRequest struct {
	FilePath string
	Key      string
	Value    string
}

type GetRemoveRequest struct {
	FilePath string
	Key      string
	Response chan Response
}

type ResetListRequest struct {
	FilePath string
	Response chan Response
}

type Response struct {
	Result string
	Err    error
}

func Add(filePath string, key string, value string) {
	addChan <- AddRequest{FilePath: filePath, Key: key, Value: value}
}

func Remove(filePath string, key string) {
	removeChan <- GetRemoveRequest{FilePath: filePath, Key: key}
}

func Get(filePath string, key string) (string, error) {
	getChan <- GetRemoveRequest{FilePath: filePath, Key: key, Response: responseChan}
	response := <-responseChan
	return response.Result, response.Err
}

func List(filePath string) (string, error) {
	listChan <- ResetListRequest{FilePath: filePath, Response: responseChan}
	response := <-responseChan
	return response.Result, response.Err
}

func Reset(filePath string) {
	resetChan <- ResetListRequest{FilePath: filePath}
}

func handleAddRequests() {
	for req := range addChan {
		f, err := os.OpenFile(req.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		_, err = f.WriteString(req.Key + ":" + req.Value + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}
		f.Close()
	}
}

func handleRemoveRequests() {
	for req := range removeChan {
		contentBytes, err := os.ReadFile(req.FilePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		content := string(contentBytes)

		lines := strings.Split(content, "\n")

		var updatedLines []string
		for _, line := range lines {
			if !strings.HasPrefix(line, req.Key+":") {
				updatedLines = append(updatedLines, line)
			}
		}

		updatedContent := strings.Join(updatedLines, "\n")

		err = os.WriteFile(req.FilePath, []byte(updatedContent), 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
}

func handleGetRequests() {
	for req := range getChan {
		f, err := os.Open(req.FilePath)
		if err != nil {
			req.Response <- Response{"", err}
			continue
		}

		scanner := bufio.NewScanner(f)
		found := false
		line := ""
		for scanner.Scan() {
			line = scanner.Text()
			if strings.HasPrefix(line, req.Key+":") {
				found = true
				line = strings.TrimPrefix(line, req.Key+":")
				break
			}
		}
		f.Close()

		if found {
			req.Response <- Response{line, nil}
		} else {
			req.Response <- Response{"", fmt.Errorf("Key not found: %s", req.Key)}
		}
	}
}

func handleListRequests() {
	for req := range listChan {
		f, err := os.Open(req.FilePath)
		if err != nil {
			req.Response <- Response{"", err}
			continue
		}

		scanner := bufio.NewScanner(f)
		var lines []string

		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			req.Response <- Response{"", err}
		} else if len(lines) == 0 {
			req.Response <- Response{"Empty", nil}
		} else {
			combinedLines := strings.Join(lines, "\n")
			req.Response <- Response{combinedLines, nil}
		}

		f.Close()
	}
}

func handleResetRequests() {
	for req := range resetChan {
		f, err := os.OpenFile(req.FilePath, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer f.Close()

		fmt.Println("File content reset successfully.")
	}
}
