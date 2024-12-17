// messy code sorry ;)
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// Config holds service configuration
type Config struct {
	UploadDir    string
	DownloadDir  string
	MaxFileSize  int64
	AllowedTypes []string
}

// Service handles PDF compression operations
type Service struct {
	config Config
}

// NewService creates a new PDF compression service
func NewService() *Service {
	uploadDir := "./uploads"
	downloadDir := "./downloads"

	for _, dir := range []string{uploadDir, downloadDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	return &Service{
		config: Config{
			UploadDir:    uploadDir,
			DownloadDir:  downloadDir,
			MaxFileSize:  10 << 20, // 10MB
			AllowedTypes: []string{"application/pdf"},
		},
	}
}

// UploadHandler handles PDF file uploads
func (s *Service) UploadHandler(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Validate file size
	if header.Size > s.config.MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	// Validate file type
	if !contains(s.config.AllowedTypes, header.Header.Get("Content-Type")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	// Generate unique filename
	filename := uuid.New().String() + ".pdf"
	uploadPath := filepath.Join(s.config.UploadDir, filename)
	downloadPath := filepath.Join(s.config.DownloadDir, filename)

	// Save the uploaded file
	out, err := os.Create(uploadPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Compress the PDF
	err = s.compressPDF(uploadPath, downloadPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to compress PDF: %v", err)})
		return
	}

	// Remove the original file after compression
	os.Remove(uploadPath)

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded and compressed successfully",
		"fileId":  filename,
	})
}

// DownloadHandler handles compressed PDF downloads
func (s *Service) DownloadHandler(c *gin.Context) {
	fileID := c.Param("fileId")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file ID provided"})
		return
	}

	filepath := filepath.Join(s.config.DownloadDir, fileID)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.FileAttachment(filepath, "compressed.pdf")
}

// Compresses a PDF file using pdfcpu
func (s *Service) compressPDF(inputPath, outputPath string) error {
	conf := model.NewDefaultConfiguration()
	return api.OptimizeFile(inputPath, outputPath, conf)
}

// Removes files older than 24 hours from both directories
func (s *Service) cleanupOldFiles() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			threshold := time.Now().Add(-24 * time.Hour)

			// Clean both directories
			for _, dir := range []string{s.config.UploadDir, s.config.DownloadDir} {
				err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() && info.ModTime().Before(threshold) {
						os.Remove(path)
					}
					return nil
				})
				if err != nil {
					log.Printf("Error cleaning up old files in %s: %v", dir, err)
				}
			}
		}
	}()
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func main() {
	service := NewService()

	// Start cleanup routine
	service.cleanupOldFiles()

	router := gin.Default()

	// Configure maximum multipart form size
	router.MaxMultipartMemory = service.config.MaxFileSize

	// Setup routes
	router.POST("/upload", service.UploadHandler)
	router.GET("/download/:fileId", service.DownloadHandler)

	// Start server
	if err := router.Run(":1337"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
