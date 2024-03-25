package encoder

import (
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func EncodeFilesInDir(sourceDir string, infoLog *log.Logger, errorLog *log.Logger) {
	err := filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			errorLog.Printf("Error accessing path %q: %v\n", filePath, err)
			return err
		}

		if !info.IsDir() {
			encodeFile(filePath, infoLog, errorLog)
		}

		return nil
	})

	if err != nil {
		errorLog.Printf("Error walking through source directory: %v\n", err)
	}
}

func encodeFile(filePath string, infoLog, errorLog *log.Logger) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		errorLog.Printf("Error reading file %s: %v\n", filePath, err)
		return
	}

	encoded := base64.StdEncoding.EncodeToString(data)

	// Strip the extension from the original file name
	baseName := filepath.Base(filePath)
	extension := filepath.Ext(baseName)
	outputFileName := strings.TrimSuffix(baseName, extension) + ".txt"

	// Use the directory of the original file
	outputDir := filepath.Dir(filePath)
	outputFilePath := filepath.Join(outputDir, outputFileName)

	err = os.WriteFile(outputFilePath, []byte(encoded), 0644)
	if err != nil {
		errorLog.Printf("Error writing encoded file %s: %v\n", outputFilePath, err)
		return
	}

	infoLog.Printf("Encoded file written: %s\n", outputFilePath)
}

func ReplaceBase64InDir(inputDir, replaceDir string, allowedExtensions map[string]bool, infoLog *log.Logger, errorLog *log.Logger) {
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != inputDir { // Skip the root directory
			// Find .txt file in the subdirectory
			err := filepath.Walk(path, func(subPath string, subInfo os.FileInfo, subErr error) error {
				if subErr != nil {
					return subErr
				}
				if !subInfo.IsDir() && strings.HasSuffix(subPath, ".txt") {
					dirName := filepath.Base(path)
					correspondingReplaceDir := filepath.Join(replaceDir, dirName)
					if _, err := os.Stat(correspondingReplaceDir); err == nil {
						replaceInMatchingDir(subPath, correspondingReplaceDir, allowedExtensions, infoLog, errorLog)
					}
					return filepath.SkipDir // Skip the rest of the directory once .txt file is found
				}
				return nil
			})
			if err != nil {
				errorLog.Printf("Error walking through subdirectory: %v\n", err)
			}
		}
		return nil
	})
	if err != nil {
		errorLog.Printf("Error walking through input directory: %v\n", err)
	}
}

func replaceInMatchingDir(encodedFilePath, replaceDir string, allowedExtensions map[string]bool, infoLog *log.Logger, errorLog *log.Logger) {
	encodedData, err := os.ReadFile(encodedFilePath)
	if err != nil {
		errorLog.Printf("Error reading encoded file %s: %v\n", encodedFilePath, err)
		return
	}
	encodedString := string(encodedData)

	var base64Pattern = regexp.MustCompile(`(?m)(?:^|[^A-Za-z0-9+/=])([A-Za-z0-9+/]{32,}(?:==|[A-Za-z0-9+/]{1,2}=)?)(?:$|[^A-Za-z0-9+/=])`)
	err = filepath.Walk(replaceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && allowedExtensions[filepath.Ext(path)] {
			replaceBase64InFile(path, base64Pattern, encodedString, infoLog, errorLog)
		}
		return nil
	})
	if err != nil {
		errorLog.Printf("Error walking through replace directory: %v\n", err)
	}
}

func replaceBase64InFile(filePath string, pattern *regexp.Regexp, replacement string, infoLog *log.Logger, errorLog *log.Logger) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		errorLog.Printf("Error reading file %s: %v\n", filePath, err)
		return
	}

	var base64PatternWithQuotes = regexp.MustCompile(`"([A-Za-z0-9+/]{32,}(?:==|[A-Za-z0-9+/]{1,2}=)?)"`)
	modifiedContent := base64PatternWithQuotes.ReplaceAllStringFunc(string(fileContent), func(s string) string {
		// Extract base64 content without the surrounding quotes
		base64Content := s[1 : len(s)-1]

		// Check if the captured content is a valid base64 string
		if _, err := base64.StdEncoding.DecodeString(base64Content); err == nil {
			// Replace only the Base64 part, keeping the surrounding quotes
			return "\"" + replacement + "\""
		}
		return s
	})

	err = os.WriteFile(filePath, []byte(modifiedContent), 0644)
	if err != nil {
		errorLog.Printf("Error writing to file %s: %v\n", filePath, err)
		return
	}

	infoLog.Printf("Base64 strings replaced in file: %s\n", filePath)
}
