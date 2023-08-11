package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s file.md\n", os.Args[0])
		return
	}

	path := os.Args[1]

	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	payload := map[string]string{
		"text": string(content),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	url := "https://api.github.com/markdown"
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	req.Header.Set("accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(respBody)
	})
	http.ListenAndServe(":8080", nil)
}
