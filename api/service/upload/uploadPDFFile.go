package upload

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func UploadPDFFile(file *multipart.FileHeader, tableID, username, messageUUID string) (string, string, error) {

	if username == "" || messageUUID == "" || tableID == "" {
		return "", "", fmt.Errorf("username, messageUUID and tableID are required")
	}

	contentType := file.Header.Get("Content-Type")
	if contentType != "application/pdf" {
		return "", "", fmt.Errorf("Content-Type must be application/PDF : %s", contentType)
	}

	if file.Size > 50<<20 {
		return "", "", fmt.Errorf("file too large %d bytes", file.Size)
	}

	normalized := strings.ReplaceAll(strings.ToLower(username), ".", "_")
	fileName := fmt.Sprintf("%s_%s.pdf", normalized, messageUUID)
	dirPath := fmt.Sprintf("./files/table_%s", tableID)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", "", fmt.Errorf("failed to create directory: %v", err)
	}

	filePath := filepath.Join(dirPath, fileName)

	return filePath, fileName, nil
}
