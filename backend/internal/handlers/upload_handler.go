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

	// ファイルタイプを検証
	if !isAllowedImageType(header.Header.Get("Content-Type")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed"})
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
	
	// セキュリティ: パス traversal を防ぐ
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filename"})
		return
	}

	filepath := filepath.Join(UploadDir, filename)
	
	// ファイルが存在するかチェック
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

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