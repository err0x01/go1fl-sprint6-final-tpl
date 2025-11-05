package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

var baseDir string

func Init() error {
	dir, err := filepath.Abs("..")
	if err != nil {
		dir, err = filepath.Abs(".")
		if err != nil {
			return err
		}
	}
	baseDir = dir
	return nil
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	indexPath := filepath.Join(baseDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		http.Error(w, fmt.Sprintf("index.html not found at %s", indexPath), http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, indexPath)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusInternalServerError)
		return
	}
	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
		return
	}
	result, err := service.AutoDetectAndConvert(string(content))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting content: %v", err), http.StatusInternalServerError)
		return
	}
	filename := generateFilename(header.Filename)
	filePath := filepath.Join(baseDir, filename)
	if err := os.WriteFile(filePath, []byte(result), 0644); err != nil {
		http.Error(w, fmt.Sprintf("Error saving result file: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf("Конвертация завершена успешно!\n\nРезультат:\n%s\n\nФайл сохранен как: %s", result, filename)
	w.Write([]byte(response))
}

func generateFilename(originalFilename string) string {
	timestamp := time.Now().UTC().Format("2006-01-02_15-04-05")
	ext := filepath.Ext(originalFilename)
	if ext == "" {
		ext = ".txt"
	}
	return fmt.Sprintf("converted_%s%s", timestamp, ext)
}
