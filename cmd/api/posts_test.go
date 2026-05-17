package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/alanloffler/social/internal/store"
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

// mockSlowStore simula una DB que nunca responde
type mockSlowStore struct{}

func (m *mockSlowStore) Create(ctx context.Context, post *store.Post) error {
	return nil
}

func (m *mockSlowStore) GetByID(ctx context.Context, id int64) (*store.Post, error) {
	return &store.Post{ID: id, Title: "test", Content: "test", Version: 1}, nil
}

func (m *mockSlowStore) Delete(ctx context.Context, id int64) error {
	return nil
}

func (m *mockSlowStore) Update(ctx context.Context, post *store.Post) error {
	<-ctx.Done()
	return ctx.Err()
}

func TestUpdatePostContextTimeout(t *testing.T) {
	app := &application{
		store: store.Storage{Posts: &mockSlowStore{}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	ctx = context.WithValue(ctx, postCtx, &store.Post{ID: 1, Title: "test", Content: "test", Version: 1})

	req := httptest.NewRequest("PATCH", "/v1/posts/1", bytes.NewBufferString(`{"title":"new"}`))
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	app.updatePostHandler(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 Internal Server Error, got %d", rr.Code)
	}

	var resp map[string]string
	json.NewDecoder(rr.Body).Decode(&resp)
	t.Logf("response: %+v", resp)
	t.Logf("status: %d", rr.Code)
}
