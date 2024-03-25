package zipper

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
)

func ZipFiles(sourceDir, outputDir string) error {
	return filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && filePath != sourceDir {
			subDirOutputPath := filepath.Join(outputDir, filepath.Base(filePath))

			err := os.MkdirAll(subDirOutputPath, os.ModePerm)
			if err != nil {
				return err
			}

			zipFileName := filepath.Join(subDirOutputPath, "application.zip")
			return createZipForDir(filePath, zipFileName)
		}

		return nil
	})
}

func createZipForDir(dirPath, zipFileName string) error {
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	err = filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil // Ignore directories
		}

		fileInZip, err := archive.Create(info.Name())
		if err != nil {
			return err
		}

		sourceFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		_, err = io.Copy(fileInZip, sourceFile)
		return err
	})

	return err
}

func UnzipFile(zipPath, destDir string, infoLog, errorLog *log.Logger) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		errorLog.Println("Error opening zip file:", err)
		return
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			errorLog.Println("Error creating directory for unzipped file:", err)
			continue
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			errorLog.Println("Error opening destination file:", err)
			continue
		}

		rc, err := f.Open()
		if err != nil {
			errorLog.Println("Error opening file inside zip:", err)
			outFile.Close()
			continue
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			errorLog.Println("Error writing unzipped data to file:", err)
			continue
		}

		infoLog.Printf("Unzipped file written: %s\n", fpath)
	}
}
