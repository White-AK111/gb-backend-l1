package app

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/White-AK111/gb-backend-l1/homework-4/internal/pkg/models"
)

const apiPort = 8000
const filePort = 8080

type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		name := r.FormValue("name")
		fmt.Fprintf(w, "Parsed query-param with key \"name\": %s", name)
	case http.MethodPost:
		defer r.Body.Close()

		var employee models.Employee
		contentType := r.Header.Get("Content-Type")

		switch contentType {
		case "application/json":
			err := json.NewDecoder(r.Body).Decode(&employee)
			if err != nil {
				http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
				return
			}
		case "application/xml":
			err := xml.NewDecoder(r.Body).Decode(&employee)
			if err != nil {
				http.Error(w, "Unable to unmarshal XML", http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "Unknown content type", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Got a new employee!\nName: %s\nAge: %dy.o.\nSalary %0.2f\n",
			employee.Name,
			employee.Age,
			employee.Salary,
		)
	}
}

type UploadHandler struct {
	HostAddr  string
	UploadDir string
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}

	filePath := h.UploadDir + "/" + header.Filename

	uniqFilePath, ok := h.GetUniqFilePath(filePath)
	if !ok {
		http.Error(w, "Unable to fet unique filename", http.StatusInternalServerError)
		return
	}

	err = ioutil.WriteFile(uniqFilePath, data, 0777)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "File %s has been successfully uploaded \n", filepath.Base(uniqFilePath))
	fileLink := h.HostAddr + "/" + filepath.Base(uniqFilePath)
	fmt.Fprintln(w, fileLink)
}

func (h *UploadHandler) GetUniqFilePath(path string) (uniqFilePath string, ok bool) {
	if _, err := os.Stat(path); err == nil {
		ext := filepath.Ext(path)
		return h.GetUniqFilePath(filepath.Dir(path) + "/" + strings.TrimSuffix(filepath.Base(path), ext) + "_copy" + ext)
	} else if os.IsNotExist(err) {
		return path, true
	} else {
		return "", false
	}
}

type ListHandler struct {
	UploadDir string
}

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		paramExt := r.FormValue("ext")
		fmt.Fprintf(w, "Parsed query-param with key \"ext\": %s\n", paramExt)

		files, err := ioutil.ReadDir(h.UploadDir)
		if err != nil {
			http.Error(w, "Unable to read upload directory", http.StatusBadRequest)
			return
		}

		for _, file := range files {
			ext := filepath.Ext(file.Name())
			if strings.EqualFold(ext, paramExt) || len(paramExt) == 0 {
				fmt.Fprintf(w, "Name:%s\tExt:%s\tSize:%d\n", strings.TrimSuffix(file.Name(), ext), ext, file.Size())
			}
		}
	default:
		http.Error(w, "Unknown content type", http.StatusMethodNotAllowed)
		return
	}
}

func App() {
	go startFileServer()
	startApiServer()
}

func startApiServer() {
	handler := &Handler{}
	http.Handle("/", handler)

	uploadHandler := &UploadHandler{
		HostAddr:  "http://localhost:" + strconv.Itoa(filePort) + "/",
		UploadDir: "../../upload",
	}
	http.Handle("/upload", uploadHandler)

	listHandler := &ListHandler{
		UploadDir: "../../upload",
	}
	http.Handle("/list", listHandler)

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", apiPort),
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("api start listening on %d port\n", apiPort)
	log.Fatal(server.ListenAndServe())
}

func startFileServer() {
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", filePort),
		Handler:      http.FileServer(http.Dir("/home/white/GolandProjects/gb-backend-l1/homework-4/upload/")),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("fs starts listening on %d port\n", filePort)
	log.Fatal(server.ListenAndServe())
}
