package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/AbdelrahmanKhaled21/musicauniversalis/config"
	"github.com/AbdelrahmanKhaled21/musicauniversalis/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type SongHandler struct {
	config *config.Config
	db     *gorm.DB
}

type UploadResponse struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Filename string `json:"filename"`
	FileSize int64  `json:"file_size"`
	Message  string `json:"message"`
}

func NewSongHandler(config *config.Config, db *gorm.DB) *SongHandler {
	return &SongHandler{
		config: config,
		db:     db,
	}
}

// UploadSong handles song file uploads
func (h *SongHandler) UploadSong(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(h.config.MaxFileSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	file, header, err := c.Request.FormFile("song")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No song file provided"})
		return
	}
	defer file.Close()

	// Validate file size
	if header.Size > h.config.MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := []string{".mp3", ".flac", ".wav", ".m4a", ".ogg"}
	allowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Allowed: mp3, flac, wav, m4a, ogg"})
		return
	}

	// Generate unique filename
	uniqueID := uuid.New().String()
	filename := uniqueID + ext
	filePath := fmt.Sprintf("users/%d/songs/%s", userID, filename)

	// Upload to MinIO
	ctx := context.Background()
	_, err = config.MinioClient.PutObject(ctx, h.config.MinIOBucket, filePath, file, header.Size, minio.PutObjectOptions{})
	if err != nil {
		log.Printf("Failed to upload to MinIO: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	// Extract metadata from form
	title := c.PostForm("title")
	if title == "" {
		title = strings.TrimSuffix(header.Filename, ext)
	}
	artist := c.PostForm("artist")
	album := c.PostForm("album")
	genre := c.PostForm("genre")
	yearStr := c.PostForm("year")
	trackNumberStr := c.PostForm("track_number")

	// Parse year and track number
	var year, trackNumber int
	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}
	if trackNumberStr != "" {
		if tn, err := strconv.Atoi(trackNumberStr); err == nil {
			trackNumber = tn
		}
	}

	// Create song record in database
	song := models.Song{
		UserID:      userID.(uint),
		Title:       title,
		Artist:      artist,
		Album:       album,
		Genre:       genre,
		Year:        year,
		TrackNumber: trackNumber,
		Filename:    header.Filename,
		FilePath:    filePath,
		FileSize:    header.Size,
		MimeType:    header.Header.Get("Content-Type"),
	}

	if err := h.db.Create(&song).Error; err != nil {
		log.Printf("Failed to create song record: %v", err)
		// Try to clean up MinIO file
		config.MinioClient.RemoveObject(ctx, h.config.MinIOBucket, filePath, minio.RemoveObjectOptions{})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save song metadata"})
		return
	}

	// Update user storage usage
	var user models.User
	if err := h.db.First(&user, userID).Error; err == nil {
		h.db.Model(&user).Update("storage_used", user.StorageUsed+header.Size)
	}

	response := UploadResponse{
		ID:       song.ID,
		Title:    song.Title,
		Artist:   song.Artist,
		Filename: song.Filename,
		FileSize: song.FileSize,
		Message:  "Song uploaded successfully",
	}

	c.JSON(http.StatusCreated, response)
}

// GetUserSongs returns all songs for the authenticated user
func (h *SongHandler) GetUserSongs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var songs []models.Song
	if err := h.db.Where("user_id = ? AND is_deleted = ?", userID, false).Find(&songs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch songs"})
		return
	}

	c.JSON(http.StatusOK, songs)
}

// GetSong returns a specific song by ID
func (h *SongHandler) GetSong(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	songID := c.Param("id")
	var song models.Song

	if err := h.db.Where("id = ? AND user_id = ? AND is_deleted = ?", songID, userID, false).First(&song).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch song"})
		return
	}

	c.JSON(http.StatusOK, song)
}

// StreamSong streams a song file from MinIO
func (h *SongHandler) StreamSong(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	songID := c.Param("id")
	var song models.Song

	if err := h.db.Where("id = ? AND user_id = ? AND is_deleted = ?", songID, userID, false).First(&song).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch song"})
		return
	}

	// Get object from MinIO
	ctx := context.Background()
	obj, err := config.MinioClient.GetObject(ctx, h.config.MinIOBucket, song.FilePath, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Failed to get object from MinIO: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream song"})
		return
	}
	defer obj.Close()

	// Set response headers
	c.Header("Content-Type", song.MimeType)
	c.Header("Content-Length", strconv.FormatInt(song.FileSize, 10))
	c.Header("Accept-Ranges", "bytes")

	// Stream the file
	io.Copy(c.Writer, obj)
}

// DeleteSong soft deletes a song
func (h *SongHandler) DeleteSong(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	songID := c.Param("id")
	var song models.Song

	if err := h.db.Where("id = ? AND user_id = ? AND is_deleted = ?", songID, userID, false).First(&song).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch song"})
		return
	}

	// Soft delete
	now := time.Now()
	song.IsDeleted = true
	song.DeletedAt = &now

	if err := h.db.Save(&song).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}

	// Update user storage usage
	var user models.User
	if err := h.db.First(&user, userID).Error; err == nil {
		h.db.Model(&user).Update("storage_used", user.StorageUsed-song.FileSize)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}
