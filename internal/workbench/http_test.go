package workbench

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer() *httptest.Server {
	store := NewFixtureStore()
	return httptest.NewServer(NewHandler(NewService(store)).Routes(nil))
}

func TestWorkbenchSectionSelection(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	response, err := http.Get(server.URL + "/api/articles?section=interview")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	var payload struct {
		Articles []Article `json:"articles"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusOK)
	}
	if len(payload.Articles) != 2 {
		t.Fatalf("articles = %d, want 2", len(payload.Articles))
	}
	for _, article := range payload.Articles {
		if article.Section != SectionInterview {
			t.Errorf("section = %q, want %q", article.Section, SectionInterview)
		}
	}
}

func TestEditorSavesArticle(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	update := ContentUpdate{
		Title:   "创客实验室周末开放安排",
		Summary: "新增预约时段已经确认。",
		Body:    "创客实验室将于周六上午和周日下午开放预约。",
	}
	body, err := json.Marshal(update)
	if err != nil {
		t.Fatal(err)
	}
	request, err := http.NewRequest(http.MethodPatch, server.URL+"/api/articles/news-lab", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	var article Article
	if err := json.NewDecoder(response.Body).Decode(&article); err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusOK)
	}
	if article.Title != update.Title || article.Summary != update.Summary || article.Body != update.Body {
		t.Fatalf("saved article = %#v, want %#v", article, update)
	}
}

func TestArchivedEditionRemainsFiled(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	transitionBody := bytes.NewBufferString(`{"status":"completed"}`)
	response, err := http.Post(server.URL+"/api/articles/archive-anniversary/transitions", "application/json", transitionBody)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusConflict {
		t.Errorf("status = %d, want %d", response.StatusCode, http.StatusConflict)
	}

	articleResponse, err := http.Get(server.URL + "/api/articles/archive-anniversary")
	if err != nil {
		t.Fatal(err)
	}
	defer articleResponse.Body.Close()
	var article Article
	if err := json.NewDecoder(articleResponse.Body).Decode(&article); err != nil {
		t.Fatal(err)
	}
	if article.Status != StatusArchived {
		t.Errorf("article status = %q, want %q", article.Status, StatusArchived)
	}

	queueResponse, err := http.Get(server.URL + "/api/queues/completed")
	if err != nil {
		t.Fatal(err)
	}
	defer queueResponse.Body.Close()
	var queue struct {
		Articles []Article `json:"articles"`
	}
	if err := json.NewDecoder(queueResponse.Body).Decode(&queue); err != nil {
		t.Fatal(err)
	}
	for _, completed := range queue.Articles {
		if completed.ID == "archive-anniversary" {
			t.Errorf("completed queue contains %q", completed.ID)
		}
	}
}
