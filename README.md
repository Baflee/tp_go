# Dictionary Service

![6ba8132e-ff5d-4417-978f-e7dde9b7b724](https://github.com/Baflee/tp_go/assets/63551775/8b4f4097-0f40-4f31-bfda-648ffb3b3153)

## Introduction
This project is a dictionary service implemented in Go, using Gorilla Mux for handling HTTP requests. It allows users to add, remove, retrieve, and list words along with their definitions in a simple dictionary application. The dictionary data is managed in a text file, and the service is designed to handle concurrent requests safely.

## Features
- Add new word definitions to the dictionary.
- Retrieve definitions of words.
- Remove words from the dictionary.
- List all words in the dictionary.

## Setup

### Requirements
- Go (version 1.15 or higher)
- Access to terminal or command line

### Installation
1. Clone the repository:
   ```bash
   git clone [URL to your repository]
   ```
2. Navigate to the project directory:
   ```bash
   cd [project-directory]
   ```
3. Build the project (optional):
   ```bash
   go build
   ```

### Running the Service
Run the service using the Go command:
```bash
go run main.go
```

## Usage

### Adding a Word
Send a `POST` request to `/dictionary/add` with the word and its definition.
```bash
curl -X POST localhost:8080/dictionary/add -d "word=example&definition=This is an example."
```

### Retrieving a Word's Definition
Send a `GET` request to `/dictionary/{word}`.
```bash
curl localhost:8080/dictionary/example
```

### Removing a Word
Send a `DELETE` request to `/dictionary/delete/{word}`.
```bash
curl -X DELETE localhost:8080/dictionary/delete/example
```

### Listing All Words
Send a `GET` request to `/dictionary`.
```bash
curl localhost:8080/dictionary
```

## Architecture

- **Gorilla Mux Router:** Manages HTTP requests and routes them to appropriate handlers.
- **Concurrent Request Handling:** Utilizes Go channels and goroutines for concurrent processing of dictionary operations.
- **DB-Based Storage:** Dictionary data is stored and managed in a database NoSQL called MongoDB.

## Notes

- The service is designed for educational and demonstration purposes and is not optimized for production use.
- Concurrency control is implemented to handle simultaneous read/write operations to the dictionary file.
