package decoder

import (
	"bufio"
	"encoding/base64"
	"github.com/NeonHermit/b64-encode-decode-ggb/pkg/zipper"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

// var base64Pattern = regexp.MustCompile(`(?m)(?:[^A-Za-z0-9+/=]|^)((?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?)(?:[^A-Za-z0-9+/=]|$)`)
var base64Pattern = regexp.MustCompile(`(?m)(?:^|[^A-Za-z0-9+/=])([A-Za-z0-9+/]{32,}(?:==|[A-Za-z0-9+/]{1,2}=)?)(?:$|[^A-Za-z0-9+/=])`)

func ProcessFile(filePath string, outputDir string, unzip bool, infoLog *log.Logger, errorLog *log.Logger) {
	file, err := os.Open(filePath)
	if err != nil {
		errorLog.Printf("Error opening file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := base64Pattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) > 1 && len(match[1]) >= 4 {
				infoLog.Printf("Found possible base64 string in file %s: %s\n", filePath, match[1])

				decodeAndSaveBase64(match[1], filePath, outputDir, unzip, infoLog, errorLog)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		errorLog.Printf("Error reading file %s: %v\n", filePath, err)
	}
}

func decodeAndSaveBase64(b64Data, filePath, outputDir string, unzip bool, infoLog, errorLog *log.Logger) {
	decodedData, err := base64.StdEncoding.DecodeString(b64Data)
	if err != nil {
		errorLog.Println("Error decoding base64 data:", err)
		return
	}

	kind, _ := filetype.Match(decodedData)
	fileExtension := ".bin" // Default extension if none is found.
	if kind != types.Unknown {
		fileExtension = "." + kind.Extension
	} else {
		errorLog.Printf("Could not determine file extension from data, defaulting to .bin\n")
	}

	dirName := filepath.Base(filepath.Dir(filePath))
	targetDir := filepath.Join(outputDir, dirName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		errorLog.Println("Error creating directory:", err)
		return
	}

	outputFileName := "application" + fileExtension
	targetFilePath := filepath.Join(targetDir, outputFileName)

	if err := os.WriteFile(targetFilePath, decodedData, 0644); err != nil {
		errorLog.Println("Error writing decoded file:", err)
		return
	}

	infoLog.Printf("Decoded file written: %s\n", targetFilePath)

	if unzip && fileExtension == ".zip" {
		zipper.UnzipFile(targetFilePath, targetDir, infoLog, errorLog)
	}
}
