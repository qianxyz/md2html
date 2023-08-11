package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Convert the markdown string to a html string.
// Both strings are represented as byte slices.
func convert(md []byte, client *http.Client) ([]byte, error) {
	payload := map[string]string{
		"text": string(md),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := "https://api.github.com/markdown"
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

// Serve the html string.
func serve(html []byte, port int) {
	fmt.Printf("Serving at http://0.0.0.0:%d ...\n", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(html)
	})

	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func main() {
	pFlag := flag.Int("p", 8080, "port (default: 8080)")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: %s [-p port] file.md\n", os.Args[0])
		return
	}

	path := flag.Arg(0)
	port := *pFlag

	md, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	html, err := convert(md, client)
	if err != nil {
		log.Fatal(err)
	}

	serve(html, port)
}
