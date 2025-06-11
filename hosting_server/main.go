package main

import (
	"database/sql"
	"encoding/base64"
	"io/ioutil"
	"log"
	"time"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	swag "github.com/swaggo/swag/example/basic/docs"
)

// @title Music Server API
// @version 1.0
// @description API to serve audio files and metadata using Gin and SQLite3.
// @host localhost:8080
// @BasePath /

// Track represents audio file metadata.
type Track struct {
	ID       string `json:"id" example:"funky-lion-pencil-heart"`
	Title    string `json:"title"`
	ArtistID string `json:"artist_id"`
	AlbumID  string `json:"album_id"`
	FilePath string `json:"file_path"`
}

var db *sql.DB

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

// @Summary Get racks by Fuzzy Search
// @Produce json
// @Param query path string true "Search Query"
// @Success 200 {array} Track
// @Failure 404 {object} map[string]string
// @Router /search/{query} [get]
func getTracksByFuzzySearchHandler(c *gin.Context) {
	query := c.Param("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query cannot be empty"})
		return
	}

	rows, err := db.Query(`SELECT human_hash_id, title, artist_id, album_id, file_path FROM audio_files WHERE title LIKE ?`, "%"+query+"%")
	if err != nil {
		log.Printf("Query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	defer rows.Close()

	var tracks []Track
	for rows.Next() {
		var t Track
		if err := rows.Scan(&t.ID, &t.Title, &t.ArtistID, &t.AlbumID, &t.FilePath); err == nil {
			tracks = append(tracks, t)
		}
	}

	if len(tracks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no tracks found"})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

// @Summary Get all Albums using Fuzzy Search
// @Produce json
// @Param query path string true "Search Query"
// @Success 200 {array} Track
// @Failure 404 {object} map[string]string
// @Router /search/album/{query} [get]
func getAlbumsByFuzzySearchHandler(c *gin.Context) {
	query := c.Param("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query cannot be empty"})
		return
	}

	// Use the `albums` table to search for albums
	log.Printf("Searching for albums with query: %s", "'%"+query+"%'")
	rows, err := db.Query(`SELECT * FROM albums WHERE title LIKE ?`, "'%"+query+"%'")
	if err != nil {
		log.Printf("Query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	defer rows.Close()

	var tracks []Track
	for rows.Next() {
		var t Track
		if err := rows.Scan(&t.ID, &t.Title, &t.ArtistID, &t.AlbumID, &t.FilePath); err == nil {
			tracks = append(tracks, t)
		}
	}

	if len(tracks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no albums found"})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

// @Summary Get all tracks
// @Produce json
// @Success 200 {array} Track
// @Router /tracks/all [get]
func getAllTracksHandler(c *gin.Context) {
	rows, err := db.Query(`SELECT human_hash_id, title, artist_id, album_id, file_path FROM audio_files`)
	if err != nil {
		log.Printf("Query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	defer rows.Close()

	var tracks []Track
	for rows.Next() {
		var t Track
		if err := rows.Scan(&t.ID, &t.Title, &t.ArtistID, &t.AlbumID, &t.FilePath); err == nil {
			tracks = append(tracks, t)
		}
	}
	c.JSON(http.StatusOK, tracks)
}

// @Summary Get all tracks by artist ID
// @Produce json
// @Param artist_id path string true "Artist ID"
// @Success 200 {array} Track
// @Failure 404 {object} map[string]string
// @Router /artist/{artist_id} [get]
func getTracksByArtistHandler(c *gin.Context) {
	artistID := c.Param("artist_id")
	rows, err := db.Query(`SELECT human_hash_id, title, artist_id, album_id, file_path FROM audio_files WHERE artist_id = ?`, artistID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	defer rows.Close()

	var tracks []Track
	for rows.Next() {
		var t Track
		if err := rows.Scan(&t.ID, &t.Title, &t.ArtistID, &t.AlbumID, &t.FilePath); err == nil {
			tracks = append(tracks, t)
		}
	}

	if len(tracks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no tracks found"})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

// @Summary Get tracks of an album
// @Produce json
// @Param id path string true "Album ID"
// @Success 200 {array} Track
// @Failure 404 {object} map[string]string
// @Router /album/{id} [get]
func getTracksByAlbumHandler(c *gin.Context) {
	id := c.Param("id")
	rows, err := db.Query(`SELECT human_hash_id, title, artist_id, album_id, file_path FROM audio_files WHERE album_id = ?`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	defer rows.Close()

	var tracks []Track
	for rows.Next() {
		var t Track
		if err := rows.Scan(&t.ID, &t.Title, &t.ArtistID, &t.AlbumID, &t.FilePath); err == nil {
			tracks = append(tracks, t)
		}
	}

	if len(tracks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no tracks found for this album"})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

// @Summary Get album cover as base64
// @Produce json
// @Param id path string true "Track HumanHash ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /cover/{id} [get]
func getAlbumCoverHandler(c *gin.Context) {
	id := c.Param("id")
	var path string
	err := db.QueryRow("SELECT file_path FROM audio_files WHERE human_hash_id = ?", id).Scan(&path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
		return
	}
	dir := filepath.Dir(path)
	candidates := []string{"cover.jpg", "cover.png", "folder.jpg", "folder.png"}

	for _, name := range candidates {
		candidatePath := filepath.Join(dir, name)
		if data, err := ioutil.ReadFile(candidatePath); err == nil {
			encoded := base64.StdEncoding.EncodeToString(data)
			c.JSON(http.StatusOK, gin.H{"image_base64": "data:image/png;base64," + encoded})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "cover image not found"})
}

// @Summary Get track metadata
// @Produce json
// @Param id path string true "Track HumanHash ID"
// @Success 200 {object} Track
// @Failure 404 {object} map[string]string
// @Router /track/{id} [get]
func getTrackHandler(c *gin.Context) {
	id := c.Param("id")
	var t Track
	err := db.QueryRow(`SELECT human_hash_id, title, artist_id, album_id, file_path FROM audio_files WHERE human_hash_id = ?`, id).
		Scan(&t.ID, &t.Title, &t.ArtistID, &t.AlbumID, &t.FilePath)
	if err != nil {
		log.Printf("Error fetching track: %v", err);
		c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
		return
	}
	c.JSON(http.StatusOK, t)
}

// @Summary Stream audio file
// @Produce audio/flac
// @Param id path string true "Track HumanHash ID"
// @Success 200 {file} string
// @Failure 404 {object} map[string]string
// @Router /stream/{id} [get]
func streamTrackHandler(c *gin.Context) {
	id := c.Param("id")
	var path string
	err := db.QueryRow("SELECT file_path FROM audio_files WHERE human_hash_id = ?", id).Scan(&path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
		return
	}

	c.Header("Content-Type", "audio/flac") // Assuming FLAC; adapt if needed
	c.File(path)
}

func indexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Music Server API"})
}

func main() {
	dbPath := getEnv("MUSIC_DB_PATH", "music_library.sqlite")
	user := getEnv("MUSIC_USER", "admin")
	pass := getEnv("MUSIC_PASS", "admin123")
	port := getEnv("PORT", "8080")

	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"*"},
    AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
    AllowHeaders:     []string{"Origin"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    AllowOriginFunc: func(origin string) bool {
      return true
    },
    MaxAge: 12 * time.Hour,
  }));

	// Conditionally apply BasicAuth
	if strings.ToLower(user) != "0null" {
		auth := r.Group("/", gin.BasicAuth(gin.Accounts{user: pass}))
		auth.GET("", indexHandler)
		auth.GET("/track/:id", getTrackHandler)
		auth.GET("/stream/:id", streamTrackHandler)
		auth.GET("/tracks/all", getAllTracksHandler)
		auth.GET("/artist/:artist_id", getTracksByArtistHandler)
		auth.GET("/cover/:id", getAlbumCoverHandler)
		auth.GET("/album/:id", getTracksByAlbumHandler)
		auth.GET("/search/track/:query", getTracksByFuzzySearchHandler)
		auth.GET("/search/album/:query", getAlbumsByFuzzySearchHandler)
	} else {
		r.GET("", indexHandler)
		r.GET("/track/:id", getTrackHandler)
		r.GET("/stream/:id", streamTrackHandler)
		r.GET("/tracks/all", getAllTracksHandler)
		r.GET("/artist/:artist_id", getTracksByArtistHandler)
		r.GET("/cover/:id", getAlbumCoverHandler)
		r.GET("/album/:id", getTracksByAlbumHandler)
		r.GET("/search/track/:query", getTracksByFuzzySearchHandler)
		r.GET("/search/album/:query", getAlbumsByFuzzySearchHandler)
	}

	// Swagger docs
	swag.SwaggerInfo.Title = "Music Server API"
	swag.SwaggerInfo.Version = "1.0"
	swag.SwaggerInfo.BasePath = "/"
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + port)
}
