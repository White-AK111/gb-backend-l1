package app

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/?name=Mikhail", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := &Handler{}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `Parsed query-param with key "name": Mikhail`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestUploadHandler(t *testing.T) {
	file, _ := os.Open("../../upload/testfile.txt")
	defer file.Close()

	var mime bytes.Buffer
	writer := multipart.NewWriter(&mime)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	_, err := io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/upload", &mime)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok!")
	}))
	defer ts.Close()

	uploadHandler := &UploadHandler{
		UploadDir: "../../upload",
		HostAddr:  ts.URL,
	}

	uploadHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `testfile.txt`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
