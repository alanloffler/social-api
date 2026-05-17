package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
)

func updatePost(t *testing.T, postID int, p UpdatePostPayload, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("http://localhost:8080/v1/posts/%d", postID)
	b, _ := json.Marshal(p)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(b))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	t.Logf("Update response status: %s", resp.Status)
}

// curl -s -o /dev/null -w "%{http_code}" -X PATCH http://localhost:8080/v1/posts/6 \
//   -H "Content-Type: application/json" \
//   -d '{"title":"test"}'

func TestUpdatePost(t *testing.T) {
	var wg sync.WaitGroup

	postID := 6

	wg.Add(2)
	content := "New content from user B"
	title := "New title from user A"

	go updatePost(t, postID, UpdatePostPayload{Title: &title}, &wg)
	go updatePost(t, postID, UpdatePostPayload{Content: &content}, &wg)

	wg.Wait()
}
