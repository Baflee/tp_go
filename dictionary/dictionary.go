package dictionary

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
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
	Http   int
}

func Add(filePath string, key string, value string) (string, error, int) {
	responseChan := make(chan Response)
	addChan <- AddRequest{FilePath: filePath, Key: key, Value: value, Response: responseChan}
	response := <-responseChan
	return response.Result, response.Err, response.Http
}

func Remove(filePath string, key string) (string, error, int) {
	responseChan := make(chan Response)
	removeChan <- GetRemoveRequest{FilePath: filePath, Key: key, Response: responseChan}
	response := <-responseChan
	return response.Result, response.Err, response.Http
}

func Get(filePath string, key string) (string, error, int) {
	responseChan := make(chan Response)
	getChan <- GetRemoveRequest{FilePath: filePath, Key: key, Response: responseChan}
	response := <-responseChan
	return response.Result, response.Err, response.Http
}

func List(filePath string) (string, error, int) {
	responseChan := make(chan Response)
	listChan <- ResetListRequest{FilePath: filePath, Response: responseChan}
	response := <-responseChan
	return response.Result, response.Err, response.Http
}

func handleAddRequests() {
	for req := range addChan {
		contentFile, err := checkFile(req.FilePath)

		if err != nil {
			req.Response <- Response{"", fmt.Errorf("File '%s' not found", req.FilePath), http.StatusNotFound}
			continue
		}

		exists, err := wordExists(contentFile, err, req.Key)

		if len(req.Key) < 3 || len(req.Key) > 20 || len(req.Value) < 5 {
			req.Response <- Response{"", fmt.Errorf("Invalid input data"), http.StatusBadRequest}
			continue
		}

		if err != nil {
			req.Response <- Response{"", err, http.StatusInternalServerError}
			continue
		}

		if exists {
			req.Response <- Response{"", fmt.Errorf("Word '%s' already exists", req.Key), http.StatusConflict}
			continue
		}

		f, err := os.OpenFile(req.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			req.Response <- Response{"", err, http.StatusInternalServerError}
			continue
		}

		_, err = f.WriteString(req.Key + ":" + req.Value + "\n")
		closeErr := f.Close()
		if err != nil || closeErr != nil {
			req.Response <- Response{"", fmt.Errorf("%v %v", err, closeErr), http.StatusInternalServerError}
			continue
		}

		req.Response <- Response{"Success", nil, http.StatusOK}
	}
}

func handleRemoveRequests() {
	for req := range removeChan {
		contentFile, err := checkFile(req.FilePath)

		if err != nil {
			req.Response <- Response{"", fmt.Errorf("File '%s' not found", req.FilePath), http.StatusNotFound}
			continue
		}

		exists, err := wordExists(contentFile, err, req.Key)

		if err != nil {
			req.Response <- Response{"", err, http.StatusInternalServerError}
			continue
		}

		if !exists {
			req.Response <- Response{"", fmt.Errorf("Word '%s' does not exist", req.Key), http.StatusNotFound}
			continue
		}

		contentBytes, err := os.ReadFile(req.FilePath)
		if err != nil {
			req.Response <- Response{"", err, http.StatusInternalServerError}
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
			req.Response <- Response{"", err, http.StatusInternalServerError}
			continue
		}

		req.Response <- Response{"Success", nil, http.StatusOK}
	}
}

func handleGetRequests() {
	for req := range getChan {
		contentFile, err := checkFile(req.FilePath)
		if err != nil {
			req.Response <- Response{"", fmt.Errorf("File '%s' not found", req.FilePath), http.StatusNotFound}
			continue
		}

		reader := bufio.NewReader(bytes.NewReader(contentFile))
		scanner := bufio.NewScanner(reader)
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

		if found {
			req.Response <- Response{line, nil, http.StatusOK}
		} else {
			req.Response <- Response{"", fmt.Errorf("%s not found", req.Key), http.StatusNotFound}
		}
	}
}

func handleListRequests() {
	for req := range listChan {
		contentFile, err := checkFile(req.FilePath)

		if err != nil {
			req.Response <- Response{"", fmt.Errorf("File '%s' not found", req.FilePath), http.StatusNotFound}
			continue
		}

		reader := bufio.NewReader(bytes.NewReader(contentFile))
		scanner := bufio.NewScanner(reader)
		var lines []string

		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			req.Response <- Response{"", err, http.StatusInternalServerError}
		} else if len(lines) == 0 {
			req.Response <- Response{"Empty", nil, http.StatusOK}
		} else {
			combinedLines := strings.Join(lines, "\n")
			req.Response <- Response{combinedLines, nil, http.StatusOK}
		}
	}
}

func checkFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(f)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	return content, nil
}

func wordExists(contentFile []byte, err error, word string) (bool, error) {
	reader := bufio.NewReader(bytes.NewReader(contentFile))
	scanner := bufio.NewScanner(reader)
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
