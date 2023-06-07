package helpers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// Extracts the file from the request and returns file, handler, and error
func ExtractFileFromResponse(r *http.Request) (multipart.File, *multipart.FileHeader, error) {
	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		return nil, nil, err
	}
	// Get the file from the form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	return file, handler, nil
}

// Saves a copy of parameter file on the server. Takes parameters from request parsing and filepath to save to (eg. ./temp/)
func SaveACopyOfTheFileOnTheServer(file multipart.File, handler *multipart.FileHeader, filePath string) error {

	// Create a new file on the server to save the parameter file using the filename from the handler
	createdFile, err := os.Create(filePath + handler.Filename)
	// If error found
	if err != nil {
		// Create folder if it doesn't exist
		fmt.Print("Couldn't find directory, so creating new: ")
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
		// Retry creating the file
		fmt.Println("Retrying to create file")
		createdFile, err = os.Create(filePath + handler.Filename)
		if err != nil {
			return err
		}
	}
	defer createdFile.Close()

	// Copy the uploaded file content to the newly created file
	_, err = io.Copy(createdFile, file)
	if err != nil {
		return err
	}
	// Return no error
	return nil
}

// Deletes a file from the server
func DeleteFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return fmt.Errorf("failed to delete file: %s", err)
	}
	return nil
}
