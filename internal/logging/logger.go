package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	// InfoLogger logs general information
	InfoLogger *log.Logger
	// ErrorLogger logs errors
	ErrorLogger *log.Logger
	// DebugLogger logs debug information
	DebugLogger *log.Logger
	// logFile holds the log file handle
	logFile *os.File
)

// InitializeLogging sets up logging to both file and console
func InitializeLogging() error {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logPath := filepath.Join(logsDir, fmt.Sprintf("tf-engine_%s.log", timestamp))

	var err error
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Create multi-writer for both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Initialize loggers
	InfoLogger = log.New(multiWriter, "INFO:  ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(multiWriter, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	InfoLogger.Printf("Logging initialized: %s", logPath)
	return nil
}

// CloseLogging closes the log file
func CloseLogging() {
	if logFile != nil {
		InfoLogger.Println("Closing log file")
		logFile.Close()
	}
}

// LogStartup logs application startup information
func LogStartup() {
	InfoLogger.Println("=== TF-Engine 2.0 Starting ===")
	InfoLogger.Printf("Working directory: %s", mustGetWd())
	InfoLogger.Printf("Executable path: %s", mustGetExecutable())
	InfoLogger.Printf("Go version: %s", mustGetGoVersion())
	InfoLogger.Printf("OS: %s", mustGetOS())
}

// Helper functions to gather system info
func mustGetWd() string {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return wd
}

func mustGetExecutable() string {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return exe
}

func mustGetGoVersion() string {
	// This will be compiled in
	return "go1.21+"
}

func mustGetOS() string {
	return fmt.Sprintf("%s/%s", os.Getenv("GOOS"), os.Getenv("GOARCH"))
}

// LogPanic logs a recovered panic
func LogPanic(r interface{}) {
	ErrorLogger.Printf("PANIC: %v", r)
}

// CleanupOldLogs removes log files older than 30 days
func CleanupOldLogs() error {
	logsDir := "logs"
	entries, err := os.ReadDir(logsDir)
	if err != nil {
		return err
	}

	cutoff := time.Now().AddDate(0, 0, -30)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			path := filepath.Join(logsDir, entry.Name())
			if err := os.Remove(path); err != nil {
				DebugLogger.Printf("Failed to remove old log: %s", path)
			} else {
				DebugLogger.Printf("Removed old log: %s", path)
			}
		}
	}

	return nil
}
