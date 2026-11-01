package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestManifestETag(t *testing.T) {
	h := routes("hyena-ufa-plugin")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/ui-manifest.json", nil))
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), `"app_id":"hyena-ufa-plugin"`) {
		t.Fatalf("manifest: %d %s", rec.Code, rec.Body)
	}
	etag := rec.Header().Get("ETag")
	if etag == "" {
		t.Fatal("no ETag")
	}
	req := httptest.NewRequest("GET", "/ui-manifest.json", nil)
	req.Header.Set("If-None-Match", etag)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotModified {
		t.Fatalf("want 304, got %d", rec.Code)
	}
}

func TestPageAndIcon(t *testing.T) {
	h := routes("hyena-ufa-plugin")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	body := rec.Body.String()
	if !strings.Contains(body, "Hello from the Hyena") || strings.Contains(body, helloPlaceholder) {
		t.Fatalf("page not filled in: %s", body)
	}
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/plugin-icon.svg", nil))
	if rec.Code != 200 {
		t.Fatalf("icon: %d", rec.Code)
	}
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/nope", nil))
	if rec.Code != 404 {
		t.Fatalf("want 404, got %d", rec.Code)
	}
}
