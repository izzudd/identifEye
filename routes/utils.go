package routes

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func saveFile(upFile *multipart.FileHeader, saveDir string, filename string) error {
	file, err := upFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	err = saveFileLocally(file, saveDir, filename)
	if err != nil {
		return err
	}

	return nil
}

func saveFileLocally(file multipart.File, saveDir string, filename string) error {
	// Create the "uploads" folder if it doesn't exist
	err := ensureDir(saveDir)
	if err != nil {
		return err
	}

	// Create the destination file
	destFile, err := createDestinationFile(saveDir, filename)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the file to the destination
	_, err = io.Copy(destFile, file)
	if err != nil {
		return err
	}

	return nil
}

func ensureDir(dirPath string) error {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func createDestinationFile(uploadPath, filename string) (*os.File, error) {
	destPath := filepath.Join(uploadPath, filename)
	destFile, err := os.Create(destPath)
	if err != nil {
		return nil, err
	}
	return destFile, nil
}
