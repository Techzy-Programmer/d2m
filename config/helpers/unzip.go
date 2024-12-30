package helpers

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UnzipFolder(zipPath, outputDir, customRootName string) error {
	gzipFile, err := os.Open(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open gzip file: %v", err)
	}
	defer gzipFile.Close()

	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzipReader.Close()

	reader := bufio.NewReader(gzipReader)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	for {
		filePath, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read file header: %v", err)
		}

		filePath = strings.TrimSpace(filePath)
		if filePath == "" {
			continue
		}

		pathParts := strings.Split(filePath, "/")
		if len(pathParts) > 0 && customRootName != "" {
			pathParts[0] = customRootName
			filePath = strings.Join(pathParts, "/")
		}

		fullPath := filepath.Join(outputDir, filePath)

		parentDir := filepath.Dir(fullPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", parentDir, err)
		}

		outputFile, err := os.Create(fullPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %v", fullPath, err)
		}

		// Reads file content until next header
		buffer := make([]byte, 4096)
		for {
			n, err := reader.Read(buffer)
			if err != nil && err != io.EOF {
				outputFile.Close()
				return fmt.Errorf("failed to read file content: %v", err)
			}

			if n == 0 {
				break
			}

			if _, err := outputFile.Write(buffer[:n]); err != nil {
				outputFile.Close()
				return fmt.Errorf("failed to write to file %s: %v", fullPath, err)
			}
		}

		outputFile.Close()
	}

	return nil
}
