package main

import (
	"flag"
	"github.com/NeonHermit/b64-encode-decode-ggb/pkg/decoder"
	"github.com/NeonHermit/b64-encode-decode-ggb/pkg/encoder"
	"github.com/NeonHermit/b64-encode-decode-ggb/pkg/logger"
	"github.com/NeonHermit/b64-encode-decode-ggb/pkg/zipper"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	inputDir := flag.String("input", "", "Directory to search for files")
	outputDir := flag.String("output", "", "Directory to output results")
	unzip := flag.Bool("unzip", false, "Unzip the decoded base64 files")
	zip := flag.Bool("zip", false, "Zip files without the parent dir")
	encode := flag.Bool("encode", false, "Encode files to Base64")
	replaceDir := flag.String("replace", "", "Directory to replace Base64 strings in matched folders")
	flag.Parse()

	infoLog, cleanupInfoLog := logger.SetupLogging(*outputDir, "info.log")
	defer cleanupInfoLog()
	errorLog, cleanupErrorLog := logger.SetupLogging(*outputDir, "error.log")
	defer cleanupErrorLog()

	allowedExtensions := map[string]bool{
		".html": true,
		".txt":  true,
	}

	if *outputDir == "" {
		*outputDir = "."
	}

	if *zip {
		err := zipper.ZipFiles(*inputDir, *outputDir)
		if err != nil {
			errorLog.Println("Error zipping files:", err)
		}
		return
	}

	if *encode {
		encoder.EncodeFilesInDir(*inputDir, infoLog, errorLog)

		if *replaceDir != "" {
			encoder.ReplaceBase64InDir(*inputDir, *replaceDir, allowedExtensions, infoLog, errorLog)
		}

		return
	}

	err := filepath.Walk(*inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errorLog.Println("Error walking through directory:", err)
			return err
		}

		ext := strings.ToLower(filepath.Ext(path))

		if !info.IsDir() && allowedExtensions[ext] {
			decoder.ProcessFile(path, *outputDir, *unzip, infoLog, errorLog)
		}

		return nil
	})

	if err != nil {
		errorLog.Println("Error walking through input directory:", err)
	}
}
