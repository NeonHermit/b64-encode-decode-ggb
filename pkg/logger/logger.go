package logger

import (
	"log"
	"os"
	"path/filepath"
)

func SetupLogging(outputDir, fileName string) (*log.Logger, func()) {
	if outputDir == "" {
		outputDir = "."
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatal("Failed to create output directory:", err)
	}

	logFile, err := os.OpenFile(filepath.Join(outputDir, fileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	cleanup := func() {
		if err := logFile.Close(); err != nil {
			log.Fatal("Failed to close log file:", err)
		}
	}

	return log.New(logFile, "", log.LstdFlags), cleanup
}
