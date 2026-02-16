package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("hello world"))
	})

	addr := ":8081"
	fmt.Printf("API server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
