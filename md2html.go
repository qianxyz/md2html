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
	"sync"

	"github.com/fsnotify/fsnotify"
)

var (
	// The HTML string to be rendered.
	rendered []byte
	mu       sync.RWMutex

	// The path of the Markdown file.
	mdPath string
	// The port to serve on.
	port int

	client = &http.Client{}
)

func init() {
	pFlag := flag.Int("p", 8080, "port (default: 8080)")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: %s [-p port] file.md\n", os.Args[0])
		os.Exit(0)
	}

	mdPath = flag.Arg(0)
	port = *pFlag
}

// Read the file and convert its content into a HTML string,
// which is stored in the global variable.
func update() {
	md, err := os.ReadFile(mdPath)
	if err != nil {
		log.Fatal(err)
	}

	payload := map[string]string{
		"text": string(md),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	url := "https://api.github.com/markdown"
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/vnd.github+json")

	log.Println("Sending API request ...")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	mu.Lock()
	rendered = respBody
	mu.Unlock()

	log.Printf("HTML updated at http://0.0.0.0:%d\n", port)
}

// Serve the html string.
func serve() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		mu.RLock()
		defer mu.RUnlock()
		w.Write(rendered)
	})

	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// Watch a file, rerender HTML when it has been modified.
func watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	addPath := func() {
		err = watcher.Add(mdPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	addPath()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) {
				log.Println(mdPath, "modified")
				update()
			}
			// Some editor does not write to the file, but creates
			// a new file and renames it (thus removing the
			// original). In this case the watcher loses track of
			// the file, so the path has to be added again.
			if event.Has(fsnotify.Remove) {
				log.Println(mdPath, "modified")
				addPath()
				update()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func main() {
	update()

	go serve()

	watch()
}
