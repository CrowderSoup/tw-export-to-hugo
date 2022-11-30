package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	files, err := getFilesRecursively("./archive/tweets-md")
	if err != nil {
		log.Fatal(err)
	}

	contents, err := getFileContents(files)
	if err != nil {
		log.Fatal(err)
	}

	for filePath, fileContent := range contents {
		fmt.Printf("%s: %s\n\n", filePath, fileContent)
	}
}

func getFilesRecursively(dir string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return files, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			additionalFiles, err := getFilesRecursively(
				fmt.Sprintf("%s/%s", dir, entry.Name()))
			if err != nil {
				return files, err
			}

			files = append(files, additionalFiles...)
		} else {
			files = append(files, fmt.Sprintf("%s/%s", dir, entry.Name()))
		}
	}

	return files, nil
}

func getFileContents(files []string) (map[string]string, error) {
	log.Fatal("Not Implemented")
	return nil, nil
}
