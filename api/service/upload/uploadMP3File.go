package upload

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func UploadMP3File(file *multipart.FileHeader, tableID, username, messageUUID string) (string, string, error) {

	if username == "" || messageUUID == "" || tableID == "" {
		return "", "", fmt.Errorf("username,messageUUID and tableID are required")
	}

	contentType := file.Header.Get("Content-Type")
	if contentType != "audio/mp3" && contentType != "audio/mpeg" {
		return "", "", fmt.Errorf("Content-Type must be audio/mp3 or audio/mpeg : %s", contentType)
	}

	if file.Size > 20<<20 {
		return "", "", fmt.Errorf("file too large %d bytes", file.Size)
	}

	normalized := strings.ReplaceAll(strings.ToLower(username), ".", "_")
	dirPath := fmt.Sprintf("./files/table_%s/chat", tableID)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", "", fmt.Errorf("failed to create directory: %v", err)
	}
	fileName := fmt.Sprintf("%s_%s.mp3", normalized, messageUUID)

	filePath := filepath.Join(dirPath, fileName)

	return filePath, fileName, nil
}
