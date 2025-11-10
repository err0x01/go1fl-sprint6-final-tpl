package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func getBaseDir() (string, error) {
	paths := []string{
		".",
		"..",
		"../..",
		"../../..",
	}
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			continue
		}
		indexPath := filepath.Join(absPath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			return absPath, nil
		}
	}
	return "", fmt.Errorf("index.html not found in any parent directory")
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	baseDir, err := getBaseDir()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	indexPath := filepath.Join(baseDir, "index.html")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
	originalText := strings.TrimSpace(string(content))
	result, err := service.AutoDetectAndConvert(originalText)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting content: %v", err), http.StatusInternalServerError)
		return
	}
	baseDir, err := getBaseDir()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	response := fmt.Sprintf("Конвертация завершена успешно!\n\nИсходный текст: %s\n\nРезультат:\n%s\n\nФайл сохранен как: %s",
		originalText, result, filename)
	w.Write([]byte(response))
}

func generateFilename(originalFilename string) string {
	timestamp := time.Now().UTC().String()
	timestamp = strings.ReplaceAll(timestamp, ":", "-")
	timestamp = strings.ReplaceAll(timestamp, " ", "_")
	ext := filepath.Ext(originalFilename)
	if ext == "" {
		ext = ".txt"
	}
	return fmt.Sprintf("converted_%s%s", timestamp, ext)
}
