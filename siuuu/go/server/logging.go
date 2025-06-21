package main

import (
	"io"
	"log"
	"os"
)

var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	archiveLog    *log.Logger
)

// setupLogging khởi tạo hệ thống logging
func setupLogging() {
	// System log file
	logFile, err := os.OpenFile("service.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open system log file: %v", err)
	}

	// Multi-writer for system log: file and stdout
	mw := io.MultiWriter(os.Stdout, logFile)

	InfoLogger = log.New(mw, "INFO: ", log.Ldate|log.Ltime)
	WarningLogger = log.New(mw, "WARNING: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile) // Thêm file và dòng cho lỗi

	// Logger mặc định sẽ dùng cho các lỗi FATAL
	log.SetOutput(mw)
	log.SetPrefix("FATAL: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Archive log file for agent messages
	archiveFile, err := os.OpenFile("archiver.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open archive log file: %v", err)
	}
	archiveLog = log.New(archiveFile, "", log.Ldate|log.Ltime) // Không có prefix, chỉ có ngày giờ
}
