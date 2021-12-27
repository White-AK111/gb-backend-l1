package app

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

	err = ioutil.WriteFile(filePath, data, 0777)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "File %s has been successfully uploaded \n", header.Filename)
	fileLink := h.HostAddr + "/" + header.Filename
	fmt.Fprintln(w, fileLink)
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
