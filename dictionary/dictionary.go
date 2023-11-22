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
}

var (
	addChan      = make(chan AddRequest)
	removeChan   = make(chan GetRemoveRequest)
	getChan      = make(chan GetRemoveRequest)
	listChan     = make(chan ResetListRequest)
	responseChan = make(chan Response)
)

type AddRequest struct {
	FilePath string
	Key      string
	Value    string
	Response chan Response
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

func Add(filePath string, key string, value string) (string, error) {
	responseChan := make(chan Response)
	addChan <- AddRequest{FilePath: filePath, Key: key, Value: value, Response: responseChan}
	response := <-responseChan
	return response.Result, response.Err
}

func Remove(filePath string, key string) (string, error) {
	responseChan := make(chan Response)
	removeChan <- GetRemoveRequest{FilePath: filePath, Key: key, Response: responseChan}
	response := <-responseChan
	return response.Result, response.Err
}

func Get(filePath string, key string) (string, error) {
	responseChan := make(chan Response)
	getChan <- GetRemoveRequest{FilePath: filePath, Key: key, Response: responseChan}
	response := <-responseChan
	return response.Result, response.Err
}

func List(filePath string) (string, error) {
	responseChan := make(chan Response)
	listChan <- ResetListRequest{FilePath: filePath, Response: responseChan}
	response := <-responseChan
	return response.Result, response.Err
}

func handleAddRequests() {
	for req := range addChan {
		exists, err := wordExists(req.FilePath, req.Key)
		if err != nil {
			req.Response <- Response{"", err}
			continue
		}
		if exists {
			req.Response <- Response{"", fmt.Errorf("Word '%s' already exists", req.Key)}
			continue
		}

		f, err := os.OpenFile(req.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			req.Response <- Response{"", err}
			continue
		}

		_, err = f.WriteString(req.Key + ":" + req.Value + "\n")
		closeErr := f.Close()
		if err != nil || closeErr != nil {
			req.Response <- Response{"", fmt.Errorf("%v %v", err, closeErr)}
			continue
		}

		req.Response <- Response{"Success", nil}
	}
}

func handleRemoveRequests() {
	for req := range removeChan {
		exists, err := wordExists(req.FilePath, req.Key)
		if err != nil {
			req.Response <- Response{"", err}
			continue
		}
		if !exists {
			req.Response <- Response{"", fmt.Errorf("Word '%s' does not exist", req.Key)}
			continue
		}

		contentBytes, err := os.ReadFile(req.FilePath)
		if err != nil {
			req.Response <- Response{"", err}
			continue
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
			req.Response <- Response{"", err}
			continue
		}

		req.Response <- Response{"Success", nil}
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
			req.Response <- Response{"", fmt.Errorf("%s not found", req.Key)}
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

func wordExists(filePath string, word string) (bool, error) {
	f, err := os.Open(filePath)
	if err != nil {
		// If the file does not exist, treat it as if the word does not exist.
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), word+":") {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}
