package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	MaxFileSize   = 10 << 20 // 10MB
	UploadDir     = "./uploads"
	AllowedTypes  = "image/jpeg,image/jpg,image/png,image/gif,image/webp"
)

type UploadResponse struct {
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
}

func init() {
	// アップロードディレクトリを作成
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}
}

func (h *Handler) UploadImage(c *gin.Context) {
	// ファイルサイズ制限を設定
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxFileSize)

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File too large or invalid"})
		}
		return
	}
	defer file.Close()

	// ファイルタイプをContent-Typeで検証
	if !isAllowedImageType(header.Header.Get("Content-Type")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed"})
		return
	}

	// マジックバイト検証でファイル内容を確認
	if !validateImageMagicBytes(file) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image file content. File content does not match the expected image format"})
		return
	}

	// ファイル拡張子を取得
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		// Content-Typeから拡張子を推定
		switch header.Header.Get("Content-Type") {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".jpg"
		}
	}

	// ユニークなファイル名を生成
	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), uuid.New().String()[:8], ext)
	filepath := filepath.Join(UploadDir, filename)

	// ファイルを保存
	dst, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		// ファイル削除
		os.Remove(filepath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// レスポンスを返す
	response := UploadResponse{
		Filename: filename,
		URL:      fmt.Sprintf("/uploads/%s", filename),
		Size:     size,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) ServeUpload(c *gin.Context) {
	filename := c.Param("filename")
	
	// Debug: Log the file request
	fmt.Printf("DEBUG: ServeUpload called for filename: %s\n", filename)
	
	// セキュリティ: パス traversal を防ぐ
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		fmt.Printf("DEBUG: Invalid filename rejected: %s\n", filename)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filename"})
		return
	}

	filepath := filepath.Join(UploadDir, filename)
	
	// Debug: Log the file path and check
	fmt.Printf("DEBUG: Looking for file at: %s\n", filepath)
	
	// ファイルが存在するかチェック
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		fmt.Printf("DEBUG: File not found: %s\n", filepath)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	fmt.Printf("DEBUG: Serving file: %s\n", filepath)
	c.File(filepath)
}

func (h *Handler) DeleteUpload(c *gin.Context) {
	filename := c.Param("filename")
	
	// セキュリティ: パス traversal を防ぐ
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filename"})
		return
	}

	filepath := filepath.Join(UploadDir, filename)
	
	// ファイルを削除
	if err := os.Remove(filepath); err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

func isAllowedImageType(contentType string) bool {
	allowedTypes := strings.Split(AllowedTypes, ",")
	for _, allowedType := range allowedTypes {
		if strings.TrimSpace(allowedType) == contentType {
			return true
		}
	}
	return false
}

// validateImageMagicBytes checks the file's magic bytes to ensure it's a valid image
func validateImageMagicBytes(file io.ReadSeeker) bool {
	// ファイルの最初の512バイトを読み取り（魔法の数値検出に十分）
	buffer := make([]byte, 512)
	bytesRead, err := file.Read(buffer)
	if err != nil {
		return false
	}

	// ファイルポインタを先頭に戻す
	_, err = file.Seek(0, 0)
	if err != nil {
		return false
	}

	// 実際に読み取ったバイト数でバッファをトリム
	if bytesRead < 512 {
		buffer = buffer[:bytesRead]
	}

	// マジックバイトで画像形式を検証
	return isValidImageMagicBytes(buffer)
}

// isValidImageMagicBytes checks if the buffer contains valid image magic bytes
func isValidImageMagicBytes(buffer []byte) bool {
	if len(buffer) < 4 {
		return false
	}

	// JPEG magic bytes: FF D8 FF
	if len(buffer) >= 3 && buffer[0] == 0xFF && buffer[1] == 0xD8 && buffer[2] == 0xFF {
		return true
	}

	// PNG magic bytes: 89 50 4E 47 0D 0A 1A 0A
	if len(buffer) >= 8 && 
		buffer[0] == 0x89 && buffer[1] == 0x50 && buffer[2] == 0x4E && buffer[3] == 0x47 &&
		buffer[4] == 0x0D && buffer[5] == 0x0A && buffer[6] == 0x1A && buffer[7] == 0x0A {
		return true
	}

	// GIF magic bytes: "GIF87a" or "GIF89a"
	if len(buffer) >= 6 && string(buffer[0:3]) == "GIF" && 
		(string(buffer[0:6]) == "GIF87a" || string(buffer[0:6]) == "GIF89a") {
		return true
	}

	// WebP magic bytes: "RIFF" followed by "WEBP" at offset 8
	if len(buffer) >= 12 && string(buffer[0:4]) == "RIFF" && string(buffer[8:12]) == "WEBP" {
		return true
	}

	return false
}