package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

// Checksum reads the file at the provided path and calculates the sha256sum
func Sha256sum(filepath string) (string, error) {
	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read file '%s' while generating sha256sum: %w", filepath, err)
	}
	sumBytes := sha256.Sum256(fileBytes)
	return hex.EncodeToString(sumBytes[:]), nil
}

// Md5sum reads the file at the provided path and calculates the md5sum
func Md5sum(filepath string) (string, error) {
	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read file '%s' while generating md5sum: %w", filepath, err)
	}
	sumBytes := md5.Sum(fileBytes)
	return hex.EncodeToString(sumBytes[:]), nil
}
