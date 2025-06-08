package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync" // Import sync package for WaitGroup
	"time"

	"github.com/dhowden/tag"
	"github.com/mattn/go-sqlite3"
	"github.com/wolfeidau/humanhash"
)

// AudioFile represents the metadata for an audio track.
type AudioFile struct {
	HumanHashID     string
	FilePath        string
	Title           string
	DurationSeconds int // Note: dhowden/tag does not directly provide duration. This will be 0 unless another library is integrated.
	Lossless        bool
	TrackNumber     int
	DiscNumber      int
	Year            int
	ArtistName      string
	AlbumTitle      string
	GenreName       string
	ArtistID        int // Foreign key after insertion
	AlbumID         int // Foreign key after insertion
	GenreID         int // Foreign key after insertion
}

// DBManager handles all database operations.
type DBManager struct {
	db *sql.DB
	mu sync.Mutex // Mutex to protect database writes from concurrent access
}

// NewDBManager creates a new DBManager and initializes the database connection.
func NewDBManager(dbPath string) (*DBManager, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	// Set a timeout for database operations (e.g., 5 seconds)
	// For SQLite, it's generally best to keep MaxOpenConns low for writes.
	// We'll manage concurrent writes using a mutex.
	db.SetMaxOpenConns(1) // Only one active connection to prevent SQLite locking issues with concurrent writes
	db.SetConnMaxLifetime(5 * time.Minute)

	mgr := &DBManager{db: db}
	if err := mgr.InitDB(); err != nil {
		db.Close() // Close on initialization failure
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}
	return mgr, nil
}

// Close closes the database connection.
func (m *DBManager) Close() error {
	return m.db.Close()
}

// InitDB creates the necessary tables if they don't exist.
func (m *DBManager) InitDB() error {
	schema := `
	PRAGMA foreign_keys = ON;

	CREATE TABLE IF NOT EXISTS artists (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL COLLATE NOCASE
	);

	CREATE TABLE IF NOT EXISTS albums (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL COLLATE NOCASE,
		artist_id INTEGER NOT NULL,
		release_year INTEGER,
		FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE,
		UNIQUE(title, artist_id)
	);

	CREATE TABLE IF NOT EXISTS genres (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL COLLATE NOCASE
	);

	CREATE TABLE IF NOT EXISTS audio_files (
		human_hash_id TEXT PRIMARY KEY NOT NULL,
		file_path TEXT UNIQUE NOT NULL,
		title TEXT NOT NULL,
		duration_seconds INTEGER,
		lossless BOOLEAN NOT NULL,
		track_number INTEGER,
		disc_number INTEGER,
		year INTEGER,
		artist_id INTEGER NOT NULL,
		album_id INTEGER NOT NULL,
		genre_id INTEGER,
		FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE,
		FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE,
		FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE SET NULL
	);
	`
	// Acquire mutex for schema creation to ensure it's not run concurrently
	m.mu.Lock()
	defer m.mu.Unlock()
	_, err := m.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("error creating database schema: %w", err)
	}
	return nil
}

// GetOrInsertArtist retrieves an artist's ID or inserts a new artist if not found.
func (m *DBManager) GetOrInsertArtist(name string) (int, error) {
	// Acquire mutex for database write operations
	m.mu.Lock()
	defer m.mu.Unlock()

	var id int
	// Try to get existing artist
	err := m.db.QueryRow("SELECT id FROM artists WHERE name = ? COLLATE NOCASE", name).Scan(&id)
	if err == nil {
		return id, nil // Artist found
	}
	if err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to query artist: %w", err)
	}

	// Artist not found, insert new one
	res, err := m.db.Exec("INSERT INTO artists (name) VALUES (?)", name)
	if err != nil {
		// Handle potential race condition if another goroutine inserted it between check and insert
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			// Try to get again if it was a unique constraint error
			err = m.db.QueryRow("SELECT id FROM artists WHERE name = ? COLLATE NOCASE", name).Scan(&id)
			if err == nil {
				return id, nil
			}
		}
		return 0, fmt.Errorf("failed to insert artist: %w", err)
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last inserted artist ID: %w", err)
	}
	return int(lastID), nil
}

// GetOrInsertAlbum retrieves an album's ID or inserts a new album if not found.
func (m *DBManager) GetOrInsertAlbum(title string, artistID int, releaseYear int) (int, error) {
	// Acquire mutex for database write operations
	m.mu.Lock()
	defer m.mu.Unlock()

	var id int
	err := m.db.QueryRow("SELECT id FROM albums WHERE title = ? COLLATE NOCASE AND artist_id = ?", title, artistID).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to query album: %w", err)
	}

	// Album not found, insert new one
	res, err := m.db.Exec("INSERT INTO albums (title, artist_id, release_year) VALUES (?, ?, ?)", title, artistID, releaseYear)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			err = m.db.QueryRow("SELECT id FROM albums WHERE title = ? COLLATE NOCASE AND artist_id = ?", title, artistID).Scan(&id)
			if err == nil {
				return id, nil
			}
		}
		return 0, fmt.Errorf("failed to insert album: %w", err)
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last inserted album ID: %w", err)
	}
	return int(lastID), nil
}

// GetOrInsertGenre retrieves a genre's ID or inserts a new genre if not found.
func (m *DBManager) GetOrInsertGenre(name string) (int, error) {
	if name == "" { // Handle empty genre gracefully
		return 0, nil
	}
	// Acquire mutex for database write operations
	m.mu.Lock()
	defer m.mu.Unlock()

	var id int
	err := m.db.QueryRow("SELECT id FROM genres WHERE name = ? COLLATE NOCASE", name).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to query genre: %w", err)
	}

	// Genre not found, insert new one
	res, err := m.db.Exec("INSERT INTO genres (name) VALUES (?)", name)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			err = m.db.QueryRow("SELECT id FROM genres WHERE name = ? COLLATE NOCASE", name).Scan(&id)
			if err == nil {
				return id, nil
			}
		}
		return 0, fmt.Errorf("failed to insert genre: %w", err)
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last inserted genre ID: %w", err)
	}
	return int(lastID), nil
}

// InsertAudioFile inserts an audio file record into the database.
func (m *DBManager) InsertAudioFile(af *AudioFile) error {
	// Acquire mutex for database write operations
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the audio file already exists by its human hash ID or file path
	var existingHumanHashID string
	err := m.db.QueryRow("SELECT human_hash_id FROM audio_files WHERE human_hash_id = ? OR file_path = ?", af.HumanHashID, af.FilePath).Scan(&existingHumanHashID)
	if err == nil {
		log.Printf("Skipping existing audio file: %s (Human Hash: %s)", af.FilePath, existingHumanHashID)
		return nil // File already exists, skip
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check for existing audio file: %w", err)
	}

	_, err = m.db.Exec(`
		INSERT INTO audio_files (human_hash_id, file_path, title, duration_seconds, lossless, track_number, disc_number, year, artist_id, album_id, genre_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, af.HumanHashID, af.FilePath, af.Title, af.DurationSeconds, af.Lossless, af.TrackNumber, af.DiscNumber, af.Year, af.ArtistID, af.AlbumID, sql.NullInt64{Int64: int64(af.GenreID), Valid: af.GenreID != 0})
	if err != nil {
		return fmt.Errorf("failed to insert audio file %s: %w", af.FilePath, err)
	}
	return nil
}

// IsLossless checks if a file extension typically denotes a lossless audio format.
func IsLossless(ext string) bool {
	switch strings.ToLower(ext) {
	case ".flac", ".wav", ".aiff", ".aif":
		return true
	default:
		return false
	}
}

// isAudioFile checks if a file has a common audio extension.
func isAudioFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".mp3", ".flac", ".wav", ".m4a", ".ogg", ".aac", ".wma", ".aiff", ".aif":
		return true
	default:
		return false
	}
}

// Indexer processes audio files and inserts their metadata into the database.
type Indexer struct {
	dbManager   *DBManager
	musicFolder string
	// Mutex to protect processedCount during concurrent updates
	countMu sync.Mutex
	processedCount int
	// Channel to send file paths to worker goroutines
	filePathChan chan string
	// WaitGroup to wait for all goroutines to finish
	wg sync.WaitGroup
	// Number of worker goroutines
	numWorkers int
}

// NewIndexer creates a new Indexer instance.
func NewIndexer(dbMgr *DBManager, folder string, numWorkers int) *Indexer {
	return &Indexer{
		dbManager:   dbMgr,
		musicFolder: folder,
		// Buffer the channel to allow some paths to be queued
		filePathChan: make(chan string, numWorkers*2), // Buffer size is a common heuristic
		numWorkers: numWorkers,
	}
}

// StartIndexing walks the music directory and processes each audio file.
func (i *Indexer) StartIndexing() error {
	log.Printf("Starting indexing of music folder: %s with %d workers", i.musicFolder, i.numWorkers)
	i.processedCount = 0

	// Start worker goroutines
	for w := 0; w < i.numWorkers; w++ {
		i.wg.Add(1) // Add one to the WaitGroup for each worker
		go i.worker()
	}

	// Walk the file path and send paths to the channel
	err := filepath.Walk(i.musicFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Preventing walk error for %q: %v", path, err)
			return err // Return the error to stop the walk
		}
		if info.IsDir() {
			return nil // Skip directories
		}

		if isAudioFile(info.Name()) {
			// Increment count safely
			i.countMu.Lock()
			i.processedCount++
			currentCount := i.processedCount
			i.countMu.Unlock()

			fmt.Printf("\rQueueing file %d: %s", currentCount, path) // Progress indicator
			i.filePathChan <- path // Send file path to the channel for a worker to pick up
		}
		return nil
	})

	// Close the channel to signal workers that no more paths will be sent
	close(i.filePathChan)

	// Wait for all worker goroutines to finish
	i.wg.Wait()

	fmt.Println("\nIndexing complete!")
	log.Printf("Processed %d audio files (attempts).", i.processedCount) // Note: This is count of files *attempted*
	return err
}

// worker processes file paths received from the filePathChan.
func (i *Indexer) worker() {
	defer i.wg.Done() // Signal that this worker is done when the function exits

	for filePath := range i.filePathChan {
		if err := i.processAudioFile(filePath); err != nil {
			log.Printf("Error processing %q: %v", filePath, err)
			// The error is logged, but the worker continues to the next file
		}
	}
}

// processAudioFile extracts metadata and inserts it into the database.
func (i *Indexer) processAudioFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", filePath, err)
	}
	defer file.Close()

	m, err := tag.ReadFrom(file)
	if err != nil {
		return fmt.Errorf("failed to read tags from %q: %w", filePath, err)
	}

	// Generate human-readable unique ID
	humanHashSource := fmt.Sprintf("%s-%s-%s-%s", m.Title(), m.Artist(), m.Album(), filePath)
	humanHashID, err := humanhash.Humanize([]byte(humanHashSource), 4)
	if err != nil {
		return fmt.Errorf("failed to generate humanhash for %q: %w", filePath, err)
	}

	trackNum, _ := m.Track()
	discNum, _ := m.Disc()

	audioFile := AudioFile{
		HumanHashID:     humanHashID,
		FilePath:        filePath,
		Title:           m.Title(),
		DurationSeconds: 0, // Set to 0, as dhowden/tag does not provide duration
		Lossless:        IsLossless(filepath.Ext(filePath)),
		TrackNumber:     trackNum,
		DiscNumber:      discNum,
		Year:            m.Year(),
		ArtistName:      m.Artist(),
		AlbumTitle:      m.Album(),
		GenreName:       m.Genre(),
	}

	// Ensure essential metadata is present
	if audioFile.Title == "" {
		audioFile.Title = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	}
	if audioFile.ArtistName == "" {
		audioFile.ArtistName = "Unknown Artist"
	}
	if audioFile.AlbumTitle == "" {
		audioFile.AlbumTitle = "Unknown Album"
	}

	// Database operations are protected by the DBManager's internal mutex
	var artistID int
	artistID, err = i.dbManager.GetOrInsertArtist(audioFile.ArtistName)
	if err != nil {
		return fmt.Errorf("failed to get/insert artist %q for %q: %w", audioFile.ArtistName, filePath, err)
	}
	audioFile.ArtistID = artistID

	var albumID int
	albumID, err = i.dbManager.GetOrInsertAlbum(audioFile.AlbumTitle, audioFile.ArtistID, audioFile.Year)
	if err != nil {
		return fmt.Errorf("failed to get/insert album %q by artist ID %d for %q: %w", audioFile.AlbumTitle, audioFile.ArtistID, filePath, err)
	}
	audioFile.AlbumID = albumID

	var genreID int
	if audioFile.GenreName != "" {
		genreID, err = i.dbManager.GetOrInsertGenre(audioFile.GenreName)
		if err != nil {
			return fmt.Errorf("failed to get/insert genre %q for %q: %w", audioFile.GenreName, filePath, err)
		}
	}
	audioFile.GenreID = genreID // Will be 0 if empty or not found

	// Insert the audio file record
	err = i.dbManager.InsertAudioFile(&audioFile)
	if err != nil {
		return fmt.Errorf("failed to insert audio file record %q: %w", filePath, err)
	}

	return nil
}

func main() {
	musicFolder := flag.String("music_folder", "", "Path to the music directory to index")
	dbPath := flag.String("db", "music_library.sqlite", "Path to the SQLite3 database file")
	numWorkers := flag.Int("workers", 4, "Number of concurrent workers for indexing") // New flag for concurrency

	flag.Parse()

	if *musicFolder == "" {
		log.Fatal("Error: --music_folder argument is required.")
	}
	if *dbPath == "" {
		log.Fatal("Error: --db argument is required.")
	}
	if *numWorkers <= 0 {
		log.Fatal("Error: --workers must be a positive integer.")
	}

	// Check if music folder exists
	if _, err := os.Stat(*musicFolder); os.IsNotExist(err) {
		log.Fatalf("Error: Music folder does not exist: %s", *musicFolder)
	}
	if stat, err := os.Stat(*musicFolder); err == nil && !stat.IsDir() {
		log.Fatalf("Error: Music folder path is not a directory: %s", *musicFolder)
	}

	log.Printf("Database path: %s", *dbPath)
	log.Printf("Music folder: %s", *musicFolder)
	log.Printf("Number of workers: %d", *numWorkers)

	dbMgr, err := NewDBManager(*dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database manager: %v", err)
	}
	defer dbMgr.Close()

	// Pass numWorkers to NewIndexer
	indexer := NewIndexer(dbMgr, *musicFolder, *numWorkers)
	if err := indexer.StartIndexing(); err != nil {
		log.Fatalf("Indexing failed: %v", err)
	}

	log.Println("Indexing process completed successfully!")
}


