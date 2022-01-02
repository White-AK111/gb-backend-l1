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
	file, _ := os.Open(uploadDir + "/testfile.txt")
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
		UploadDir: uploadDir,
		HostAddr:  ts.URL,
	}

	newName, ok := uploadHandler.GetUniqFilePath(uploadDir + "/" + "testfile.txt")
	if !ok {
		t.Errorf("error on get new file name")
	}

	uploadHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "File " + filepath.Base(newName) + " has been successfully uploaded"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	if _, err := os.Stat(newName); err == nil {
		return
	} else if os.IsNotExist(err) {
		t.Errorf("uploaded file " + filepath.Base(newName) + " not found in directory " + uploadDir)
	} else {
		t.Errorf("error on find  file " + filepath.Base(newName) + " in directory " + uploadDir)
	}
}

func TestListHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/list?ext=.jpg", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := &ListHandler{
		UploadDir: uploadDir,
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Name:jpgFile	Ext:.jpg	Size:0"

	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
